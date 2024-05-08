package rt

import (
	"strings"
	"time"
)

// Today returns a time.Time value which represents the current day.
func Today() time.Time {
	return TrimTime(time.Now().UTC())
}

// CurrentYear returns the current year as a int.
func CurrentYear() int {
	return time.Now().Year()
}

// TrimTime erases the hh:mm:ss:ns of the given time.Time.
func TrimTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// TrimZero removes the trailing ".00" from a string.
func TrimZero(s string) string {
	return strings.Replace(s, ".00", "", -1)
}

// Ordered represents integers and floating-point values.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// Plural returns the singular term if n == 1, and returns
// the plural value if n != 1.
func Plural[T Ordered](n T, singular, plural string) string {
	if n == T(1) {
		return singular
	}
	return plural
}

// WrapItem wraps a value of any type in a func.
// This is used to inject values into the template.FuncMap.
func WrapItem[T any](val T) func() T {
	return func() T {
		return val
	}
}

// WrapSlice wraps a slice of any type in a func.
// This is used to inject values into the template.FuncMap.
func WrapSlice[T any](val []T) func() []T {
	return func() []T {
		return val
	}
}

// WrapMap wraps a map of type [string]any in a func.
// This is used to inject values into the template.FuncMap.
func WrapMap[T any](val map[string]T) func() map[string]T {
	return func() map[string]T {
		return val
	}
}
