package dom

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/tunedmystic/rio/internal/assert"
)

func TestHtmlElement(t *testing.T) {
	t.Run("CreateElement/Render/String", func(t *testing.T) {
		r := CreateElement("div", &htmlAttr{Name: "class", Value: "test"})
		assert.Equal(t, render(r), `<div class="test"></div>`)
		assert.Equal(t, fmt.Sprint(r), `<div class="test"></div>`)
	})

	t.Run("CreateElementVoid/Render/String", func(t *testing.T) {
		r := CreateElementVoid("img", &htmlAttr{Name: "src", Value: "image.jpg"})
		assert.Equal(t, render(r), `<img src="image.jpg">`)
		assert.Equal(t, fmt.Sprint(r), `<img src="image.jpg">`)
	})

	t.Run("Render error on attribute", func(t *testing.T) {
		r := Div(errorAttr{}, Span(Text("test")))
		// w := ErrorWriter{}
		var b bytes.Buffer
		err := r.Render(&b)
		assert.Equal(t, err.Error(), "error from errorAttr.RenderAttribute")
	})

	t.Run("Render error on children", func(t *testing.T) {
		r := Div(Class("test"), errorNode{})
		var b bytes.Buffer
		err := r.Render(&b)
		assert.Equal(t, err.Error(), "error from errorNode.Render")
	})
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func TestElements(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		tests := []struct {
			Name     string
			ElemFunc func(...Node) Node
		}{
			{Name: "a", ElemFunc: A},
			{Name: "abbr", ElemFunc: Abbr},
			{Name: "address", ElemFunc: Address},
			{Name: "article", ElemFunc: Article},
			{Name: "aside", ElemFunc: Aside},
			{Name: "audio", ElemFunc: Audio},
			{Name: "b", ElemFunc: B},
			{Name: "bdi", ElemFunc: Bdi},
			{Name: "bdo", ElemFunc: Bdo},
			{Name: "blockquote", ElemFunc: Blockquote},
			{Name: "body", ElemFunc: Body},
			{Name: "button", ElemFunc: Button},
			{Name: "canvas", ElemFunc: Canvas},
			{Name: "caption", ElemFunc: Caption},
			{Name: "cite", ElemFunc: Cite},
			{Name: "code", ElemFunc: Code},
			{Name: "colgroup", ElemFunc: Colgroup},
			{Name: "data", ElemFunc: DataEl},
			{Name: "datalist", ElemFunc: Datalist},
			{Name: "dd", ElemFunc: Dd},
			{Name: "del", ElemFunc: Del},
			{Name: "details", ElemFunc: Details},
			{Name: "dfn", ElemFunc: Dfn},
			{Name: "dialog", ElemFunc: Dialog},
			{Name: "div", ElemFunc: Div},
			{Name: "dl", ElemFunc: Dl},
			{Name: "dt", ElemFunc: Dt},
			{Name: "em", ElemFunc: Em},
			{Name: "fieldset", ElemFunc: Fieldset},
			{Name: "figcaption", ElemFunc: Figcaption},
			{Name: "figure", ElemFunc: Figure},
			{Name: "footer", ElemFunc: Footer},
			{Name: "form", ElemFunc: Form},
			{Name: "h1", ElemFunc: H1},
			{Name: "h2", ElemFunc: H2},
			{Name: "h3", ElemFunc: H3},
			{Name: "h4", ElemFunc: H4},
			{Name: "h5", ElemFunc: H5},
			{Name: "h6", ElemFunc: H6},
			{Name: "head", ElemFunc: Head},
			{Name: "header", ElemFunc: Header},
			{Name: "html", ElemFunc: Html},
			{Name: "i", ElemFunc: I},
			{Name: "ins", ElemFunc: Ins},
			{Name: "kbd", ElemFunc: Kbd},
			{Name: "label", ElemFunc: Label},
			{Name: "legend", ElemFunc: Legend},
			{Name: "li", ElemFunc: Li},
			{Name: "main", ElemFunc: Main},
			{Name: "map", ElemFunc: MapEl},
			{Name: "mark", ElemFunc: Mark},
			{Name: "menu", ElemFunc: Menu},
			{Name: "meter", ElemFunc: Meter},
			{Name: "nav", ElemFunc: Nav},
			{Name: "noscript", ElemFunc: Noscript},
			{Name: "object", ElemFunc: Object},
			{Name: "ol", ElemFunc: Ol},
			{Name: "optgroup", ElemFunc: Optgroup},
			{Name: "option", ElemFunc: Option},
			{Name: "output", ElemFunc: Output},
			{Name: "p", ElemFunc: P},
			{Name: "picture", ElemFunc: Picture},
			{Name: "pre", ElemFunc: Pre},
			{Name: "progress", ElemFunc: Progress},
			{Name: "q", ElemFunc: Q},
			{Name: "s", ElemFunc: S},
			{Name: "samp", ElemFunc: Samp},
			{Name: "script", ElemFunc: Script},
			{Name: "section", ElemFunc: Section},
			{Name: "select", ElemFunc: Select},
			{Name: "slot", ElemFunc: SlotEl},
			{Name: "small", ElemFunc: Small},
			{Name: "span", ElemFunc: Span},
			{Name: "strong", ElemFunc: Strong},
			{Name: "style", ElemFunc: StyleEl},
			{Name: "sub", ElemFunc: Sub},
			{Name: "summary", ElemFunc: Summary},
			{Name: "sup", ElemFunc: Sup},
			{Name: "table", ElemFunc: Table},
			{Name: "tbody", ElemFunc: Tbody},
			{Name: "td", ElemFunc: Td},
			{Name: "template", ElemFunc: Template},
			{Name: "textarea", ElemFunc: Textarea},
			{Name: "tfoot", ElemFunc: Tfoot},
			{Name: "th", ElemFunc: Th},
			{Name: "thead", ElemFunc: Thead},
			{Name: "time", ElemFunc: Time},
			{Name: "title", ElemFunc: TitleEl},
			{Name: "tr", ElemFunc: Tr},
			{Name: "u", ElemFunc: U},
			{Name: "ul", ElemFunc: Ul},
			{Name: "var", ElemFunc: Var},
			{Name: "video", ElemFunc: Video},
			{Name: "wbr", ElemFunc: Wbr},
		}

		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				r := test.ElemFunc(CreateAttr("name", "test"))
				assert.Equal(t, render(r), fmt.Sprintf(`<%s name="test"></%s>`, test.Name, test.Name))
			})
		}
	})

	t.Run("void", func(t *testing.T) {
		tests := []struct {
			Name     string
			ElemFunc func(...Node) Node
		}{
			{Name: "area", ElemFunc: Area},
			{Name: "base", ElemFunc: Base},
			{Name: "br", ElemFunc: Br},
			{Name: "col", ElemFunc: Col},
			{Name: "embed", ElemFunc: Embed},
			{Name: "hr", ElemFunc: Hr},
			{Name: "iframe", ElemFunc: Iframe},
			{Name: "img", ElemFunc: Img},
			{Name: "input", ElemFunc: Input},
			{Name: "link", ElemFunc: Link},
			{Name: "meta", ElemFunc: Meta},
			{Name: "param", ElemFunc: Param},
			{Name: "source", ElemFunc: Source},
			{Name: "track", ElemFunc: Track},
		}

		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				r := test.ElemFunc(CreateAttr("name", "test"))
				assert.Equal(t, render(r), fmt.Sprintf(`<%s name="test">`, test.Name))
			})
		}
	})

	t.Run("doctype", func(t *testing.T) {
		r := Doctype(Html(Lang("en"), Head(Meta(Charset("UTF-8")))))
		assert.Equal(t, render(r), `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"></head></html>`)
	})

	t.Run("text", func(t *testing.T) {
		r := Text("test < testing")
		assert.Equal(t, render(r), "test &lt; testing")
	})

	t.Run("raw", func(t *testing.T) {
		r := Raw("test < testing")
		assert.Equal(t, render(r), "test < testing")
	})
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Benchmark_Elements(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = Div(
			Class("container"),
			Id("main-content"),
			P(Text("Hello, World!")),
			Span(Text("Another element")),
		)
	}
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Benchmark_ElementsDeeplyNested(b *testing.B) {
	b.ReportAllocs()

	const depth = 100
	var node Node = Div()

	for j := 0; j < depth; j++ {
		// Create a new Div that wraps the previous node
		node = Div(Class("inner"), node)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = node.Render(io.Discard)
	}
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Benchmark_ElementWithManyAttributes(b *testing.B) {
	b.ReportAllocs()

	const numAttrs = 100
	attrs := make([]Node, numAttrs)

	for j := 0; j < numAttrs; j++ {
		attrs[j] = CreateAttr(fmt.Sprintf("data-attr%d", j), "value")
	}
	node := Div(attrs...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = node.Render(io.Discard)
	}
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Benchmark_StringEscaped(b *testing.B) {
	b.ReportAllocs()

	// A string that contains characters requiring HTML escaping
	longString := strings.Repeat("Text with <html> tags & special chars like < > & \" ' ", 50)
	node := Text(longString)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = node.Render(io.Discard)
	}
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Benchmark_StringRaw(b *testing.B) {
	b.ReportAllocs()

	longString := strings.Repeat("Text with <html> tags & special chars like < > & \" ' ", 50)
	node := Raw(longString)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = node.Render(io.Discard)
	}
}
