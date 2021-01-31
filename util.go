package soup

import (
	"regexp"
	"strings"
	"unsafe"

	"golang.org/x/net/html"
)

type Node html.Node
type Nodes []*Node

func AsHTMLNode(n *Node) *html.Node  { return (*html.Node)(unsafe.Pointer(n)) }
func AsNode(n *html.Node) *Node      { return (*Node)(unsafe.Pointer(n)) }
func AsNodes(ns *[]*html.Node) Nodes { return *(*[]*Node)(unsafe.Pointer(ns)) }

var duplicateWhitespace = regexp.MustCompile(`\s+(\n)\s*|\s*(\n)\s+|(\s)\s+`)

func appendText(out *strings.Builder, n *html.Node) {
	switch {
	case n == nil || n.Type == html.CommentNode:
		return
	case n.Type == html.TextNode:
		out.WriteString(n.Data)
	default:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			appendText(out, c)
		}
	}
}

func trimmed(s string) string {
	return duplicateWhitespace.ReplaceAllString(strings.TrimSpace(s), "$1")
}
