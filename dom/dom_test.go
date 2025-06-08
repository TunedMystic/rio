package dom

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tunedmystic/rio/internal/assert"
)

func render(n Node) string {
	var b strings.Builder
	n.Render(&b)
	return b.String()
}

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

func TestAttributes(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		tests := []struct {
			Name     string
			AttrFunc func(string) Node
		}{
			{Name: "accesskey", AttrFunc: Accesskey},
			{Name: "accept", AttrFunc: Accept},
			{Name: "action", AttrFunc: Action},
			{Name: "alt", AttrFunc: Alt},
			{Name: "as", AttrFunc: As},
			{Name: "autocapitalize", AttrFunc: Autocapitalize},
			{Name: "autocomplete", AttrFunc: Autocomplete},
			{Name: "charset", AttrFunc: Charset},
			{Name: "cite", AttrFunc: CiteAttr},
			{Name: "class", AttrFunc: Class},
			{Name: "cols", AttrFunc: Cols},
			{Name: "colspan", AttrFunc: Colspan},
			{Name: "content", AttrFunc: Content},
			{Name: "contenteditable", AttrFunc: Contenteditable},
			{Name: "crossorigin", AttrFunc: Crossorigin},
			{Name: "datetime", AttrFunc: Datetime},
			{Name: "dir", AttrFunc: Dir},
			{Name: "download", AttrFunc: Download},
			{Name: "draggable", AttrFunc: Draggable},
			{Name: "enctype", AttrFunc: Enctype},
			{Name: "enterkeyhint", AttrFunc: Enterkeyhint},
			{Name: "for", AttrFunc: For},
			{Name: "form", AttrFunc: FormAttr},
			{Name: "headers", AttrFunc: Headers},
			{Name: "height", AttrFunc: Height},
			{Name: "href", AttrFunc: Href},
			{Name: "hreflang", AttrFunc: Hreflang},
			{Name: "http-equiv", AttrFunc: Httpequiv},
			{Name: "id", AttrFunc: Id},
			{Name: "inert", AttrFunc: Inert},
			{Name: "inputmode", AttrFunc: Inputmode},
			{Name: "integrity", AttrFunc: Integrity},
			{Name: "itemid", AttrFunc: Itemid},
			{Name: "itemprop", AttrFunc: Itemprop},
			{Name: "itemref", AttrFunc: Itemref},
			{Name: "itemtype", AttrFunc: Itemtype},
			{Name: "lang", AttrFunc: Lang},
			{Name: "list", AttrFunc: List},
			{Name: "loading", AttrFunc: Loading},
			{Name: "max", AttrFunc: Max},
			{Name: "maxlength", AttrFunc: Maxlength},
			{Name: "media", AttrFunc: Media},
			{Name: "method", AttrFunc: Method},
			{Name: "min", AttrFunc: Min},
			{Name: "minlength", AttrFunc: Minlength},
			{Name: "name", AttrFunc: Name},
			{Name: "nonce", AttrFunc: Nonce},
			{Name: "onload", AttrFunc: Onload},
			{Name: "part", AttrFunc: Part},
			{Name: "pattern", AttrFunc: Pattern},
			{Name: "placeholder", AttrFunc: Placeholder},
			{Name: "popover", AttrFunc: Popover},
			{Name: "poster", AttrFunc: Poster},
			{Name: "preload", AttrFunc: Preload},
			{Name: "property", AttrFunc: Property},
			{Name: "referrerpolicy", AttrFunc: Referrerpolicy},
			{Name: "rel", AttrFunc: Rel},
			{Name: "role", AttrFunc: Role},
			{Name: "rows", AttrFunc: Rows},
			{Name: "rowspan", AttrFunc: Rowspan},
			{Name: "sandbox", AttrFunc: Sandbox},
			{Name: "scope", AttrFunc: Scope},
			{Name: "size", AttrFunc: Size},
			{Name: "sizes", AttrFunc: Sizes},
			{Name: "slot", AttrFunc: Slot},
			{Name: "spellcheck", AttrFunc: Spellcheck},
			{Name: "src", AttrFunc: Src},
			{Name: "srcset", AttrFunc: Srcset},
			{Name: "step", AttrFunc: Step},
			{Name: "style", AttrFunc: Style},
			{Name: "tabindex", AttrFunc: Tabindex},
			{Name: "target", AttrFunc: Target},
			{Name: "title", AttrFunc: Title},
			{Name: "translate", AttrFunc: Translate},
			{Name: "type", AttrFunc: Type},
			{Name: "value", AttrFunc: Value},
			{Name: "width", AttrFunc: Width},
		}

		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				r := Div(test.AttrFunc("test"))
				assert.Equal(t, render(r), fmt.Sprintf(`<div %s="test"></div>`, test.Name))
			})
		}
	})

	t.Run("boolean", func(t *testing.T) {
		tests := []struct {
			Name     string
			AttrFunc func() Node
		}{
			{Name: "allowfullscreen", AttrFunc: Allowfullscreen},
			{Name: "async", AttrFunc: Async},
			{Name: "autofocus", AttrFunc: Autofocus},
			{Name: "autoplay", AttrFunc: Autoplay},
			{Name: "checked", AttrFunc: Checked},
			{Name: "controls", AttrFunc: Controls},
			{Name: "defer", AttrFunc: Defer},
			{Name: "disabled", AttrFunc: Disabled},
			{Name: "formnovalidate", AttrFunc: Formnovalidate},
			{Name: "hidden", AttrFunc: Hidden},
			{Name: "itemscope", AttrFunc: Itemscope},
			{Name: "loop", AttrFunc: Loop},
			{Name: "multiple", AttrFunc: Multiple},
			{Name: "muted", AttrFunc: Muted},
			{Name: "nomodule", AttrFunc: Nomodule},
			{Name: "novalidate", AttrFunc: Novalidate},
			{Name: "open", AttrFunc: Open},
			{Name: "readonly", AttrFunc: Readonly},
			{Name: "required", AttrFunc: Required},
			{Name: "selected", AttrFunc: Selected},
		}

		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				r := Div(test.AttrFunc())
				assert.Equal(t, render(r), fmt.Sprintf(`<div %s></div>`, test.Name))
			})
		}
	})

	t.Run("data attr", func(t *testing.T) {
		r := Div(Data("test", "value"))
		assert.Equal(t, render(r), `<div data-test="value"></div>`)
	})

	t.Run("aria attr", func(t *testing.T) {
		r := Div(Aria("label", "test"))
		assert.Equal(t, render(r), `<div aria-label="test"></div>`)
	})
}

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
		r := Div(ErrorAttr{}, Span(Text("test")))
		// w := ErrorWriter{}
		var b bytes.Buffer
		err := r.Render(&b)
		assert.Equal(t, err.Error(), "error rendering attribute")
	})

	t.Run("Render error on children", func(t *testing.T) {
		r := Div(Class("test"), ErrNode{})
		var b bytes.Buffer
		err := r.Render(&b)
		assert.Equal(t, err.Error(), "error rendering node")
	})
}

