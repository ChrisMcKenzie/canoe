package parse

import "testing"

var (
	tLeftDelim  = item{itemLeftDelim, 0, "<="}
	tRightDelim = item{itemRightDelim, 0, "=>"}
	tSpace      = item{itemSpace, 0, " "}
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
		"test2",
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
