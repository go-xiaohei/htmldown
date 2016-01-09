package ng

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type (
	// Node defines a simple description iof html node
	Node interface {
		New(t html.Token, parent Node) Node
		Tag() string                   // node's tag name
		Text(str ...string) string     // get or append node's text
		Attr() map[string]string       // node's attributes
		Parent() Node                  // node's parent
		Children(nodes ...Node) []Node // append or get node's children nodes
		Markdown() string              // render nodes to markdown
	}
	baseNode struct {
		tag        string
		text       string
		attributes map[string]string
		parent     Node
		children   []Node
	}
	// RootNode is root node of html document
	RootNode struct {
		baseNode
	}
	// CommonNode is a node that print as normal html, not markdown
	CommonNode struct {
		baseNode
	}
)

// New returns new baseNode
func (bn *baseNode) New(t html.Token, parent Node) Node {
	b := &baseNode{
		tag:        strings.ToLower(t.Data),
		attributes: make(map[string]string),
		parent:     parent,
	}
	for _, attr := range t.Attr {
		b.attributes[attr.Key] = attr.Val
	}
	return b
}

func (bn *baseNode) Tag() string { return bn.tag }

func (bn *baseNode) Attr() map[string]string { return bn.attributes }

func (bn *baseNode) Text(str ...string) string {
	if len(str) > 0 {
		bn.text += strings.Join(str, "")
	}
	return bn.text
}

func (bn *baseNode) Parent() Node { return bn.parent }

func (bn *baseNode) Children(nodes ...Node) []Node {
	if len(nodes) > 0 {
		bn.children = append(bn.children, nodes...)
		bn.Text(strings.Repeat("@node@", len(nodes)))
	}
	return bn.children
}

func (bn *baseNode) Markdown() string {
	content := bn.text
	for _, child := range bn.children {
		content = strings.Replace(content, "@node@", child.Markdown(), 1)
	}
	return content
}

// Tag returns RootNode's tag ,"root"
func (rt *RootNode) Tag() string { return "root" }

// Markdown renders CommonNode as html tag string ,<tag>text</tag>
func (cn *CommonNode) Markdown() string {
	return fmt.Sprintf("<%s>%s</%s>", cn.Tag(), cn.baseNode.Markdown(), cn.Tag())
}

// New returns new CommonNode
func (cn *CommonNode) New(t html.Token, parent Node) Node {
	bs := cn.baseNode.New(t, parent)
	return &CommonNode{
		baseNode: *(bs.(*baseNode)),
	}
}

type (
	// ParagraphNode is used to <p>
	ParagraphNode struct{ baseNode }
	// HrefNode is used to <a>
	HrefNode struct{ baseNode }
	// ImageNode is used to <img>
	ImageNode struct{ baseNode }
	// HorizontalNode is used to <hr/>
	HorizontalNode struct{ baseNode }
	// BreakNode is used to <br/>
	BreakNode struct{ baseNode }
	// BlockquoteNode is used to <blockquote>
	BlockquoteNode struct{ baseNode }
	// CodeNode is used to <code>
	CodeNode struct{ baseNode }
)

// Markdown renders ParagraphNode with \n
func (pn *ParagraphNode) Markdown() string {
	if pn.Text() == "" {
		return ""
	}
	return pn.baseNode.Markdown() + "\n\n"
}

// New returns new ParagraphNode
func (pn *ParagraphNode) New(t html.Token, parent Node) Node {
	bs := pn.baseNode.New(t, parent)
	return &ParagraphNode{
		baseNode: *(bs.(*baseNode)),
	}
}

// Markdown renders HrefNode as [text](href)
func (hn *HrefNode) Markdown() string {
	return fmt.Sprintf("[%s](%s)", hn.baseNode.Markdown(), hn.attributes["href"])
}

// New returns new HrefNode
func (hn *HrefNode) New(t html.Token, parent Node) Node {
	bs := hn.baseNode.New(t, parent)
	return &HrefNode{
		baseNode: *(bs.(*baseNode)),
	}
}

// Markdown renders ImageNode as ![alt](src)
func (in *ImageNode) Markdown() string {
	return fmt.Sprintf("![%s](%s)", in.attributes["alt"], in.attributes["src"])
}

// New returns new ImageNode
func (in *ImageNode) New(t html.Token, parent Node) Node {
	bs := in.baseNode.New(t, parent)
	return &ImageNode{
		baseNode: *(bs.(*baseNode)),
	}
}

// Markdown renders HorizontalNode as `---\n`
func (hn *HorizontalNode) Markdown() string {
	return "---\n\n"
}

// New returns HorizontalNode
func (hn *HorizontalNode) New(t html.Token, parent Node) Node {
	bs := hn.baseNode.New(t, parent)
	return &HorizontalNode{
		baseNode: *(bs.(*baseNode)),
	}
}

// Markdown renders BreakNode as "\n\n"
func (bn *BreakNode) Markdown() string {
	return "\n\n"
}

// New returns new BreakNode
func (bn *BreakNode) New(t html.Token, parent Node) Node {
	bs := bn.baseNode.New(t, parent)
	return &BreakNode{
		baseNode: *(bs.(*baseNode)),
	}
}

// Markdown renders BlockquoteNode as  \n```\n text \n```\n
func (bn *BlockquoteNode) Markdown() string {
	return fmt.Sprintf("\n```\n%s\n```\n", bn.baseNode.Markdown())
}

// New returns new BlockquoteNode
func (bn *BlockquoteNode) New(t html.Token, parent Node) Node {
	bs := bn.baseNode.New(t, parent)
	return &BlockquoteNode{
		baseNode: *(bs.(*baseNode)),
	}
}

// Markdown renders CodeNode as <code>text</code>
// if element class exist, render as <code class="class">text</code>
func (cn *CodeNode) Markdown() string {
	class := cn.attributes["class"]
	if class != "" {
		return fmt.Sprintf(`<code class="%s">%s</code>`, class, cn.baseNode.Markdown())
	}
	return fmt.Sprintf("<code>%s</code>", cn.baseNode.Markdown())
}

// New returns new CodeNode
func (cn *CodeNode) New(t html.Token, parent Node) Node {
	bs := cn.baseNode.New(t, parent)
	return &CodeNode{
		baseNode: *(bs.(*baseNode)),
	}
}
