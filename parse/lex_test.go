package parse

import "testing"

var (
	tLeftDelim  = item{itemLeftDelim, 0, "<="}
	tRightDelim = item{itemRightDelim, 0, "=>"}
	tSpace      = item{itemSpace, 0, " "}
	tLeftParen  = item{itemLeftParen, 0, "("}
	tRightParen = item{itemRightParen, 0, ")"}
	tLeftBrace  = item{itemLeftBrace, 0, "{"}
	tRightBrace = item{itemRightBrace, 0, "}"}
	tColonEqual = item{itemColonEqual, 0, ":="}
	tFunc       = item{itemFunc, 0, "func"}
	tIf         = item{itemIf, 0, "if"}
	tEOF        = item{itemEOF, 0, ""}
)

type lexTest struct {
	name  string
	input string
	items []item
}

var tests = []lexTest{
	{
		"test1",
		"",
		[]item{
			tEOF,
		},
	},
	{
		"test2",
		"<==>",
		[]item{
			tLeftDelim,
			tRightDelim,
			tEOF,
		},
	},
	{
		"test3",
		`<= "hello" =>`,
		[]item{
			tLeftDelim,
			tSpace,
			item{itemString, 3, "\"hello\""},
			tSpace,
			tRightDelim,
			item{typ: itemEOF},
		},
	},
	{
		"test4",
		`<= ("hello") =>`,
		[]item{
			tLeftDelim,
			tSpace,
			tLeftParen,
			item{itemString, 3, "\"hello\""},
			tRightParen,
			tSpace,
			tRightDelim,
			item{typ: itemEOF},
		},
	},
	{
		"test5",
		`<= func ("hello") =>`,
		[]item{
			tLeftDelim,
			tSpace,
			tFunc,
			tSpace,
			tLeftParen,
			item{itemString, 3, "\"hello\""},
			tRightParen,
			tSpace,
			tRightDelim,
			item{typ: itemEOF},
		},
	},
	{
		"test6",
		`<= func test(hello) {} =>`,
		[]item{
			tLeftDelim,
			tSpace,
			tFunc,
			tSpace,
			item{itemIdentifier, 3, "test"},
			tLeftParen,
			item{itemIdentifier, 3, "hello"},
			tRightParen,
			tSpace,
			tLeftBrace,
			tRightBrace,
			tSpace,
			tRightDelim,
			item{typ: itemEOF},
		},
	},
	{
		"test7",
		`<= hello := "world" =>`,
		[]item{
			tLeftDelim,
			tSpace,
			item{itemIdentifier, 0, "hello"},
			tSpace,
			tColonEqual,
			tSpace,
			item{itemString, 3, "\"world\""},
			tSpace,
			tRightDelim,
			item{typ: itemEOF},
		},
	},
	{
		"Test Complex Numbers",
		"<= 1+2i =>",
		[]item{
			tLeftDelim,
			tSpace,
			item{itemComplex, 0, "1+2i"},
			tSpace,
			tRightDelim,
			tEOF,
		},
	},
	{
		"Test decimal numbers",
		"<= 10.2 =>",
		[]item{
			tLeftDelim,
			tSpace,
			item{itemNumber, 0, "10.2"},
			tSpace,
			tRightDelim,
			tEOF,
		},
	},
	{
		"Test hexadecimal numbers",
		"<= 0x000fff =>",
		[]item{
			tLeftDelim,
			tSpace,
			item{itemNumber, 0, "0x000fff"},
			tSpace,
			tRightDelim,
			tEOF,
		},
	},
	{
		"Test Operators",
		"<= if (1 > 2) {} =>",
		[]item{
			tLeftDelim,
			tSpace,
			tIf,
			tSpace,
			tLeftParen,
			item{itemNumber, 0, "1"},
			tSpace,
			item{itemOperator, 0, ">"},
			tSpace,
			item{itemNumber, 0, "2"},
			tRightParen,
			tSpace,
			tLeftBrace,
			tRightBrace,
			tSpace,
			tRightDelim,
			tEOF,
		},
	},
	{
		"Test Unexpected }",
		"<= } =>",
		[]item{
			tLeftDelim,
			tSpace,
			tRightBrace,
			item{itemError, 0, "unexpected closing brace [U+007D '}']"},
		},
	},
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest) (items []item) {
	l := lex(t.name, t.input)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

func equal(i1, i2 []item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val {
			return false
		}
		if checkPos && i1[k].pos != i2[k].pos {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range tests {
		items := collect(&test)
		if !equal(items, test.items, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, items, test.items)
		}
	}
}
