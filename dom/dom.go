package dom

import (
	"html"
	"io"
	"strings"
)

type Node interface {
	Render(w io.Writer) error
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

// htmlElement represents an HTML element
type htmlElement struct {
	Name     string
	IsVoid   bool
	Children []Node
}

func (e htmlElement) Render(w io.Writer) error {
	// Render opening tag
	w.Write([]byte("<"))
	w.Write([]byte(e.Name))

	// Render attributes
	for _, c := range e.Children {
		if attr, ok := c.(htmlAttr); ok {
			if err := attr.Render(w); err != nil {
				return err
			}
		}
	}

	w.Write([]byte(">"))

	// Void elements have no children or closing tags.
	if e.IsVoid {
		return nil
	}

	// Render children
	for _, c := range e.Children {
		if _, ok := c.(htmlAttr); !ok {
			if err := c.Render(w); err != nil {
				return err
			}
		}
	}

	// Render closing tag
	w.Write([]byte("</"))
	w.Write([]byte(e.Name))
	w.Write([]byte(">"))

	return nil
}

func (e htmlElement) String() string {
	var b strings.Builder
	e.Render(&b)
	return b.String()
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

// htmlAttr represents an HTML attribute
type htmlAttr struct {
	Name  string
	Value string
}

func (a htmlAttr) Render(w io.Writer) error {
	w.Write([]byte(" "))
	w.Write([]byte(a.Name))

	if a.Value != "" {
		w.Write([]byte(`="`))
		w.Write([]byte(html.EscapeString(a.Value)))
		w.Write([]byte(`"`))
	}
	return nil
}

func (a htmlAttr) String() string {
	var b strings.Builder
	a.Render(&b)
	return b.String()
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

// htmlString represents a piece of HTML content
type htmlString struct {
	Value string
	Raw   bool
}

func (s htmlString) Render(w io.Writer) error {
	if s.Raw {
		w.Write([]byte(s.Value))
	} else {
		w.Write([]byte(html.EscapeString(s.Value)))
	}
	return nil
}

func (s htmlString) String() string {
	var b strings.Builder
	s.Render(&b)
	return b.String()
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

// Group is a convenience type for grouping multiple Nodes together.
type Group []Node

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

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

// Map is a utility function that maps a slice of items to a Group of Nodes.
func Map[T any](items []T, fn func(T) Node) Group {
	nodes := make([]Node, 0, len(items))
	for _, item := range items {
		nodes = append(nodes, fn(item))
	}
	return nodes
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

// If is a utility function that returns a Node if the condition is true,
// otherwise returns nil.
func If(condition bool, a Node) Node {
	if condition {
		return a
	}
	return nil
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
//
//
// ------------------------------------------------------------------

func CreateElement(name string, children ...Node) Node {
	return htmlElement{
		Name:     name,
		Children: children,
	}
}

func CreateElementVoid(name string, children ...Node) Node {
	return htmlElement{
		Name:     name,
		IsVoid:   true,
		Children: children,
	}
}

func CreateAttr(name, value string) Node {
	return htmlAttr{
		Name:  name,
		Value: value,
	}
}

func CreateAttrBoolean(name string) Node {
	return htmlAttr{
		Name:  name,
		Value: "",
	}
}

func CreateString(value string) Node {
	return htmlString{
		Value: value,
	}
}

func CreateStringRaw(value string) Node {
	return htmlString{
		Value: value,
		Raw:   true,
	}
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Text(value string) Node {
	return CreateString(value)
}

func Raw(value string) Node {
	return CreateStringRaw(value)
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Doctype(sibling Node) Node {
	return Group([]Node{
		Raw("<!DOCTYPE html>"),
		sibling,
	})
}

// ------------------------------------------------------------------
//
// Dom Elements
//
// ------------------------------------------------------------------

func A(children ...Node) Node {
	return CreateElement("a", children...)
}

func Abbr(children ...Node) Node {
	return CreateElement("abbr", children...)
}

func Address(children ...Node) Node {
	return CreateElement("address", children...)
}

func Area(children ...Node) Node {
	return CreateElementVoid("area", children...)
}

func Article(children ...Node) Node {
	return CreateElement("article", children...)
}

func Aside(children ...Node) Node {
	return CreateElement("aside", children...)
}

func Audio(children ...Node) Node {
	return CreateElement("audio", children...)
}

func B(children ...Node) Node {
	return CreateElement("b", children...)
}

func Base(children ...Node) Node {
	return CreateElementVoid("base", children...)
}

func Bdi(children ...Node) Node {
	return CreateElement("bdi", children...)
}

func Bdo(children ...Node) Node {
	return CreateElement("bdo", children...)
}

func Blockquote(children ...Node) Node {
	return CreateElement("blockquote", children...)
}

func Body(children ...Node) Node {
	return CreateElement("body", children...)
}

func Br(children ...Node) Node {
	return CreateElementVoid("br", children...)
}

func Button(children ...Node) Node {
	return CreateElement("button", children...)
}

func Canvas(children ...Node) Node {
	return CreateElement("canvas", children...)
}

func Caption(children ...Node) Node {
	return CreateElement("caption", children...)
}

func Cite(children ...Node) Node {
	return CreateElement("cite", children...)
}

func Code(children ...Node) Node {
	return CreateElement("code", children...)
}

func Col(children ...Node) Node {
	return CreateElementVoid("col", children...)
}

func Colgroup(children ...Node) Node {
	return CreateElement("colgroup", children...)
}

func DataEl(children ...Node) Node {
	return CreateElement("data", children...)
}

func Datalist(children ...Node) Node {
	return CreateElement("datalist", children...)
}

func Dd(children ...Node) Node {
	return CreateElement("dd", children...)
}

func Del(children ...Node) Node {
	return CreateElement("del", children...)
}

func Details(children ...Node) Node {
	return CreateElement("details", children...)
}

func Dfn(children ...Node) Node {
	return CreateElement("dfn", children...)
}

func Dialog(children ...Node) Node {
	return CreateElement("dialog", children...)
}

func Div(children ...Node) Node {
	return CreateElement("div", children...)
}

func Dl(children ...Node) Node {
	return CreateElement("dl", children...)
}

func Dt(children ...Node) Node {
	return CreateElement("dt", children...)
}

func Em(children ...Node) Node {
	return CreateElement("em", children...)
}

func Embed(children ...Node) Node {
	return CreateElementVoid("embed", children...)
}

func Fieldset(children ...Node) Node {
	return CreateElement("fieldset", children...)
}

func Figcaption(children ...Node) Node {
	return CreateElement("figcaption", children...)
}

func Figure(children ...Node) Node {
	return CreateElement("figure", children...)
}

func Footer(children ...Node) Node {
	return CreateElement("footer", children...)
}

func Form(children ...Node) Node {
	return CreateElement("form", children...)
}

func H1(children ...Node) Node {
	return CreateElement("h1", children...)
}

func H2(children ...Node) Node {
	return CreateElement("h2", children...)
}

func H3(children ...Node) Node {
	return CreateElement("h3", children...)
}

func H4(children ...Node) Node {
	return CreateElement("h4", children...)
}

func H5(children ...Node) Node {
	return CreateElement("h5", children...)
}

func H6(children ...Node) Node {
	return CreateElement("h6", children...)
}

func Head(children ...Node) Node {
	return CreateElement("head", children...)
}

func Header(children ...Node) Node {
	return CreateElement("header", children...)
}

func Hr(children ...Node) Node {
	return CreateElementVoid("hr", children...)
}

func Html(children ...Node) Node {
	return CreateElement("html", children...)
}

func I(children ...Node) Node {
	return CreateElement("i", children...)
}

func Iframe(children ...Node) Node {
	return CreateElementVoid("iframe", children...)
}

func Img(children ...Node) Node {
	return CreateElementVoid("img", children...)
}

func Input(children ...Node) Node {
	return CreateElementVoid("input", children...)
}

func Ins(children ...Node) Node {
	return CreateElement("ins", children...)
}

func Kbd(children ...Node) Node {
	return CreateElement("kbd", children...)
}

func Label(children ...Node) Node {
	return CreateElement("label", children...)
}

func Legend(children ...Node) Node {
	return CreateElement("legend", children...)
}

func Li(children ...Node) Node {
	return CreateElement("li", children...)
}

func Link(children ...Node) Node {
	return CreateElementVoid("link", children...)
}

func Main(children ...Node) Node {
	return CreateElement("main", children...)
}

func MapEl(children ...Node) Node {
	return CreateElement("map", children...)
}

func Mark(children ...Node) Node {
	return CreateElement("mark", children...)
}

func Menu(children ...Node) Node {
	return CreateElement("menu", children...)
}

func Meta(children ...Node) Node {
	return CreateElementVoid("meta", children...)
}

func Meter(children ...Node) Node {
	return CreateElement("meter", children...)
}

func Nav(children ...Node) Node {
	return CreateElement("nav", children...)
}

func Noscript(children ...Node) Node {
	return CreateElement("noscript", children...)
}

func Object(children ...Node) Node {
	return CreateElement("object", children...)
}

func Ol(children ...Node) Node {
	return CreateElement("ol", children...)
}

func Optgroup(children ...Node) Node {
	return CreateElement("optgroup", children...)
}

func Option(children ...Node) Node {
	return CreateElement("option", children...)
}

func Output(children ...Node) Node {
	return CreateElement("output", children...)
}

func P(children ...Node) Node {
	return CreateElement("p", children...)
}

func Param(children ...Node) Node {
	return CreateElementVoid("param", children...)
}

func Picture(children ...Node) Node {
	return CreateElement("picture", children...)
}

func Pre(children ...Node) Node {
	return CreateElement("pre", children...)
}

func Progress(children ...Node) Node {
	return CreateElement("progress", children...)
}

func Q(children ...Node) Node {
	return CreateElement("q", children...)
}

func S(children ...Node) Node {
	return CreateElement("s", children...)
}

func Samp(children ...Node) Node {
	return CreateElement("samp", children...)
}

func Script(children ...Node) Node {
	return CreateElement("script", children...)
}

func Section(children ...Node) Node {
	return CreateElement("section", children...)
}

func Select(children ...Node) Node {
	return CreateElement("select", children...)
}

func SlotEl(children ...Node) Node {
	return CreateElement("slot", children...)
}

func Small(children ...Node) Node {
	return CreateElement("small", children...)
}

func Source(children ...Node) Node {
	return CreateElementVoid("source", children...)
}

func Span(children ...Node) Node {
	return CreateElement("span", children...)
}

func Strong(children ...Node) Node {
	return CreateElement("strong", children...)
}

func StyleEl(children ...Node) Node {
	return CreateElement("style", children...)
}

func Sub(children ...Node) Node {
	return CreateElement("sub", children...)
}

func Summary(children ...Node) Node {
	return CreateElement("summary", children...)
}

func Sup(children ...Node) Node {
	return CreateElement("sup", children...)
}

func Table(children ...Node) Node {
	return CreateElement("table", children...)
}

func Tbody(children ...Node) Node {
	return CreateElement("tbody", children...)
}

func Td(children ...Node) Node {
	return CreateElement("td", children...)
}

func Template(children ...Node) Node {
	return CreateElement("template", children...)
}

func Textarea(children ...Node) Node {
	return CreateElement("textarea", children...)
}

func Tfoot(children ...Node) Node {
	return CreateElement("tfoot", children...)
}

func Th(children ...Node) Node {
	return CreateElement("th", children...)
}

func Thead(children ...Node) Node {
	return CreateElement("thead", children...)
}

func Time(children ...Node) Node {
	return CreateElement("time", children...)
}

func TitleEl(children ...Node) Node {
	return CreateElement("title", children...)
}

func Tr(children ...Node) Node {
	return CreateElement("tr", children...)
}

func Track(children ...Node) Node {
	return CreateElementVoid("track", children...)
}

func U(children ...Node) Node {
	return CreateElement("u", children...)
}

func Ul(children ...Node) Node {
	return CreateElement("ul", children...)
}

func Var(children ...Node) Node {
	return CreateElement("var", children...)
}

func Video(children ...Node) Node {
	return CreateElement("video", children...)
}

func Wbr(children ...Node) Node {
	return CreateElement("wbr", children...)
}

// ------------------------------------------------------------------
//
// Dom Attributes
//
// ------------------------------------------------------------------

func Accesskey(v string) Node {
	return CreateAttr("accesskey", v)
}

func Accept(v string) Node {
	return CreateAttr("accept", v)
}

func Action(v string) Node {
	return CreateAttr("action", v)
}

func Allowfullscreen() Node {
	return CreateAttrBoolean("allowfullscreen")
}

func Alt(v string) Node {
	return CreateAttr("alt", v)
}

func Aria(name, v string) Node {
	return CreateAttr("aria-"+name, v)
}

func As(v string) Node {
	return CreateAttr("as", v)
}

func Async() Node {
	return CreateAttrBoolean("async")
}

func Autocapitalize(v string) Node {
	return CreateAttr("autocapitalize", v)
}

func Autocomplete(v string) Node {
	return CreateAttr("autocomplete", v)
}

func Autofocus() Node {
	return CreateAttrBoolean("autofocus")
}

func Autoplay() Node {
	return CreateAttrBoolean("autoplay")
}

func Charset(v string) Node {
	return CreateAttr("charset", v)
}

func Checked() Node {
	return CreateAttrBoolean("checked")
}

func CiteAttr(v string) Node {
	return CreateAttr("cite", v)
}

func Class(v string) Node {
	return CreateAttr("class", v)
}

func Cols(v string) Node {
	return CreateAttr("cols", v)
}

func Colspan(v string) Node {
	return CreateAttr("colspan", v)
}

func Content(v string) Node {
	return CreateAttr("content", v)
}

func Contenteditable(v string) Node {
	return CreateAttr("contenteditable", v)
}

func Controls() Node {
	return CreateAttrBoolean("controls")
}

func Crossorigin(v string) Node {
	return CreateAttr("crossorigin", v)
}

func Data(name, v string) Node {
	return CreateAttr("data-"+name, v)
}

func Datetime(v string) Node {
	return CreateAttr("datetime", v)
}

func Defer() Node {
	return CreateAttrBoolean("defer")
}

func Dir(v string) Node {
	return CreateAttr("dir", v)
}

func Disabled() Node {
	return CreateAttrBoolean("disabled")
}

func Download(v string) Node {
	return CreateAttr("download", v)
}

func Draggable(v string) Node {
	return CreateAttr("draggable", v)
}

func Enctype(v string) Node {
	return CreateAttr("enctype", v)
}

func Enterkeyhint(v string) Node {
	return CreateAttr("enterkeyhint", v)
}

func For(v string) Node {
	return CreateAttr("for", v)
}

func FormAttr(v string) Node {
	return CreateAttr("form", v)
}

func Formnovalidate() Node {
	return CreateAttrBoolean("formnovalidate")
}

func Headers(v string) Node {
	return CreateAttr("headers", v)
}

func Height(v string) Node {
	return CreateAttr("height", v)
}

func Hidden() Node {
	return CreateAttrBoolean("hidden")
}

func Href(v string) Node {
	return CreateAttr("href", v)
}

func Hreflang(v string) Node {
	return CreateAttr("hreflang", v)
}

func Httpequiv(v string) Node {
	return CreateAttr("http-equiv", v)
}

func Id(v string) Node {
	return CreateAttr("id", v)
}

func Inert(v string) Node {
	return CreateAttr("inert", v)
}

func Inputmode(v string) Node {
	return CreateAttr("inputmode", v)
}

func Integrity(v string) Node {
	return CreateAttr("integrity", v)
}

func Itemid(v string) Node {
	return CreateAttr("itemid", v)
}

func Itemprop(v string) Node {
	return CreateAttr("itemprop", v)
}

func Itemref(v string) Node {
	return CreateAttr("itemref", v)
}

func Itemscope() Node {
	return CreateAttrBoolean("itemscope")
}

func Itemtype(v string) Node {
	return CreateAttr("itemtype", v)
}

func Lang(v string) Node {
	return CreateAttr("lang", v)
}

func List(v string) Node {
	return CreateAttr("list", v)
}

func Loading(v string) Node {
	return CreateAttr("loading", v)
}

func Loop() Node {
	return CreateAttrBoolean("loop")
}

func Max(v string) Node {
	return CreateAttr("max", v)
}

func Maxlength(v string) Node {
	return CreateAttr("maxlength", v)
}

func Media(v string) Node {
	return CreateAttr("media", v)
}

func Method(v string) Node {
	return CreateAttr("method", v)
}

func Min(v string) Node {
	return CreateAttr("min", v)
}

func Minlength(v string) Node {
	return CreateAttr("minlength", v)
}

func Multiple() Node {
	return CreateAttrBoolean("multiple")
}

func Muted() Node {
	return CreateAttrBoolean("muted")
}

func Name(v string) Node {
	return CreateAttr("name", v)
}

func Nomodule() Node {
	return CreateAttrBoolean("nomodule")
}

func Nonce(v string) Node {
	return CreateAttr("nonce", v)
}

func Novalidate() Node {
	return CreateAttrBoolean("novalidate")
}

func Onload(v string) Node {
	return CreateAttr("onload", v)
}

func Open() Node {
	return CreateAttrBoolean("open")
}

func Part(v string) Node {
	return CreateAttr("part", v)
}

func Pattern(v string) Node {
	return CreateAttr("pattern", v)
}

func Placeholder(v string) Node {
	return CreateAttr("placeholder", v)
}

func Popover(v string) Node {
	return CreateAttr("popover", v)
}

func Poster(v string) Node {
	return CreateAttr("poster", v)
}

func Preload(v string) Node {
	return CreateAttr("preload", v)
}

func Property(v string) Node {
	return CreateAttr("property", v)
}

func Readonly() Node {
	return CreateAttrBoolean("readonly")
}

func Referrerpolicy(v string) Node {
	return CreateAttr("referrerpolicy", v)
}

func Rel(v string) Node {
	return CreateAttr("rel", v)
}

func Required() Node {
	return CreateAttrBoolean("required")
}

func Role(v string) Node {
	return CreateAttr("role", v)
}

func Rows(v string) Node {
	return CreateAttr("rows", v)
}

func Rowspan(v string) Node {
	return CreateAttr("rowspan", v)
}

func Sandbox(v string) Node {
	return CreateAttr("sandbox", v)
}

func Scope(v string) Node {
	return CreateAttr("scope", v)
}

func Selected() Node {
	return CreateAttrBoolean("selected")
}

func Size(v string) Node {
	return CreateAttr("size", v)
}

func Sizes(v string) Node {
	return CreateAttr("sizes", v)
}

func Slot(v string) Node {
	return CreateAttr("slot", v)
}

func Spellcheck(v string) Node {
	return CreateAttr("spellcheck", v)
}

func Src(v string) Node {
	return CreateAttr("src", v)
}

func Srcset(v string) Node {
	return CreateAttr("srcset", v)
}

func Step(v string) Node {
	return CreateAttr("step", v)
}

func Style(v string) Node {
	return CreateAttr("style", v)
}

func Tabindex(v string) Node {
	return CreateAttr("tabindex", v)
}

func Target(v string) Node {
	return CreateAttr("target", v)
}

func Title(v string) Node {
	return CreateAttr("title", v)
}

func Translate(v string) Node {
	return CreateAttr("translate", v)
}

func Type(v string) Node {
	return CreateAttr("type", v)
}

func Value(v string) Node {
	return CreateAttr("value", v)
}

func Width(v string) Node {
	return CreateAttr("width", v)
}
