package soup

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/niklasfasching/css"
	"golang.org/x/net/html"
)

func Parse(r io.Reader) (*Node, error) {
	htmlNode, err := html.Parse(r)
	return AsNode(htmlNode), err
}

func MustParse(r io.Reader) *Node {
	n, err := Parse(r)
	if err != nil {
		panic(err)
	}
	return n
}

func Load(client *http.Client, url string) (*Node, error) {
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return Parse(res.Body)
}

func MustLoad(client *http.Client, url string) *Node {
	n, err := Load(client, url)
	if err != nil {
		panic(err)
	}
	return n
}

func (n *Node) First(s string) *Node { return n.FirstSel(css.MustCompile(s)) }
func (n *Node) FirstSel(s css.Selector) *Node {
	if n == nil {
		return nil
	}
	return AsNode(css.First(s, AsHTMLNode(n)))
}

func (n *Node) All(s string) Nodes { return n.AllSel(css.MustCompile(s)) }
func (n *Node) AllSel(s css.Selector) Nodes {
	if n == nil {
		return nil
	}
	htmlNodes := css.All(s, AsHTMLNode(n))
	return AsNodes(&htmlNodes)
}

func (n *Node) Text() string {
	var out strings.Builder
	appendText(&out, AsHTMLNode(n))
	return out.String()
}

func (n *Node) TrimmedText() string {
	return trimmed(n.Text())
}

func (n *Node) OuterHTML() string {
	if n == nil {
		return ""
	}
	var out strings.Builder
	if err := html.Render(&out, AsHTMLNode(n)); err != nil {
		panic(fmt.Sprintf("Could not render html: %s", err))
	}
	return out.String()
}

func (n *Node) HTML() string {
	if n == nil {
		return ""
	}
	var out strings.Builder
	for n := AsHTMLNode(n).FirstChild; n != nil; n = n.NextSibling {
		if err := html.Render(&out, n); err != nil {
			panic(fmt.Sprintf("Could not render html: %s", err))
		}
	}
	return out.String()
}

func (n *Node) Attribute(key string) string {
	if n == nil {
		return ""
	}
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func (ns Nodes) Eq(i int) *Node {
	if len(ns) >= i {
		return nil
	}
	return ns[i]
}

func (ns Nodes) Len() int {
	return len(ns)
}

func (ns Nodes) Text(sep string) string {
	ss := make([]string, len(ns))
	for i, n := range ns {
		ss[i] = n.Text()
	}
	return strings.Join(ss, sep)
}

func (ns Nodes) Attribute(key string) []string {
	as := make([]string, len(ns))
	for i, n := range ns {
		as[i] = n.Attribute(key)
	}
	return as
}

func (ns Nodes) First(s string) *Node { return ns.FirstSel(css.MustCompile(s)) }
func (ns Nodes) FirstSel(s css.Selector) *Node {
	for _, n := range ns {
		if f := n.FirstSel(s); f != nil {
			return f
		}
	}
	return nil
}

func (ns Nodes) All(s string) Nodes { return ns.AllSel(css.MustCompile(s)) }
func (ns Nodes) AllSel(s css.Selector) Nodes {
	all := []*Node{}
	for _, n := range ns {
		all = append(all, n.AllSel(s)...)
	}
	return all
}

func (ns Nodes) HTML() string {
	ss := make([]string, len(ns))
	for i, n := range ns {
		ss[i] = n.OuterHTML()
	}
	return strings.Join(ss, "\n")
}
