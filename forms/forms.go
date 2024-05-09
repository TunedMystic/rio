// Package forms implements utilities to parse and validate data.
package forms

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/tunedmystic/rio/format"
)

// ------------------------------------------------------------------
//
//
// Type: Form
//
//
// ------------------------------------------------------------------

// Form is a type which parses and validates form data.
//
// Values can be parsed into a desired type, like an int, date, bool, etc.
// Validation functions can also be provided, which ensures that
// the parsed type is properly vetted before being retrieved.
type Form struct {
	values  map[string]string
	cleaned map[string]Field
	errors  map[string]error

	extraerrors []error
}

// New constructs and returns a Form.
func New() *Form {
	return &Form{}
}

// ------------------------------------------------------------------
//
//
// Form Cleaners
//
//
// ------------------------------------------------------------------

// CleanString cleans the given value as a string.
func (f *Form) CleanString(name, value string, funcs ...CheckFunc) {
	f.parseAndClean(name, value, parseString, funcs...)
}

// CleanInteger cleans the given value as a integer.
func (f *Form) CleanInteger(name, value string, funcs ...CheckFunc) {
	f.parseAndClean(name, value, parseInteger, funcs...)
}

// CleanFloat cleans the given value as a float.
func (f *Form) CleanFloat(name, value string, funcs ...CheckFunc) {
	f.parseAndClean(name, value, parseFloat, funcs...)
}

// CleanBool cleans the given value as a bool.
func (f *Form) CleanBool(name, value string, funcs ...CheckFunc) {
	f.parseAndClean(name, value, parseBool, funcs...)
}

// CleanDate cleans the given value as a date.
func (f *Form) CleanDate(name, value string, funcs ...CheckFunc) {
	f.parseAndClean(name, value, parseDate, funcs...)
}

// CleanDecimal cleans the given value as a decimal.
func (f *Form) CleanDecimal(name, value string, funcs ...CheckFunc) {
	f.parseAndClean(name, value, parseDecimal, funcs...)
}

// CleanExtra adds the error to the errors list if the condition is true.
func (f *Form) CleanExtra(cond bool, err error) {
	if cond {
		f.addExtraError(err)
	}
}

// ------------------------------------------------------------------
//
//
// Form Getters for Cleaned Fields
//
//
// ------------------------------------------------------------------

// CleanedString retrieves the named field as a string.
func (f *Form) CleanedString(name string) string {
	return f.MustField(name).String
}

// CleanedInteger retrieves the named field as an integer.
func (f *Form) CleanedInteger(name string) int {
	return f.MustField(name).Integer
}

// CleanedFloat retrieves the named field as a float.
func (f *Form) CleanedFloat(name string) float64 {
	return f.MustField(name).Float
}

// CleanedBool retrieves the named field as a bool.
func (f *Form) CleanedBool(name string) bool {
	return f.MustField(name).Bool
}

// CleanedDate retrieves the named field as a date.
func (f *Form) CleanedDate(name string) time.Time {
	return f.MustField(name).Date
}

// CleanedDecimal retrieves the named field as a decimal.
func (f *Form) CleanedDecimal(name string) big.Rat {
	return f.MustField(name).Decimal
}

// ------------------------------------------------------------------
//
//
// Form Introspection
//
//
// ------------------------------------------------------------------

// IsValid returns true if the errors map contains no errors.
func (f *Form) IsValid() bool {
	return len(f.errors) == 0 && len(f.extraerrors) == 0
}

