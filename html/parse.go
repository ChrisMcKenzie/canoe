package html

import (
	"fmt"
	"io"

	"github.com/google/uuid"

	"golang.org/x/net/html"
)

const (
	// DefaultFragmentTag defines the tag name that will used by default in the
	// parser
	DefaultFragmentTag = "canoe-fragment"
)

var (
	uuidSpace = uuid.Must(uuid.Parse("8C6CDB3E-A15F-4C3B-992B-58D923D50BD6"))
)

// Fragment ...
type Fragment struct {
	ID   string
	Href string
	node *html.Node
}

// Parser represents a html parse that can stream parse a canoe template
type Parser struct {
	doc   *html.Node
	frags chan Fragment
	tag   string
}

// OptionFunc is a type used to modify options in the Parser
type OptionFunc func(*Parser)

// NewParser takes in an io.Reader of html to be parsed.
func NewParser(t io.Reader, opts ...OptionFunc) (*Parser, error) {
	doc, err := html.Parse(t)
	if err != nil {
		return nil, fmt.Errorf("unable to load template: %v", err)
	}

	p := &Parser{doc, make(chan Fragment), DefaultFragmentTag}
	for _, opt := range opts {
		opt(p)
	}

	go p.parse()

	return p, nil
}

func (p *Parser) Render(w io.Writer) error {
	return html.Render(w, p.doc)
}

func (p *Parser) parse() {
	p.parseFragment(p.doc)
	defer close(p.frags)
}

func (p *Parser) parseFragment(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == p.tag {
		f := Fragment{
			node: n,
		}
		for i, a := range n.Attr {
			if a.Key == "href" {
				f.Href = a.Val

				// remove the href from the tag
				n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			}
		}

		f.ID = generateID(f.Href)
		n.Attr = append(n.Attr, html.Attribute{
			Key: "fragment",
			Val: f.ID,
		})
		p.frags <- f
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.parseFragment(c)
	}
}

// Fragments returns a read-only channel of Fragment that have been parsed
func (p *Parser) Fragments() <-chan Fragment {
	return p.frags
}

func generateID(href string) string {
	id := uuid.NewSHA1(uuidSpace, []byte(href))
	return id.String()
}
