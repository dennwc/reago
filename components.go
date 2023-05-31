package reago

import (
	"bytes"
	"html/template"
	"io"
	"os"

	"golang.org/x/net/html"
)

// Context is passed to each component when rendering.
//
// Context has a generic parameter for data interface, which can be customized by client code.
type Context[T any] struct {
	// Content is an inner HTML of a component tag.
	Content template.HTML
	// Attrs is a map of HTML attributes of the component tag.
	Attrs map[string]string
	// DB is an interface for accessing external data.
	DB T
}

func NewEngine[T any](compPath string, db T) (*Engine[T], error) {
	e := &Engine[T]{db: db}
	if err := e.readComponents(compPath); err != nil {
		return nil, err
	}
	return e, nil
}

// Engine renders HTML page with components. See RenderPage for details.
type Engine[T any] struct {
	root *template.Template
	db   T
}

// readComponents reads gohtml components from a given directory.
func (e *Engine[T]) readComponents(path string) error {
	tmpl, err := template.ParseGlob(path + "/*.gohtml")
	if err != nil {
		return err
	}
	e.root = tmpl
	return nil
}

// RenderPage renders an HTML page by replacing component tags with rendered versions.
//
// It scans HTML page for tags that have a matching <tag>.gohtml file in the component directory.
// This component is then rendered from a template file and is appended back to the main HTML.
func (e *Engine[T]) RenderPage(w io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// TODO: Support gohtml templates for main pages.
	// TODO: This is pretty inefficient. Ideally we should use HTML tokenizer.
	root, err := html.Parse(f)
	if err != nil {
		return err
	}
	root, err = e.renderNode(root)
	if err != nil {
		return err
	}
	return html.Render(w, root)
}

// renderNode scans a single HTML node for component tags and replaces them with rendered versions.
func (e *Engine[T]) renderNode(n *html.Node) (*html.Node, error) {
	// TODO: Cache components with the same parameters?
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		c2, err := e.renderNode(c)
		if err != nil {
			return nil, err
		}
		if c != c2 {
			n.InsertBefore(c2, c)
			n.RemoveChild(c)
			c = c2
		}
	}
	if n.Type != html.ElementNode {
		return n, nil
	}
	t := e.root.Lookup(n.Data)
	if t == nil {
		return n, nil
	}

	m := make(map[string]string, len(n.Attr))
	for _, a := range n.Attr {
		m[a.Key] = a.Val
	}
	var content template.HTML
	if n.FirstChild != nil {
		inner := &html.Node{Type: html.DocumentNode}
		for n.FirstChild != nil {
			c := n.FirstChild
			n.RemoveChild(c)
			inner.AppendChild(c)
		}
		var buf bytes.Buffer
		err := html.Render(&buf, inner)
		if err != nil {
			return nil, err
		}
		content = template.HTML(buf.String())
	}
	ctx := &Context[T]{Content: content, Attrs: m, DB: e.db}
	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		return nil, err
	}
	nroot, err := html.Parse(&buf)
	if err != nil {
		return nil, err
	}
	n2 := nroot.FirstChild.LastChild.FirstChild
	n2.Parent.RemoveChild(n2)
	return n2, nil
}