// HasError returns true if the errors map contains the target error.
func (f *Form) HasError(target any) bool {
	for _, err := range f.errors {
		if errors.As(err, target) {
			return true
		}
	}
	for _, err := range f.extraerrors {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

// ------------------------------------------------------------------
//
//
// Form Getters for Map Items
//
//
// ------------------------------------------------------------------

// Value returns the original, uncleaned value of a field.
func (f *Form) Value(name string) string {
	if f.values == nil {
		return ""
	}
	return f.values[name]
}

// Field returns the cleaned Field.
//
// When a clean function is successful, the cleaned map will be
// populated with the parsed value, as a Field.
func (f *Form) Field(name string) (Field, bool) {
	if f.cleaned == nil {
		return Field{}, false
	}
	field, ok := f.cleaned[name]
	return field, ok
}

// MustField returns the desired Field and panics if it does not exist.
func (f *Form) MustField(name string) Field {
	field, ok := f.Field(name)
	if !ok {
		panic(fmt.Sprintf("failed to get field %s", name))
	}
	return field
}

// ------------------------------------------------------------------
//
//
// Form Getters for Errors
//
//
// ------------------------------------------------------------------

// Error returns the error for a field.
//
// When a clean function fails, the errors map will be
// populated with the error.
func (f *Form) Error(name string) error {
	if f.errors == nil {
		return nil
	}
	return f.errors[name]
}

// ErrorNames returns the names of the fields with errors.
func (f *Form) ErrorNames() []string {
	var names []string
	for name := range f.errors {
		names = append(names, name)
	}
	return names
}

// ExtraErrors returns the errors slice.
//
// When a custom check fails, the errors slice will be
// populated with the error.
func (f *Form) ExtraErrors() []error {
	return f.extraerrors
}

// ------------------------------------------------------------------
//
//
// Form Internals
//
//
// ------------------------------------------------------------------

// parseAndClean is the internal function for the cleaning workflow.
//
// First, the value is parsed into the desired type, as a Field.
// Second, the Field is validated against the provided check functions.
//
// If the First and Second steps are succesful, the Field is added to the cleaned map.
// If the First or Second steps are not successful, the error is added to the errors map.
func (f *Form) parseAndClean(name, value string, parse ParseFunc, checks ...CheckFunc) {
	// If the value is already processed, then skip.
	if f.isProcessed(name) {
		return
	}

	// Keep a copy of the value.
	f.addValue(name, value)

	// Parse the value.
	field, err := parse(value)
	if err != nil {
		f.addError(name, err)
		return
	}

	// Validate the field.
	for i := range checks {
		if err := checks[i](field); err != nil {
			f.addError(name, err)
			return
		}
	}

	// At this point, the field has passed all check functions.
	// In this case, add the field to the cleaned map.
	f.addField(name, field)
}

// addValue adds the value into the values map.
func (f *Form) addValue(name, val string) {
	if f.values == nil {
		f.values = make(map[string]string)
	}

	if _, exists := f.values[name]; !exists {
		f.values[name] = val
	}
}

// addField adds the field into the cleaned map.
func (f *Form) addField(name string, val Field) {
	if f.cleaned == nil {
		f.cleaned = make(map[string]Field)
	}

	if _, exists := f.cleaned[name]; !exists {
		f.cleaned[name] = val
	}
}

// addError adds the error into the errors map.
func (f *Form) addError(name string, err error) {
	if f.errors == nil {
		f.errors = make(map[string]error)
	}

	if _, exists := f.errors[name]; !exists {
		f.errors[name] = err
	}
}

// addExtraError adds the error to the extra errors list.
func (f *Form) addExtraError(err error) {
	f.extraerrors = append(f.extraerrors, err)
}

// isProcessed returns true if the value exists in the cleaned or errors map.
func (f *Form) isProcessed(name string) bool {
	_, inCleaned := f.cleaned[name]
	_, inErrors := f.errors[name]
	return inCleaned || inErrors
}

// ------------------------------------------------------------------
//
//
// Type: Field
//
//
// ------------------------------------------------------------------

// Field represents a parsed value.
type Field struct {
	String  string
	Integer int
	Float   float64
	Bool    bool
	Decimal big.Rat
	Date    time.Time

	isBlank bool
}

// ParseFunc is a function which parses a value into the desired type, as a Field.
type ParseFunc func(string) (Field, error)

// CheckFunc is a function which validates a Field.
type CheckFunc func(Field) error

// ------------------------------------------------------------------
//
//
// Type: ParseError
//
//
// ------------------------------------------------------------------

// ParseError represents when a value fails to parse into a specific type.
type ParseError struct {
	Msg string
}

func (p ParseError) Error() string {
	return p.String()
}

func (p ParseError) String() string {
	return p.Msg
}

// ------------------------------------------------------------------
//
//
// Errors and Constants
//
//
// ------------------------------------------------------------------

var (
	errParseInt    = ParseError{"must be a valid integer"}
	errParseDate   = ParseError{"must be a valid date"}
	errParseBool   = ParseError{"must be a valid boolean"}
	errParseBigRat = ParseError{"must be a valid decimal"}
	errParseFloat  = ParseError{"must be a valid float"}

	errInvalidChoice = errors.New("must be a valid choice")
	errInvalidConfig = errors.New("invalid validation config")
	errBlankValue    = errors.New("cannot be blank")

	errLessThan           = "must be less than %v"
	errLessThanOrEqual    = "must be less than or equal to %v"
	errGreaterThan        = "must be more than %v"
	errGreaterThanOrEqual = "must be more than or equal to %v"
	errBetween            = "must be between %v and %v"
	errBefore             = "must be before %v"
	errAfter              = "must be after %v"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	urlRegex   = regexp.MustCompile(`^(http(s)?://)?([\da-z\.-]+)\.([a-z\.]{2,6})([/\w \.-]*)*/?$`)
)

// ------------------------------------------------------------------
//
//
// Parse Functions
//
//
// ------------------------------------------------------------------

// parseString parses the value into a string Field.
func parseString(val string) (Field, error) {
	if val == "" {
		return Field{isBlank: true}, nil
	}
	return Field{String: strings.TrimSpace(val)}, nil
}

// parseInteger parses the value into an integer Field.
func parseInteger(val string) (Field, error) {
	if val == "" {
		return Field{isBlank: true}, nil
	}

	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return Field{}, errParseInt
	}

	return Field{Integer: int(num)}, nil
}

// parseFloat parses the value into a float64 Field.
func parseFloat(val string) (Field, error) {
	if val == "" {
		return Field{isBlank: true}, nil
	}

	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return Field{}, errParseFloat
	}

	return Field{Float: num}, nil
}

