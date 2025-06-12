package dom

func Text(value string) Node {
	return CreateString(value)
}

func Raw(value string) Node {
	return CreateStringRaw(value)
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
