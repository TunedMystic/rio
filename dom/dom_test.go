package dom

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/tunedmystic/rio/internal/assert"
)

func TestHtmlElement_Render_Errors(t *testing.T) {
	errWriterSentinel := errors.New("writer error")
	// Use error messages that match the modified errorAttr and existing errorNode
	errAttrRenderSentinel := errors.New("error from errorAttr.RenderAttribute")
	errChildRenderSentinel := errors.New("error from errorNode.Render")

	tests := []struct {
		name        string
		element     Node
		writer      *errorWriter // nil if not testing writer failure
		expectedErr error
	}{
		{
			name:        "fail write opening <",
			element:     Div(),
			writer:      &errorWriter{targetErr: errWriterSentinel, failOnNthWrite: 1},
			expectedErr: errWriterSentinel,
		},
		{
			name:        "fail write element name",
			element:     Div(),
			writer:      &errorWriter{targetErr: errWriterSentinel, failOnNthWrite: 2},
			expectedErr: errWriterSentinel,
		},
		{
			name:        "fail attribute's RenderAttribute method",
			element:     Div(errorAttr{}),
			writer:      nil,
			expectedErr: errAttrRenderSentinel,
		},
		{
			name:        "fail write > after attributes",
			element:     Div(Class("foo")),                                                         // Class("foo") makes 5 writes
			writer:      &errorWriter{targetErr: errWriterSentinel, failOnNthWrite: 1 + 1 + 5 + 1}, // <, div, (attr: " ", class, =, "foo", "), >
			expectedErr: errWriterSentinel,
		},
		{
			name:        "fail write > for void element after attributes",
			element:     Img(Src("test.jpg")),                                                      // Img is void, Src attribute makes 5 writes
			writer:      &errorWriter{targetErr: errWriterSentinel, failOnNthWrite: 1 + 1 + 5 + 1}, // <, img, (attr: " ", src, =, "test.jpg", "), >
			expectedErr: errWriterSentinel,
		},
		{
			name:        "fail child's Render method",
			element:     Div(errorNode{}),
			writer:      nil,
			expectedErr: errChildRenderSentinel,
		},
		{
			name:        "fail write </ for closing tag",
			element:     Div(Text("hi")),                                                               // Text("hi") makes 1 write
			writer:      &errorWriter{targetErr: errWriterSentinel, failOnNthWrite: 1 + 1 + 1 + 1 + 1}, // <, div, >, "hi", </
			expectedErr: errWriterSentinel,
		},
		{
			name:        "fail write closing element name",
			element:     Div(Text("hi")),
			writer:      &errorWriter{targetErr: errWriterSentinel, failOnNthWrite: 1 + 1 + 1 + 1 + 1 + 1}, // <, div, >, "hi", </, div
			expectedErr: errWriterSentinel,
		},
		{
			name:        "fail write final > for closing tag",
			element:     Div(Text("hi")),
			writer:      &errorWriter{targetErr: errWriterSentinel, failOnNthWrite: 1 + 1 + 1 + 1 + 1 + 1 + 1}, // <, div, >, "hi", </, div, >
			expectedErr: errWriterSentinel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w io.Writer = io.Discard
			if tt.writer != nil {
				w = tt.writer
			}
			err := tt.element.Render(w)
			assert.Equal(t, err.Error(), tt.expectedErr.Error())
		})
	}
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func TestHtmlString(t *testing.T) {
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

func TestControlStructures(t *testing.T) {
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
		_ = doc.Render(io.Discard)
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
		// Map items to Nodes
		group := Map(items, func(s string) Node {
			return Li(Text(s))
		})
		// Render the group
		_ = group.Render(io.Discard)
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
		_ = group.Render(io.Discard)
	}
}

// ------------------------------------------------------------------
//
// Test Helpers
//
// ------------------------------------------------------------------

func render(n Node) string {
	var b strings.Builder
	n.Render(&b)
	return b.String()
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
	failOnNthWrite int // 1-based index of Write call to fail on
	writeCount     int
	// buf can be used to inspect what was written before the error, if necessary.
	// For these tests, we primarily care about the error propagation.
	// buf bytes.Buffer
}

func (ew *errorWriter) Write(p []byte) (n int, err error) {
	ew.writeCount++
	if ew.writeCount == ew.failOnNthWrite {
		return 0, ew.targetErr
	}
	// ew.buf.Write(p)
	return len(p), nil // Simulate successful write for non-failing calls
}

var _ io.Writer = (*errorWriter)(nil)
