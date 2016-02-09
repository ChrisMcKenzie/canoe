package parse

type Tree struct {
	Name      string
	Root      *ListNode
	text      string
	funcs     []map[string]interface{}
	lex       *lexer
	token     [3]item
	peekCount int
	vars      []string
}

func Parse(name, text string, funcs ...map[string]interface{}) (treeSet map[string]*Tree, err error) {
	treeSet = make(map[string]*Tree)
	t := New(name)
	t.text = text
	// _, err = t.Parse(text, treeSet, funcs...)

	return
}

func (t *Tree) Copy() *Tree {
	if t == nil {
		return nil
	}

	return &Tree{
		Name: t.Name,
		// Root: t.Root.CopyList()
		text: t.text,
	}
}

func (t *Tree) next() item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lex.nextItem()
	}

	return t.token[t.peekCount]
}

func (t *Tree) backup() {
	t.peekCount++
}

func (t *Tree) backup2(t1 item) {
	t.token[1] = t1
	t.peekCount = 2
}

func (t *Tree) backup3(t1, t2 item) {
	t.token[1] = t1
	t.token[2] = t2
	t.peekCount = 3
}

func (t *Tree) peek() item {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}

	t.peekCount = 1
	t.token[0] = t.lex.nextItem()
	return t.token[0]
}

func (t *Tree) nextNonSpace() (token item) {
	for {
		token = t.next()
		if token.typ == itemSpace {
			break
		}
	}

	return token
}

func (t *Tree) peekNonSpace() (token item) {
	for {
		token = t.next()
		if token.typ == itemSpace {
			break
		}
	}
	t.backup()
	return token
}

func New(name string, funcs ...map[string]interface{}) *Tree {
	return &Tree{
		Name:  name,
		funcs: funcs,
	}
}