// parseBool parses the value into a bool Field.
func parseBool(val string) (Field, error) {
	if val == "" {
		return Field{isBlank: true}, nil
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		return Field{}, errParseBool
	}

	return Field{Bool: b}, nil
}

// parseDecimal parses the value into a decimal Field.
func parseDecimal(val string) (Field, error) {
	if val == "" {
		return Field{isBlank: true}, nil
	}

	br, ok := new(big.Rat).SetString(val)
	if !ok {
		return Field{}, errParseBigRat
	}

	return Field{Decimal: *br}, nil
}

// parseDate parses the value into a date Field.
func parseDate(val string) (Field, error) {
	if val == "" {
		return Field{isBlank: true}, nil
	}

	date, err := format.ParseDate(val)
	if err != nil {
		return Field{}, errParseDate
	}

	return Field{Date: date}, nil
}

// ------------------------------------------------------------------
//
//
// String Check Functions
//
//
// ------------------------------------------------------------------

// Checks that a string is not blank.
func StrRequired() CheckFunc {
	return func(v Field) error {
		if v.isBlank {
			return errBlankValue
		}
		return nil
	}
}

// Checks that a string's length is less than n.
func StrLt(n int) CheckFunc {
	err := fmt.Errorf("must be less than %d characters", n)

	return func(v Field) error {
		if utf8.RuneCountInString(v.String) >= n {
			return err
		}
		return nil
	}
}

// Checks that a string's length is less than or equal to n.
func StrLte(n int) CheckFunc {
	err := fmt.Errorf("must be less than or equal to %d characters", n)

	return func(v Field) error {
		if utf8.RuneCountInString(v.String) > n {
			return err
		}
		return nil
	}
}

// Checks that a string's length is greater than n.
func StrGt(n int) CheckFunc {
	err := fmt.Errorf("must be more than %d characters", n)

	return func(v Field) error {
		if utf8.RuneCountInString(v.String) <= n {
			return err
		}
		return nil
	}
}

// Checks that a string's length is greater than or equal to n.
func StrGte(n int) CheckFunc {
	err := fmt.Errorf("must be more than or equal to %d characters", n)

	return func(v Field) error {
		if utf8.RuneCountInString(v.String) <= n {
			return err
		}
		return nil
	}
}

// Checks that a string's length is between n and m.
func StrBtw(n, m int) CheckFunc {
	err := fmt.Errorf("must be between %d and %d characters", n, m)

	return func(v Field) error {
		if (n == m) || (n > m) {
			return errInvalidConfig
		}
		if strLen := utf8.RuneCountInString(v.String); strLen < n || strLen > m {
			return err
		}
		return nil
	}
}

// Checks that a string is a member of the given choices.
func StrIn(choices []string) CheckFunc {
	return func(v Field) error {
		if !slices.Contains(choices, v.String) {
			return errInvalidChoice
		}
		return nil
	}
}

