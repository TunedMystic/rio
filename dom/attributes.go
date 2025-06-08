package dom

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
