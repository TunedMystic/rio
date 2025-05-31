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
	fields      map[string]Field
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
	f.cleanField(name, parseString(value), funcs...)
}

// CleanInteger cleans the given value as a integer.
func (f *Form) CleanInteger(name, value string, funcs ...CheckFunc) {
	f.cleanField(name, parseInteger(value), funcs...)
}

// CleanFloat cleans the given value as a float.
func (f *Form) CleanFloat(name, value string, funcs ...CheckFunc) {
	f.cleanField(name, parseFloat(value), funcs...)
}

// CleanBool cleans the given value as a bool.
func (f *Form) CleanBool(name, value string, funcs ...CheckFunc) {
	f.cleanField(name, parseBool(value), funcs...)
}

// CleanDate cleans the given value as a date.
func (f *Form) CleanDate(name, value string, funcs ...CheckFunc) {
	f.cleanField(name, parseDate(value), funcs...)
}

// CleanDecimal cleans the given value as a decimal.
func (f *Form) CleanDecimal(name, value string, funcs ...CheckFunc) {
	f.cleanField(name, parseDecimal(value), funcs...)
}

// CleanExtra adds the error to the extra errors list if the condition is true.
func (f *Form) CleanExtra(cond bool, err error) {
	if cond {
		f.extraerrors = append(f.extraerrors, err)
	}
}

// cleanField is the internal function for the cleaning workflow.
//
// First, the Field is checked for parse errors, and if they exist then
// then workflow is halted.
//
// Then, the Field is validated against the provided check functions, and
// any error encountered is added to the Field.
//
// If all validation funcs are successful with no errors, then
// the field is determined to be valid.
// .
func (f *Form) cleanField(name string, field Field, checks ...CheckFunc) {
	if field.Err() != nil {
		f.addField(name, field)
		return
	}

	// Validate the field.
	for i := range checks {
		if err := checks[i](field); err != nil {
			field.addError(err)
			break
		}
	}

	f.addField(name, field)
}