// Checks that a string matches the given regex.
func StrMatches(rx *regexp.Regexp, errMsg string) CheckFunc {
	err := fmt.Errorf(errMsg)

	return func(v Field) error {
		if !rx.MatchString(v.String) {
			return err
		}
		return nil
	}
}

// Checks that a string is an email.
func StrEmail() CheckFunc {
	return StrMatches(emailRegex, "must be a valid email")
}

// Checks that a string is a URL.
func StrUrl() CheckFunc {
	return StrMatches(urlRegex, "must be a valid url")
}

// ------------------------------------------------------------------
//
//
// Integer Check Functions
//
//
// ------------------------------------------------------------------

// Checks that an int is not blank.
func IntRequired() CheckFunc {
	return func(v Field) error {
		if v.isBlank {
			return errBlankValue
		}
		return nil
	}
}

// Checks that an int is less than n.
func IntLt(n int) CheckFunc {
	return func(v Field) error {
		if v.Integer >= n {
			return fmt.Errorf(errLessThan, n)
		}
		return nil
	}
}

// Checks that an int is less than or equal to n.
func IntLte(n int) CheckFunc {
	return func(v Field) error {
		if v.Integer > n {
			return fmt.Errorf(errLessThanOrEqual, n)
		}
		return nil
	}
}

// Checks that an int is more than n.
func IntGt(n int) CheckFunc {
	return func(v Field) error {
		if v.Integer <= n {
			return fmt.Errorf(errGreaterThan, n)
		}
		return nil
	}
}

// Checks that an int is more than or equal to n.
func IntGte(n int) CheckFunc {
	return func(v Field) error {
		if v.Integer < n {
			return fmt.Errorf(errGreaterThanOrEqual, n)
		}
		return nil
	}
}

// Checks that an int is between n and m.
func IntBtw(n, m int) CheckFunc {
	return func(v Field) error {
		if (n == m) || (n > m) {
			return errInvalidConfig
		}
		if (v.Integer < n) || (v.Integer > m) {
			return fmt.Errorf(errBetween, n, m)
		}
		return nil
	}
}

// Checks that an int is a member of the given choices.
func IntIn(choices []int) CheckFunc {
	return func(v Field) error {
		if !slices.Contains(choices, v.Integer) {
			return errInvalidChoice
		}
		return nil
	}
}

// ------------------------------------------------------------------
//
//
// Float Check Functions
//
//
// ------------------------------------------------------------------

// Checks that a float is not blank.
func FltRequired() CheckFunc {
	return func(v Field) error {
		if v.isBlank {
			return errBlankValue
		}
		return nil
	}
}

// Checks that a float is less than n.
func FltLt(n float64) CheckFunc {
	return func(v Field) error {
		if v.Float >= n {
			return fmt.Errorf(errLessThan, n)
		}
		return nil
	}
}

// Checks that a float is less than or equal to n.
func FltLte(n float64) CheckFunc {
	return func(v Field) error {
		if v.Float > n {
			return fmt.Errorf(errLessThanOrEqual, n)
		}
		return nil
	}
}

// Checks that a float is more than n.
func FltGt(n float64) CheckFunc {
	return func(v Field) error {
		if v.Float <= n {
			return fmt.Errorf(errGreaterThan, n)
		}
		return nil
	}
}

// Checks that a float is more than or equal to n.
func FltGte(n float64) CheckFunc {
	return func(v Field) error {
		if v.Float < n {
			return fmt.Errorf(errGreaterThanOrEqual, n)
		}
		return nil
	}
}

// Checks that a float is between n and m.
func FltBtw(n, m float64) CheckFunc {
	return func(v Field) error {
		if (n == m) || (n > m) {
			return errInvalidConfig
		}
		if (v.Float < n) || (v.Float > m) {
			return fmt.Errorf(errBetween, n, m)
		}
		return nil
	}
}

// Checks that a float is a member of the given choices.
func FltIn(choices []float64) CheckFunc {
	return func(v Field) error {
		if !slices.Contains(choices, v.Float) {
			return errInvalidChoice
		}
		return nil
	}
}

// ------------------------------------------------------------------
//
//
// Bool Check Functions
//
//
// ------------------------------------------------------------------

// Checks that a bool is not blank.
func BoolRequired() CheckFunc {
	return func(v Field) error {
		if v.Bool {
			return errBlankValue
		}
		return nil
	}
}

// ------------------------------------------------------------------
//
//
// Decimal Check Functions
//
//
// ------------------------------------------------------------------

