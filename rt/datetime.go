package rt

import (
	"strings"
	"time"
)

// Date formats a time.Time to a string like "2006-01-02".
func Date(d time.Time) string {
	return d.Format("2006-01-02")
}

// Time formats a time.Time to a string like "3:04 PM".
func Time(d time.Time) string {
	return d.Format("3:04 PM")
}

// DateTime formats a time.Time to a string like "January 02, 2006, 3:04 PM".
func DateTime(d time.Time) string {
	return d.Format("January 02, 2006, 3:04 PM")
}

// DateNatural formats a time.Time to a string like "January 2, 2006".
func DateNatural(d time.Time) string {
	return d.Format("January 2, 2006")
}

// DateSlug formats a time.Time to a string like "January-2-2006".
func DateSlug(d time.Time) string {
	return strings.ToLower(d.Format("January-2-2006"))
}

// DateSlugIntl formats a time.Time to a string like "2-January-2006".
func DateSlugIntl(d time.Time) string {
	return strings.ToLower(d.Format("2-January-2006"))
}

var dateFormats = []string{
	"2006-01-02",      // Short
	"January-2-2006",  // Slug
	"2-January-2006",  // Slug (international)
	"January 2, 2006", // Display
}

// ParseDate parses a date string from multiple layouts.
//
// The following layouts are tried:
//   - "2006-01-02"
//   - "January-2-2006"
//   - "2-January-2006"
//   - "January 2, 2006"
func ParseDate(val string) (time.Time, error) {
	var date time.Time
	var err error
	for _, format := range dateFormats {
		date, err = time.Parse(format, val)
		// If the time parsing fails,
		// then try the next format.
		if err != nil {
			continue
		}
		// If there is no error, then
		// return the parsed time.Time.
		return date, nil
	}
	return time.Time{}, err
}