// addField adds the field into the fields map.
func (f *Form) addField(name string, val Field) {
	if f.fields == nil {
		f.fields = make(map[string]Field)
	}

	if _, exists := f.fields[name]; !exists {
		f.fields[name] = val
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
// Form Getters for Fields and Errors
//
//
// ------------------------------------------------------------------

// Field returns the Field mapped to the given name.
func (f *Form) Field(name string) (Field, bool) {
	if f.fields == nil {
		return Field{}, false
	}
	field, ok := f.fields[name]
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

// Names returns the field names.
func (f *Form) Names() []string {
	names := make([]string, 0, len(f.fields))
	for name := range f.fields {
		names = append(names, name)
	}
	return names
}

// ExtraErrors returns the extra errors slice.
func (f *Form) ExtraErrors() []error {
	return f.extraerrors
}

// Errors returns a slice of all field and non-field errors.
// Field error messages are prepared as "{field name} - {error message}".
// Non-Field errors messages are collected as is.
func (f *Form) Errors() []error {
	var errs []error

	// Collect field errors
	for _, name := range f.Names() {
		field := f.MustField(name)

		if err := field.Err(); err != nil {
			if errs == nil {
				errs = make([]error, 0, len(f.Names()))
			}
			errs = append(errs, fmt.Errorf("%s %w", name, err))
		}
	}

	// Collect non-field errors
	errs = append(errs, f.ExtraErrors()...)

	return errs
}

// ------------------------------------------------------------------
//
//
// Form Introspection
//
//
// ------------------------------------------------------------------

// IsValid returns true if there are no field errors and no extra errors.
func (f *Form) IsValid() bool {
	for _, field := range f.fields {
		if field.Err() != nil {
			return false
		}
	}
	return len(f.extraerrors) == 0
}

// HasError returns true if the errors map contains the target error.
func (f *Form) HasError(target any) bool {
	for _, field := range f.fields {
		if errors.As(field.Err(), target) {
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
	errInvalidEmail  = errors.New("must be a valid email")
	errInvalidUrl    = errors.New("must be a valid url")
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
// Type: Field
//
//
// ------------------------------------------------------------------

// Field represents a parsed value.
type Field struct {
	val     string
	err     error
	isBlank bool

	String  string
	Integer int
	Float   float64
	Bool    bool
	Decimal big.Rat
	Date    time.Time
}

// Value returns the field's original value.
func (f Field) Value() string {
	return f.val
}

// IsBlank returns true if the field's original value is empty.
// This is used to distinguish between the field's zero value and if the field is empty.
func (f Field) IsBlank() bool {
	return f.isBlank
}

// Err returns the error associated with the field's parsing or validation.
func (f Field) Err() error {
	return f.err
}

// addError adds the error to the field.
func (f *Field) addError(err error) {
	if f.err != nil {
		return
	}
	f.err = err
}

// ParseFunc is a function which parses a value into the desired type, as a Field.
type ParseFunc func(string) Field

// CheckFunc is a function which validates a Field.
type CheckFunc func(Field) error

// ------------------------------------------------------------------
//
//
// Field Parse Functions
//
//
// ------------------------------------------------------------------

// parseString parses the value into a string Field.
func parseString(val string) Field {
	if val == "" {
		return Field{val: val, isBlank: true}
	}
	return Field{val: val, String: strings.TrimSpace(val)}
}

// parseInteger parses the value into an integer Field.
func parseInteger(val string) Field {
	if val == "" {
		return Field{val: val, isBlank: true}
	}

	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return Field{val: val, err: errParseInt}
	}

	return Field{val: val, Integer: int(num)}
}

// parseFloat parses the value into a float64 Field.
func parseFloat(val string) Field {
	if val == "" {
		return Field{val: val, isBlank: true}
	}

	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return Field{val: val, err: errParseFloat}
	}

	return Field{val: val, Float: num}
}

// parseBool parses the value into a bool Field.
func parseBool(val string) Field {
	if val == "" {
		return Field{val: val, isBlank: true}
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		return Field{val: val, err: errParseBool}
	}

	return Field{val: val, Bool: b}
}

// parseDecimal parses the value into a decimal Field.
func parseDecimal(val string) Field {
	if val == "" {
		return Field{val: val, isBlank: true}
	}

	br, ok := new(big.Rat).SetString(val)
	if !ok {
		return Field{val: val, err: errParseBigRat}
	}

	return Field{val: val, Decimal: *br}
}

// parseDate parses the value into a date Field.
func parseDate(val string) Field {
	if val == "" {
		return Field{val: val, isBlank: true}
	}

	date, err := format.ParseDate(val)
	if err != nil {
		return Field{val: val, err: errParseDate}
	}

	return Field{val: val, Date: date}
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
		if v.IsBlank() {
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
		if utf8.RuneCountInString(v.String) < n {
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
func StrMatches(rx *regexp.Regexp, err error) CheckFunc {
	return func(v Field) error {
		if !rx.MatchString(v.String) {
			return err
		}
		return nil
	}
}

// Checks that a string is an email.
func StrEmail() CheckFunc {
	return StrMatches(emailRegex, errInvalidEmail)
}

// Checks that a string is a URL.
func StrUrl() CheckFunc {
	return StrMatches(urlRegex, errInvalidUrl)
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
		if v.IsBlank() {
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
		if v.IsBlank() {
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
		if v.IsBlank() {
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
		if v.IsBlank() {
			return errBlankValue
		}
		return nil
	}
}

// Checks that a decimal is less than n.
func DecLt(n string) CheckFunc {
	nn := parseDecimal(n)

	return func(v Field) error {
		if nn.Err() != nil {
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
	nn := parseDecimal(n)

	return func(v Field) error {
		if nn.Err() != nil {
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
	nn := parseDecimal(n)

	return func(v Field) error {
		if nn.Err() != nil {
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
	nn := parseDecimal(n)

	return func(v Field) error {
		if nn.Err() != nil {
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
	nn := parseDecimal(n)
	mm := parseDecimal(m)

	return func(v Field) error {
		if (nn.Err() != nil) || (mm.Err() != nil) {
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
		if v.IsBlank() {
			return errBlankValue
		}
		return nil
	}
}

// Checks that a date is before n (yyyy-mm-dd).
func DtBefore(n string) CheckFunc {
	nn := parseDate(n)

	return func(v Field) error {
		if nn.Err() != nil {
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
	nn := parseDate(n)

	return func(v Field) error {
		if nn.Err() != nil {
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
