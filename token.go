package htmldown

import (
	"strings"

	"golang.org/x/net/html"
)

func Markdown(htmlContent string) string {
	z := html.NewTokenizer(strings.NewReader(htmlContent))
	tree := NewNode("root")

	var currentNode *Node
	for {
		next := z.Next()
		if next == html.ErrorToken {
			break
		}
		t := z.Token()
		n := NewNode(strings.ToLower(t.Data))
		for _, attr := range t.Attr {
			n.Attributes[strings.ToLower(attr.Key)] = attr.Val
		}
		n.Parent = tree
		if currentNode != nil {
			n.Parent = currentNode
		}

		if next == html.StartTagToken {
			// n.Parent.AddNode(n)
			currentNode = n
			continue
		}
		if next == html.TextToken {
			if currentNode == nil {
				tree.AddText(t.Data)
			} else {
				currentNode.AddText(t.Data)
			}
			continue
		}
		if next == html.SelfClosingTagToken {
			currentNode.AddNode(n)
			continue
		}
		if next == html.EndTagToken {
			currentNode.Parent.AddNode(currentNode)
			currentNode = currentNode.Parent
		}
	}

	if currentNode != nil && currentNode.Parent != nil && !currentNode.IsEnd {
		currentNode.Parent.AddNode(currentNode)
	}

	return tree.Markdown()
}
