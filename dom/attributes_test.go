package dom

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tunedmystic/rio/internal/assert"
)

func Test_CreateAttr(t *testing.T) {
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

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Test_Attributes(t *testing.T) {
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

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Test_Attr_Render_Errors(t *testing.T) {
	errWriterSentinel := errors.New("writer error from htmlAttr test")

	tests := []struct {
		name           string
		attr           *htmlAttr
		failOnNthWrite int
		expectedErr    error
		expectNilError bool
	}{
		{
			name:           "non-escaping: fail write space",
			attr:           &htmlAttr{Name: "class", Value: "test"},
			failOnNthWrite: 1,
			expectedErr:    errWriterSentinel,
		},
		{
			name:           "non-escaping: fail write name",
			attr:           &htmlAttr{Name: "class", Value: "test"},
			failOnNthWrite: 2,
			expectedErr:    errWriterSentinel,
		},
		{
			name:           "non-escaping: fail write equals-quote",
			attr:           &htmlAttr{Name: "class", Value: "test"},
			failOnNthWrite: 3,
			expectedErr:    errWriterSentinel,
		},
		{
			name:           "non-escaping: fail write value",
			attr:           &htmlAttr{Name: "class", Value: "test"},
			failOnNthWrite: 4,
			expectedErr:    errWriterSentinel,
		},
		{
			name:           "non-escaping: fail write closing quote",
			attr:           &htmlAttr{Name: "class", Value: "test"},
			failOnNthWrite: 5,
			expectedErr:    errWriterSentinel,
		},
		{
			name:           "boolean: fail write space",
			attr:           &htmlAttr{Name: "hidden", Value: ""},
			failOnNthWrite: 1,
			expectedErr:    errWriterSentinel,
		},
		{
			name:           "boolean: fail write name",
			attr:           &htmlAttr{Name: "hidden", Value: ""},
			failOnNthWrite: 2,
			expectedErr:    errWriterSentinel,
		},
		{
			name:           "escaping: fail write space",
			attr:           &htmlAttr{Name: "id", Value: "<"},
			failOnNthWrite: 1,
			expectedErr:    errWriterSentinel,
		},
		{
			name:           "escaping: fail write name",
			attr:           &htmlAttr{Name: "id", Value: "<"},
			failOnNthWrite: 2,
			expectedErr:    errWriterSentinel,
		},
		{
			name:           "escaping: fail write equals-quote",
			attr:           &htmlAttr{Name: "id", Value: "<"},
			failOnNthWrite: 3,
			expectedErr:    errWriterSentinel,
		},
		{
			name:           "escaping: fail write value (during HTMLEscape - demonstrates SUT bug)",
			attr:           &htmlAttr{Name: "id", Value: "<"},
			failOnNthWrite: 4,
			expectedErr:    nil,
			expectNilError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &errorWriter{targetErr: tt.expectedErr, failOnNthWrite: tt.failOnNthWrite}
			err := tt.attr.Render(writer)
			assert.Error(t, err, tt.expectedErr)
		})
	}
}

// ------------------------------------------------------------------
//
//
//
// ------------------------------------------------------------------

func Benchmark_Attributes(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = Class("my-class")
		_ = Id("my-id")
		_ = Href("/some/url")
		_ = Src("image.png")
	}
}
