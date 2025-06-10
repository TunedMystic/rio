// Package dom provides a library for programmatic HTML generation.
package dom

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

// Node is the basic interface for all renderable entities in the DOM.
type Node interface {
	Render(w io.Writer) error
}

// ------------------------------------------------------------------
//
// DOM elements
//
// ------------------------------------------------------------------

// htmlElement represents an HTML element.
type htmlElement struct {
	Name     string
	IsVoid   bool
	Children []Node
}

var _ Node = (*htmlElement)(nil)
var _ fmt.Stringer = (*htmlElement)(nil)

func (e *htmlElement) Render(w io.Writer) error {
	var err error
	if _, err = w.Write(bLt); err != nil {
		return err
	}
	if _, err = io.WriteString(w, e.Name); err != nil {
		return err
	}

	// Render attributes
	for _, child := range e.Children {
		if attr, ok := child.(HtmlAttributer); ok {
			if err = attr.RenderAttribute(w); err != nil {
				return err
			}
		}
	}

	if _, err = w.Write(bGt); err != nil {
		return err
	}

	// Void elements have no children or closing tags.
	if e.IsVoid {
		return nil
	}

	// Render children
	for _, child := range e.Children {
		if _, ok := child.(HtmlAttributer); !ok {
			if err = child.Render(w); err != nil {
				return err
			}
		}
	}

	if _, err = w.Write(bLtSlash); err != nil {
		return err
	}
	if _, err = io.WriteString(w, e.Name); err != nil {
		return err
	}
	if _, err = w.Write(bGt); err != nil {
		return err
	}
	return nil
}

func (e *htmlElement) String() string {
	var b strings.Builder
	e.Render(&b)
	return b.String()
}

// ------------------------------------------------------------------
//
// DOM attributes
//
// ------------------------------------------------------------------

// HtmlAttributer defines the interface for nodes that can be rendered as HTML attributes.
type HtmlAttributer interface {
	RenderAttribute(w io.Writer) error
}

// htmlAttr represents an HTML attribute.
type htmlAttr struct {
	Name  string
	Value string
}

var _ Node = (*htmlAttr)(nil)
var _ fmt.Stringer = (*htmlAttr)(nil)
var _ HtmlAttributer = (*htmlAttr)(nil)

func (a *htmlAttr) Render(w io.Writer) error {
	var err error
	if _, err = w.Write(bSpace); err != nil {
		return err
	}
	if _, err = io.WriteString(w, a.Name); err != nil {
		return err
	}

	// Boolean attributes have no value.
	if a.Value == "" {
		return nil
	}

	if _, err = w.Write(bEqualsQuote); err != nil {
		return err
	}

	if needsEscaping(a.Value) {
		template.HTMLEscape(w, []byte(a.Value))
	} else {
		_, err = io.WriteString(w, a.Value)
	}

	if err != nil {
		return err
	}

	if _, err = w.Write(bQuote); err != nil {
		return err
	}
	return nil
}

func (a *htmlAttr) RenderAttribute(w io.Writer) error {
	return a.Render(w)
}

func (a *htmlAttr) String() string {
	var b strings.Builder
	a.Render(&b)
	return b.String()
}

// ------------------------------------------------------------------
//
// DOM strings
//
// ------------------------------------------------------------------

// htmlSafe represents a string that will be HTML-escaped upon rendering.
type htmlSafe string

var _ Node = htmlSafe("")
var _ fmt.Stringer = htmlSafe("")

func (s htmlSafe) Render(w io.Writer) error {
	val := string(s)
	if needsEscaping(val) {
		template.HTMLEscape(w, []byte(val))
		return nil
	}
	_, err := io.WriteString(w, val)
	return err
}

func (s htmlSafe) String() string {
	var b strings.Builder
	s.Render(&b)
	return b.String()
}

// htmlRaw represents a string that will be rendered as-is, without HTML escaping.
type htmlRaw string

var _ Node = htmlRaw("")
var _ fmt.Stringer = htmlRaw("")

func (s htmlRaw) Render(w io.Writer) error {
	_, err := io.WriteString(w, string(s))
	return err
}

func (s htmlRaw) String() string {
	return string(s)
}

// ------------------------------------------------------------------
//
// DOM control structures
//
// ------------------------------------------------------------------

// Group is a convenience type for grouping multiple Nodes together.
type Group []Node

var _ Node = (Group)(nil)
var _ fmt.Stringer = (Group)(nil)

func (g Group) Render(w io.Writer) error {
	for _, node := range g {
		if err := node.Render(w); err != nil {
			return err
		}
	}
	return nil
}

func (g Group) String() string {
	var b strings.Builder
	g.Render(&b)
	return b.String()
}

// Map is a utility function that maps a slice of items to a Group of Nodes.
func Map[T any](items []T, fn func(T) Node) Group {
	nodes := make([]Node, 0, len(items))
	for _, item := range items {
		nodes = append(nodes, fn(item))
	}
	return nodes
}

var emptyNode Node = htmlRaw("")

// If is a utility function that returns a Node if the condition is true,
// otherwise returns a shared empty node.
func If(condition bool, a Node) Node {
	if condition {
		return a
	}
	return emptyNode
}

// Ifelse is a utility function that returns Node a if the condition is true,
// otherwise returns Node b.
func Ifelse(condition bool, a, b Node) Node {
	if condition {
		return a
	}
	return b
}

// ------------------------------------------------------------------
//
// Dom creation functions
//
// ------------------------------------------------------------------

func CreateElement(name string, children ...Node) Node {
	return &htmlElement{
		Name:     name,
		Children: children,
	}
}

func CreateElementVoid(name string, children ...Node) Node {
	return &htmlElement{
		Name:     name,
		IsVoid:   true,
		Children: children,
	}
}

func CreateAttr(name, value string) Node {
	return &htmlAttr{
		Name:  name,
		Value: value,
	}
}

func CreateAttrBoolean(name string) Node {
	return &htmlAttr{
		Name:  name,
		Value: "",
	}
}

func CreateString(value string) Node {
	return htmlSafe(value)
}

func CreateStringRaw(value string) Node {
	return htmlRaw(value)
}

// ------------------------------------------------------------------
//
// Helper functions
//
// ------------------------------------------------------------------

// Common byte slices for HTML rendering to reduce allocations.
var (
	bLt          = []byte("<")
	bGt          = []byte(">")
	bLtSlash     = []byte("</")
	bSpace       = []byte(" ")
	bEqualsQuote = []byte(`="`)
	bQuote       = []byte(`"`)
)

// needsEscaping checks if a string contains characters that require HTML escaping.
// This is based on the characters handled by html/template.HTMLEscapeString.
func needsEscaping(s string) bool {
	return strings.ContainsAny(s, "'\"&<>\000")
}
