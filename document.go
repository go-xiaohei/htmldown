package htmldown

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type (
	// Document maintains nodes of html content
	Document struct {
		htmlReader io.Reader

		nodes    map[string]Node
		rootNode Node
		isParsed bool
	}
)

// NewDocument returns new document with io.Reader
func NewDocument(r io.Reader) *Document {
	return &Document{
		htmlReader: r,
		nodes: map[string]Node{
			"p":          new(ParagraphNode),
			"a":          new(HrefNode),
			"img":        new(ImageNode),
			"hr":         new(HorizontalNode),
			"br":         new(BreakNode),
			"blockquote": new(BlockquoteNode),
			"code":       new(CodeNode),
			"pre":        new(PreNode),
			"h1":         new(HeaderNode),
			"h2":         new(HeaderNode),
			"h3":         new(HeaderNode),
			"h4":         new(HeaderNode),
			"h5":         new(HeaderNode),
			"h6":         new(HeaderNode),
			"strong":     new(StrongNode),
			"em":         new(StrongNode),
		},
		rootNode: &RootNode{},
	}
}

func (doc *Document) parse() error {
	if doc.isParsed {
		return nil
	}
	var (
		z       = html.NewTokenizer(doc.htmlReader)
		current = doc.rootNode
		err     error
	)

	for {
		next := z.Next()
		if next == html.ErrorToken {
			if z.Err() != io.EOF {
				err = z.Err()
			}
			break
		}
		t := z.Token()

		// start tag, self close tag
		if next == html.StartTagToken || next == html.SelfClosingTagToken {
			tag := strings.ToLower(t.Data)
			n := new(CommonNode).New(t, current)
			if n2, ok := doc.nodes[tag]; ok {
				n = n2.New(t, current)
			}
			if next == html.StartTagToken {
				current = n
			} else {
				current.Children(n)
			}
			continue
		}

		// text tag
		if next == html.TextToken {
			if len(t.Data) > 0 && t.Data != "\n" {
				current.Text(t.Data)
			}
			continue
		}

		// end tag
		if next == html.EndTagToken {
			current.Parent().Children(current)
			current = current.Parent()
		}
	}

	if current != nil && current.Parent() != nil {
		current.Parent().Children(current)
	}

	doc.isParsed = true
	return err
}

// Root returns root node of the document
func (doc *Document) Root() Node {
	return doc.rootNode
}

// SetNode sets node by tag name.
// so it parses element by this node
func (doc *Document) SetNode(tag string, n Node) {
	doc.nodes[tag] = n
}

// Markdown renders this document to markdown content
func (doc *Document) Markdown() (string, error) {
	if !doc.isParsed {
		if err := doc.parse(); err != nil {
			return "", err
		}
	}
	return doc.rootNode.Markdown(), nil
}
