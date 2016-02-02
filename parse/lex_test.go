package parse

import "testing"

func TestLex(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result []itemType
	}{
		{
			"test1",
			"",
			[]itemType{
				itemEOF,
			},
		},
		{
			"test2",
			"<==>",
			[]itemType{
				itemLeftDelim,
				itemRightDelim,
				itemEOF,
			},
		},
		{
			"test2",
			"<=if=>",
			[]itemType{
				itemLeftDelim,
				itemIf,
				itemRightDelim,
				itemEOF,
			},
		},
	}

	for _, test := range tests {
		l := lex(test.name, test.input)

		for _, typ := range test.result {
			i := <-l.items
			if i.typ != typ {
				t.Errorf("expected %q token got %q", i, typ)
			}
		}
	}
}
