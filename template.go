package rio

import (
	"html/template"
	"time"
)

// ------------------------------------------------------------------
//
//
// Functions for template.FuncMap
//
//
// ------------------------------------------------------------------

func DisplaySafeHTML(content string) template.HTML {
	return template.HTML(content)
}

func DisplayTime(d time.Time) string {
	return d.Format("3:04 PM")
}

func DisplayDate(d time.Time) string {
	return d.Format("January 02, 2006")
}

func DisplayDateTime(d time.Time) string {
	return d.Format("January 02, 2006, 3:04 PM")
}

// ------------------------------------------------------------------
//
//
// Helper Functions for template.FuncMap
//
//
// ------------------------------------------------------------------

func WrapString(val string) func() string {
	return func() string {
		return val
	}
}

func WrapBool(val bool) func() bool {
	return func() bool {
		return val
	}
}

func WrapInt(val int) func() int {
	return func() int {
		return val
	}
}

func WrapFloat(val float64) func() float64 {
	return func() float64 {
		return val
	}
}

func WrapTime(val time.Time) func() time.Time {
	return func() time.Time {
		return val
	}
}
