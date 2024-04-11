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

// DisplaySafeHTML converts a string into an HTML fragment, so that
// it can be rendered verbatim in the template.
func DisplaySafeHTML(content string) template.HTML {
	return template.HTML(content)
}

// DisplayTime formats a time.Time value like "3:04 PM".
func DisplayTime(d time.Time) string {
	return d.Format("3:04 PM")
}

// DisplayDate formats a time.Time value like "January 02, 2006".
func DisplayDate(d time.Time) string {
	return d.Format("January 02, 2006")
}

// DisplayDateTime formats a time.Time value like "January 02, 2006, 3:04 PM".
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

// WrapString wraps a string value in a func.
// This is used to inject values into the template.FuncMap.
func WrapString(val string) func() string {
	return func() string {
		return val
	}
}

// WrapBool wraps a boolean value in a func.
// This is used to inject values into the template.FuncMap.
func WrapBool(val bool) func() bool {
	return func() bool {
		return val
	}
}

// WrapInt wraps an integer value in a func.
// This is used to inject values into the template.FuncMap.
func WrapInt(val int) func() int {
	return func() int {
		return val
	}
}

// WrapFloat wraps a float value in a func.
// This is used to inject values into the template.FuncMap.
func WrapFloat(val float64) func() float64 {
	return func() float64 {
		return val
	}
}

// WrapTime wraps a time.Time value in a func.
// This is used to inject values into the template.FuncMap.
func WrapTime(val time.Time) func() time.Time {
	return func() time.Time {
		return val
	}
}

// WrapItems wraps a list of any type in a func.
// This is used to inject values into the template.FuncMap.
func WrapItems[T any](vals []T) func() []T {
	return func() []T {
		return vals
	}
}