// Checks that a decimal is not blank.
func DecRequired() CheckFunc {
	return func(v Field) error {
		if v.isBlank {
			return errBlankValue
		}
		return nil
	}
}

// Checks that a decimal is less than n.
func DecLt(n string) CheckFunc {
	nn, err := parseDecimal(n)

	return func(v Field) error {
		if err != nil {
			return errInvalidConfig
		}
		// Check if v >= nn.
		if r := v.Decimal.Cmp(&nn.Decimal); (r == 0) || (r == 1) {
			return fmt.Errorf(errLessThan, format.Decimal(nn.Decimal))
		}
		return nil
	}
}

// Checks that a decimal is less than or equal to n.
func DecLte(n string) CheckFunc {
	nn, err := parseDecimal(n)

	return func(v Field) error {
		if err != nil {
			return errInvalidConfig
		}
		// Check if v > nn.
		if r := v.Decimal.Cmp(&nn.Decimal); r == 1 {
			return fmt.Errorf(errLessThanOrEqual, format.Decimal(nn.Decimal))
		}
		return nil
	}
}

// Checks that a decimal is more than n.
func DecGt(n string) CheckFunc {
	nn, err := parseDecimal(n)

	return func(v Field) error {
		if err != nil {
			return errInvalidConfig
		}
		// Check if v <= nn.
		if r := v.Decimal.Cmp(&nn.Decimal); (r == -1) || (r == 0) {
			return fmt.Errorf(errGreaterThan, format.Decimal(nn.Decimal))
		}
		return nil
	}
}

// Checks that a decimal is greater than or equal to n.
func DecGte(n string) CheckFunc {
	nn, err := parseDecimal(n)

	return func(v Field) error {
		if err != nil {
			return errInvalidConfig
		}
		// Check if v < nn.
		if r := v.Decimal.Cmp(&nn.Decimal); r == -1 {
			return fmt.Errorf(errGreaterThanOrEqual, format.Decimal(nn.Decimal))
		}
		return nil
	}
}

// Checks that a decimal is between n and m.
func DecBtw(n, m string) CheckFunc {
	nn, nerr := parseDecimal(n)
	mm, merr := parseDecimal(m)

	return func(v Field) error {
		if (nerr != nil) || (merr != nil) {
			return errInvalidConfig
		}

		// Check if (n == m) || (n > m)
		r := nn.Decimal.Cmp(&mm.Decimal)
		if (r == 0) || (r == 1) {
			return errInvalidConfig
		}

		// Check if (v < n) || (v > m)
		nr := v.Decimal.Cmp(&nn.Decimal)
		mr := v.Decimal.Cmp(&mm.Decimal)
		if (nr == -1) || (mr == 1) {
			return fmt.Errorf(errBetween, format.Decimal(nn.Decimal), format.Decimal(mm.Decimal))
		}
		return nil
	}
}

// ------------------------------------------------------------------
//
//
// Date Check Functions
//
//
// ------------------------------------------------------------------

// Checks that a date is not blank.
func DtRequired() CheckFunc {
	return func(v Field) error {
		if v.isBlank {
			return errBlankValue
		}
		return nil
	}
}

// Checks that a date is before n (yyyy-mm-dd).
func DtBefore(n string) CheckFunc {
	nn, err := parseDate(n)

	return func(v Field) error {
		if err != nil {
			return errInvalidConfig
		}
		if !v.Date.Before(nn.Date) {
			return fmt.Errorf(errBefore, format.DateNatural(nn.Date))
		}
		return nil
	}
}

// Checks that a date is after n (yyyy-mm-dd).
func DtAfter(n string) CheckFunc {
	nn, err := parseDate(n)

	return func(v Field) error {
		if err != nil {
			return errInvalidConfig
		}
		if !v.Date.After(nn.Date) {
			return fmt.Errorf(errAfter, format.DateNatural(nn.Date))
		}
		return nil
	}
}

// Checks that a date is in the past.
func DtInPast() CheckFunc {
	return func(v Field) error {
		if !v.Date.Before(format.Today()) {
			return fmt.Errorf(errBefore, "the current date")
		}
		return nil
	}
}

// Checks that a date is in the future.
func DtInFuture() CheckFunc {
	return func(v Field) error {
		if !v.Date.After(format.Today()) {
			return fmt.Errorf(errAfter, "the current date")
		}
		return nil
	}
}
