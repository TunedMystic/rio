package dom

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/tunedmystic/rio/internal/assert"
)

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Test_Group_Render_Errors(t *testing.T) {
	errNodeRenderSentinel := errors.New("error from errorNode.Render")
	errWriterSentinel := errors.New("writer error from group test")

	tests := []struct {
		name        string
		group       Group
		writer      *errorWriter
		expectedErr error
	}{
		{
			name:        "empty group",
			group:       Group{},
			writer:      nil,
			expectedErr: nil,
		},
		{
			name:        "single failing node",
			group:       Group{errorNode{}},
			writer:      nil,
			expectedErr: errNodeRenderSentinel,
		},
		{
			name:        "first node fails in a multi-node group",
			group:       Group{errorNode{}, Text("this should not render")},
			writer:      nil,
			expectedErr: errNodeRenderSentinel,
		},
		{
			name:        "middle node fails in a multi-node group",
			group:       Group{Text("first"), errorNode{}, Text("this should not render")},
			writer:      nil,
			expectedErr: errNodeRenderSentinel,
		},
		{
			name:        "writer error during rendering of a child node",
			group:       Group{Text("try to write this")},
			writer:      &errorWriter{targetErr: errWriterSentinel, failOnNthWrite: 1},
			expectedErr: errWriterSentinel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w io.Writer = io.Discard
			if tt.writer != nil {
				w = tt.writer
			}
			err := tt.group.Render(w)
			assert.Error(t, err, tt.expectedErr)
		})
	}
}