type ErrorAttr struct{}

func (a ErrorAttr) Render(w io.Writer) error {
	return errors.New("error rendering attribute")
}

type ErrNode struct{}

func (e ErrNode) Render(w io.Writer) error {
	return errors.New("error rendering node")
}

func TestHtmlAttribute(t *testing.T) {
	t.Run("CreateAttr/Render/String", func(t *testing.T) {
		r := CreateAttr("name", "test")
		assert.Equal(t, render(r), ` name="test"`)
		assert.Equal(t, fmt.Sprint(r), ` name="test"`)
	})

	t.Run("CreateAttrBoolean", func(t *testing.T) {
		r := CreateAttrBoolean("hidden")
		assert.Equal(t, render(r), ` hidden`)
		assert.Equal(t, fmt.Sprint(r), ` hidden`)
	})
}

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

func TestHandlerMiddleware(t *testing.T) {
	t.Run("DomHandler", func(t *testing.T) {
		fn := Handler(func(w http.ResponseWriter, r *http.Request) Node {
			return Div(Text("Hello, World!"))
		})

		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		fn.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusOK)
		assert.Equal(t, rr.Body.String(), `<div>Hello, World!</div>`)
	})

	t.Run("DomHandler with nil Node", func(t *testing.T) {
		fn := Handler(func(w http.ResponseWriter, r *http.Request) Node {
			return nil
		})

		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		fn.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusInternalServerError)
		assert.Equal(t, rr.Body.String(), "Internal Server Error\n")
	})

	t.Run("DomHandler with error", func(t *testing.T) {
		fn := Handler(func(w http.ResponseWriter, r *http.Request) Node {
			return ErrNode{}
		})

		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		fn.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusInternalServerError)
		assert.Equal(t, rr.Body.String(), "Internal Server Error\n")
	})
}

func BenchmarkLargeHTMLDocument(b *testing.B) {
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

		var sb strings.Builder
		_ = doc.Render(&sb)
	}
}

func BenchmarkElementCreation(b *testing.B) {
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

func BenchmarkAttributeCreation(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Class("my-class")
		_ = Id("my-id")
		_ = Href("/some/url")
		_ = Src("image.png")
	}
}

func BenchmarkRenderDeeplyNestedElements(b *testing.B) {
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

func BenchmarkRenderElementWithManyAttributes(b *testing.B) {
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

func BenchmarkRenderHtmlStringEscaped(b *testing.B) {
	b.ReportAllocs()
	// A string that contains characters requiring HTML escaping
	longString := strings.Repeat("Text with <html> tags & special chars like < > & \" ' ", 50)
	node := Text(longString)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = node.Render(io.Discard)
	}
}

func BenchmarkRenderHtmlStringRaw(b *testing.B) {
	b.ReportAllocs()
	longString := strings.Repeat("Text with <html> tags & special chars like < > & \" ' ", 50)
	node := Raw(longString)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = node.Render(io.Discard)
	}
}

func BenchmarkMapAndRender(b *testing.B) {
	b.ReportAllocs()
	const numItems = 500
	items := make([]string, numItems)
	for j := 0; j < numItems; j++ {
		items[j] = fmt.Sprintf("item-%d", j)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Map items to Nodes
		groupNode := Map(items, func(s string) Node {
			return Li(Text(s))
		})
		// Render the resulting group
		_ = groupNode.Render(io.Discard)
	}
}

func BenchmarkGroupRenderLarge(b *testing.B) {
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
