package htmldown

import (
	"fmt"
	"strings"
)

type Node struct {
	Tag        string
	Text       string
	Attributes map[string]string
	Children   []*Node
	Parent     *Node

	IsEnd bool
}

func NewNode(tag string) *Node {
	return &Node{
		Tag:        tag,
		Attributes: make(map[string]string),
	}
}

func (node *Node) AddText(text string) {
	node.Text += text
}

func (node *Node) AddNode(n *Node) {
	node.Children = append(node.Children, n)
	node.Text += fmt.Sprintf("@node%d@", len(node.Children)-1)
	n.IsEnd = true
}

func (node *Node) Markdown() string {
	for i, child := range node.Children {
		node.Text = strings.Replace(
			node.Text,
			fmt.Sprintf("@node%d@", i),
			child.Markdown(),
			-1,
		)
	}
	if node.Tag == "code" && node.Attributes["class"] != "" {
		return fmt.Sprintf(`<%s class="%s">%s</%s>`, node.Tag, node.Attributes["class"], node.Text, node.Tag)
	}
	switch node.Tag {
	case "a":
		return fmt.Sprintf("[%s](%s)", node.Text, node.Attributes["href"])
	case "img":
		return fmt.Sprintf("\n![%s](%s)\n", node.Attributes["alt"], node.Attributes["src"])
	case "p", "br":
		return node.Text + "\n"
	case "strong", "b":
		return fmt.Sprintf("**%s**", node.Text)
	case "em":
		return fmt.Sprintf("**%s**", node.Text)
	case "hr":
		return "\n---\n"
	case "blockquote":
		return fmt.Sprintf("\n```\n%s\n```\n", node.Text)
	case "h1", "h2", "h3", "h4", "h5", "h6":
		repeat := int(node.Tag[len(node.Tag)-1]) - 48
		return fmt.Sprintf("%s %s", strings.Repeat("#", repeat), node.Text)
	default:
		if node.Tag != "" && node.Tag != "root" {
			return fmt.Sprintf("<%s>%s</%s>", node.Tag, node.Text, node.Tag)
		}
		if node.Text == "\n" {
			return ""
		}
		return node.Text
	}
	return node.Text
}