func Test_NodeMapper_Render_Errors(t *testing.T) {
	errNodeRenderSentinel := errors.New("error from errorNode.Render")
	errWriterSentinel := errors.New("writer error from nodeMapper test")

	tests := []struct {
		name        string
		node        Node
		writer      *errorWriter
		expectedErr error
	}{
		{
			name:        "empty items",
			node:        Map([]string{}, func(s string) Node { return Text(s) }),
			writer:      nil,
			expectedErr: nil,
		},
		{
			name:        "single failing node",
			node:        Map([]string{"fail"}, func(s string) Node { return errorNode{} }),
			writer:      nil,
			expectedErr: errNodeRenderSentinel,
		},
		{
			name: "first node fails in a multi-item map",
			node: Map([]string{"fail", "ok"}, func(s string) Node {
				if s == "fail" {
					return errorNode{}
				}
				return Text(s)
			}),
			writer:      nil,
			expectedErr: errNodeRenderSentinel,
		},
		{
			name: "middle node fails in a multi-item map",
			node: Map([]string{"ok", "fail", "another_ok"}, func(s string) Node {
				if s == "fail" {
					return errorNode{}
				}
				return Text(s)
			}),
			writer:      nil,
			expectedErr: errNodeRenderSentinel,
		},
		{
			name:        "writer error during rendering of a mapped node",
			node:        Map([]string{"write_me"}, func(s string) Node { return Text(s) }),
			writer:      &errorWriter{targetErr: errWriterSentinel, failOnNthWrite: 1},
			expectedErr: errWriterSentinel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w io.Writer = io.Discard
			if tt.writer != nil {
				w = tt.writer
			}
			err := tt.node.Render(w)
			assert.Error(t, err, tt.expectedErr)
		})
	}
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Test_HtmlString(t *testing.T) {
	t.Run("CreateString/Render/String", func(t *testing.T) {
		r := CreateString("test < testing")
		assert.Equal(t, render(r), "test &lt; testing")
		assert.Equal(t, fmt.Sprint(r), "test &lt; testing")
	})

	t.Run("CreateStringRaw/Render/String", func(t *testing.T) {
		r := CreateStringRaw("test < testing")
		assert.Equal(t, render(r), "test < testing")
		assert.Equal(t, fmt.Sprint(r), "test < testing")
	})
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Test_Doctype(t *testing.T) {
	t.Run("Render/String", func(t *testing.T) {
		r := Doctype(Html(Lang("en"), Head(Meta(Charset("UTF-8")))))
		assert.Equal(t, render(r), `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"></head></html>`)
		assert.Equal(t, fmt.Sprint(r), `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"></head></html>`)
	})

	t.Run("Render error", func(t *testing.T) {
		r := Doctype(Html(Lang("en"), Head(Meta(Charset("UTF-8")))))
		errWriterSentinel := errors.New("writer error from htmlDoctype test")
		writer := &errorWriter{
			targetErr:      errWriterSentinel,
			failOnNthWrite: 1,
		}

		err := r.Render(writer)
		assert.Error(t, err, errors.New("writer error from htmlDoctype test"))
	})
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Test_ControlStructures(t *testing.T) {
	t.Run("Group", func(t *testing.T) {
		r := Group{
			Div(Text("foo")),
			Div(Text("bar")),
		}
		assert.Equal(t, render(r), `<div>foo</div><div>bar</div>`)
		assert.Equal(t, fmt.Sprint(r), `<div>foo</div><div>bar</div>`)
	})

	t.Run("Map", func(t *testing.T) {
		r := Map(
			[]string{"foo", "bar"},
			func(s string) Node {
				return Div(Text(s))
			},
		)
		assert.Equal(t, render(r), `<div>foo</div><div>bar</div>`)
		assert.Equal(t, fmt.Sprint(r), `<div>foo</div><div>bar</div>`)
	})

	t.Run("If", func(t *testing.T) {
		r1 := If(true, Div(Text("foo")))
		assert.Equal(t, render(r1), `<div>foo</div>`)

		r2 := If(false, Div(Text("foo")))
		assert.Equal(t, render(r2), ``)
	})

	t.Run("Ifelse", func(t *testing.T) {
		r1 := Ifelse(true, Div(Text("foo")), Div(Text("bar")))
		assert.Equal(t, render(r1), `<div>foo</div>`)

		r2 := Ifelse(false, Div(Text("foo")), Div(Text("bar")))
		assert.Equal(t, render(r2), `<div>bar</div>`)
	})

}

// ------------------------------------------------------------------
//
// Test Nil Checks
//
// ------------------------------------------------------------------

func Test_Rendering_NilChecks(t *testing.T) {
	t.Run("htmlElement attribute is (*htmlAttr)(nil)", func(t *testing.T) {
		var nilConcreteAttribute Node = (*htmlAttr)(nil)
		element := Div(Class("test"), nilConcreteAttribute, Id("my-id"))

		var sb strings.Builder
		err := element.Render(&sb)
		assert.Error(t, err, nil)
		assert.Equal(t, sb.String(), `<div class="test" id="my-id"></div>`)
	})

	t.Run("htmlElement child node is (Node)(nil)", func(t *testing.T) {
		var nilInterfaceNode Node = nil
		element := Div(Text("hello"), nilInterfaceNode, Text("world"))

		var sb strings.Builder
		err := element.Render(&sb)
		assert.Error(t, err, nil)
		assert.Equal(t, sb.String(), `<div>helloworld</div>`)
	})

	t.Run("Group contains (Node)(nil)", func(t *testing.T) {
		var nilInterfaceNode Node = nil
		group := Group{
			Text("first"),
			nilInterfaceNode,
			Text("third"),
		}

		var sb strings.Builder
		err := group.Render(&sb)
		assert.Error(t, err, nil)
		assert.Equal(t, sb.String(), `firstthird`)
	})

	t.Run("nodeMapper function returns (Node)(nil)", func(t *testing.T) {
		items := []string{"one", "two", "three"}
		mapper := Map(items, func(item string) Node {
			if item == "two" {
				return nil
			}
			return Text(item)
		})

		var sb strings.Builder
		err := mapper.Render(&sb)
		assert.Error(t, err, nil)
		assert.Equal(t, sb.String(), `onethree`)
	})

	t.Run("htmlDoctype sibling is (Node)(nil)", func(t *testing.T) {
		var nilSiblingNode Node = nil
		doctypeNode := Doctype(nilSiblingNode)

		var sb strings.Builder
		err := doctypeNode.Render(&sb)
		assert.Error(t, err, nil)
		assert.Equal(t, sb.String(), `<!DOCTYPE html>`)
	})
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Benchmark_Document(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		elements := make([]Node, 0, 10000)

		for i := 0; i < 5000; i++ {
			elements = append(elements,
				Div(Class("foo")),
				Span(Class("bar")),
			)
		}
		doc := Div(elements...)
		renderNull(doc)
	}
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Benchmark_Map(b *testing.B) {
	b.ReportAllocs()

	const numItems = 500
	items := make([]string, numItems)

	for j := 0; j < numItems; j++ {
		items[j] = fmt.Sprintf("item-%d", j)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		group := Map(items, func(s string) Node {
			return Li(Text(s))
		})
		renderNull(group)
	}
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Benchmark_Group(b *testing.B) {
	b.ReportAllocs()

	const numNodes = 1000
	nodes := make([]Node, numNodes)

	for j := 0; j < numNodes; j++ {
		nodes[j] = Span(Text(fmt.Sprintf("span-node-%d", j)))
	}
	group := Group(nodes)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderNull(group)
	}
}

// ------------------------------------------------------------------
//
// Test Helpers
//
// ------------------------------------------------------------------

func render(n Node) string {
	var b strings.Builder
	_ = n.Render(&b)
	return b.String()
}

func renderNull(n Node) {
	_ = n.Render(io.Discard)
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

// errorAttr is a mock attribute that always returns an error on rendering.
type errorAttr struct{}

func (a errorAttr) Render(w io.Writer) error {
	return errors.New("error from errorAttr.Render")
}

func (a errorAttr) RenderAttribute(w io.Writer) error {
	return errors.New("error from errorAttr.RenderAttribute")
}

var _ Node = (*errorAttr)(nil)
var _ HtmlAttributer = (*errorAttr)(nil)

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

// errorNode is a mock Node that always returns an error on rendering.
type errorNode struct{}

func (e errorNode) Render(w io.Writer) error {
	return errors.New("error from errorNode.Render")
}

var _ Node = (*errorNode)(nil)

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

// errorWriter is a mock io.Writer that fails on the Nth call to Write.
type errorWriter struct {
	targetErr      error
	failOnNthWrite int
	writeCount     int
}

func (ew *errorWriter) Write(p []byte) (n int, err error) {
	ew.writeCount++
	if ew.writeCount == ew.failOnNthWrite {
		return 0, ew.targetErr
	}
	return len(p), nil
}

var _ io.Writer = (*errorWriter)(nil)
