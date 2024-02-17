package form

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
)

var (
	errInvalidChoice = errors.New("must be a valid choice")
	errInvalidConfig = errors.New("invalid validation config")
	errBlankValue    = errors.New("cannot be blank")
	errParseInt      = errors.New("must be an integer value")
	errParseDate     = errors.New("must be a date value")
	errParseBool     = errors.New("must be a boolean value")
	errParseBigRat   = errors.New("must be a decimal value")
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	urlRegex   = regexp.MustCompile(`^(http(s)?://)?([\da-z\.-]+)\.([a-z\.]{2,6})([/\w \.-]*)*/?$`)
)

// ------------------------------------------------------------------
//
//
// Custom Parsed Types
//
//
// ------------------------------------------------------------------

// A String represents a parsed string.
type String struct {
	Value string
	Blank bool
}

func (s String) Len() int {
	return utf8.RuneCountInString(s.Value)
}

// A Integer represents a parsed integer.
type Integer struct {
	Value int
	Blank bool
}

// A Bool represents a parsed boolean.
type Bool struct {
	Value bool
	Blank bool
}

// A Date represents a parsed date.
type Date struct {
	Value time.Time
	Blank bool
}

func (d Date) String() string {
	return d.Value.Format("2006-01-02")
}

// A Decimal represents a parsed decimal.
type Decimal struct {
	Value big.Rat
	Blank bool
}

func (d Decimal) String() string {
	return d.Value.FloatString(2)
}

// ------------------------------------------------------------------
//
//
// Validation Function Types
//
//
// ------------------------------------------------------------------

type StringFunc func(String) error

type IntegerFunc func(Integer) error

type BoolFunc func(Bool) error

type DateFunc func(Date) error

type DecimalFunc func(Decimal) error

// ------------------------------------------------------------------
//
//
// Validator
//
//
// ------------------------------------------------------------------

// Validator is a type which parses and validates form data.
// It is initialized with a data map, containing values to process.
// Values can be parsed into a desired type, like an int, date, bool, etc.
// Validation functions can also be provided, which ensures that
// the parsed type is properly vetted before being retrieved.
type Validator struct {
	data   map[string]string
	errors map[string]error
}

// NewValidator constructs a new Validator with the given data map.
// The provided data is normalized, which means
// keys and values are trimmed, and
// key/value pairs are omitted if the value is an empty string.
func NewValidator(data map[string]string) *Validator {
	v := Validator{
		data:   make(map[string]string, len(data)),
		errors: make(map[string]error),
	}
	// Trim spaces and omit empty values.
	for key, val := range data {
		kk := strings.TrimSpace(key)
		vv := strings.TrimSpace(val)
		if kk != "" && vv != "" {
			v.data[kk] = vv
		}
	}
	return &v
}

func (v *Validator) AddError(key string, err error) {
	if _, exists := v.errors[key]; !exists {
		v.errors[key] = err
	}
}

func (v *Validator) GetError(key string) error {
	return v.errors[key]
}

func (v *Validator) Errors() map[string]error {
	return v.errors
}

func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

// ------------------------------------------------------------------
//
//
// Form Parse Functions
//
//
// ------------------------------------------------------------------

// Convert the field into a string and perform validations.
func (v *Validator) String(key string, funcs ...StringFunc) string {
	// Read the value.
	w, err := NewString(v.data[key])
	if err != nil {
		v.AddError(key, err)
		return w.Value
	}

	// Validate the value.
	for i := range funcs {
		if err := funcs[i](w); err != nil {
			v.AddError(key, err)
			return w.Value
		}
	}
	return w.Value
}

// Convert the field into an integer and perform validations.
func (v *Validator) Integer(key string, funcs ...IntegerFunc) int {
	// Read the value.
	w, err := NewInteger(v.data[key])
	if err != nil {
		v.AddError(key, err)
		return w.Value
	}

	// Validate the value.
	for i := range funcs {
		if err := funcs[i](w); err != nil {
			v.AddError(key, err)
			return w.Value
		}
	}
	return w.Value
}

// Convert the field into a date and perform validations.
func (v *Validator) Date(key string, funcs ...DateFunc) time.Time {
	// Read the value.
	w, err := NewDate(v.data[key])
	if err != nil {
		v.AddError(key, err)
		return w.Value
	}

	// Validate the value.
	for i := range funcs {
		if err := funcs[i](w); err != nil {
			v.AddError(key, err)
			return w.Value
		}
	}
	return w.Value
}

// Convert the field into a bool and perform validations.
func (v *Validator) Bool(key string, funcs ...BoolFunc) bool {
	// Read the value.
	w, err := NewBool(v.data[key])
	if err != nil {
		v.AddError(key, err)
		return w.Value
	}

	// Validate the value.
	for i := range funcs {
		if err := funcs[i](w); err != nil {
			v.AddError(key, err)
			return w.Value
		}
	}
	return w.Value
}

// Convert the field into a decimal and perform validations.
func (v *Validator) Decimal(key string, funcs ...DecimalFunc) big.Rat {
	// Read the value.
	w, err := NewDecimal(v.data[key])
	if err != nil {
		v.AddError(key, err)
		return w.Value
	}

	// Validate the value.
	for i := range funcs {
		if err := funcs[i](w); err != nil {
			v.AddError(key, err)
			return w.Value
		}
	}
	return w.Value
}

// ------------------------------------------------------------------
//
//
// Read Functions
//
//
// ------------------------------------------------------------------

// Parse the given value into a custom String type.
func NewString(val string) (String, error) {
	if val == "" {
		return String{Blank: true}, nil
	}
	return String{Value: strings.TrimSpace(val)}, nil
}

// Parse the given value into a custom Integer type.
func NewInteger(val string) (Integer, error) {
	if val == "" {
		return Integer{Blank: true}, nil
	}

	// Parse the string into an int.-
	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return Integer{}, errParseInt
	}

	return Integer{Value: int(num)}, nil
}

// Parse the given value into a custom Bool type.
func NewBool(val string) (Bool, error) {
	if val == "" {
		return Bool{Blank: true}, nil
	}

	// Parse the string into a bool.
	b, err := strconv.ParseBool(val)
	if err != nil {
		return Bool{}, errParseBool
	}

	return Bool{Value: b}, nil
}

// Parse the given value into a custom Date type.
func NewDate(val string) (Date, error) {
	if val == "" {
		return Date{Blank: true}, nil
	}

	// Parse the string into a date.
	date, err := time.Parse("2006-01-02", val)
	if err != nil {
		return Date{}, errParseDate
	}

	return Date{Value: date}, nil
}

// Parse the given value into a custom Decimal type.
func NewDecimal(val string) (Decimal, error) {
	if val == "" {
		return Decimal{Blank: true}, nil
	}

	// Parse the string into a big.Rat.
	br, ok := new(big.Rat).SetString(val)
	if !ok {
		return Decimal{}, errParseBigRat
	}

	return Decimal{Value: *br}, nil
}

// ------------------------------------------------------------------
//
//
// Integer Check Functions
//
//
// ------------------------------------------------------------------

// Checks that an int is not blank.
func IntRequired() IntegerFunc {
	return func(i Integer) error {
		if i.Blank {
			return errBlankValue
		}
		return nil
	}
}

// Checks that an int is less than n.
func IntLt(n int) IntegerFunc {
	err := fmt.Errorf("must be less than %d", n)

	return func(i Integer) error {
		if i.Value >= n {
			return err
		}
		return nil
	}
}

// Checks that an int is less than or equal to n.
func IntLte(n int) IntegerFunc {
	err := fmt.Errorf("must be less than or equal to %d", n)

	return func(i Integer) error {
		if i.Value > n {
			return err
		}
		return nil
	}
}

// Checks that an int is more than n.
func IntGt(n int) IntegerFunc {
	err := fmt.Errorf("must be more than %d", n)

	return func(i Integer) error {
		if i.Value <= n {
			return err
		}
		return nil
	}
}

// Checks that an int is more than or equal to n.
func IntGte(n int) IntegerFunc {
	err := fmt.Errorf("must be more than or equal to %d", n)

	return func(i Integer) error {
		if i.Value < n {
			return err
		}
		return nil
	}
}

// Checks that an int is between n and m.
func IntBtw(n, m int) IntegerFunc {
	err := fmt.Errorf("must be between %d and %d", n, m)

	return func(i Integer) error {
		if (n == m) || (n > m) {
			return errInvalidConfig
		}
		if (i.Value < n) || (i.Value > m) {
			return err
		}
		return nil
	}
}

// Checks that an int is a member of the given choices.
func IntIn(choices []int) IntegerFunc {
	return func(i Integer) error {
		if !slices.Contains(choices, i.Value) {
			return errInvalidChoice
		}
		return nil
	}
}

// ------------------------------------------------------------------
//
//
// String Check Functions
//
//
// ------------------------------------------------------------------

// Checks that a string is not blank.
func StrRequired() StringFunc {
	return func(s String) error {
		if s.Blank {
			return errBlankValue
		}
		return nil
	}
}

// Sets a default value for the string.
func StrDefault(n string) StringFunc {
	return func(s String) error {
		if s.Blank {
			s.Value = n
		}
		return nil
	}
}

// Checks that a string's length is less than n.
func StrLt(n int) StringFunc {
	err := fmt.Errorf("must be less than %d characters", n)

	return func(s String) error {
		if len(s.Value) >= n {
			return err
		}
		return nil
	}
}

// Checks that a string's length is less than or equal to n.
func StrLte(n int) StringFunc {
	err := fmt.Errorf("must be less than or equal to %d characters", n)

	return func(s String) error {
		if utf8.RuneCountInString(s.Value) > n {
			return err
		}
		return nil
	}
}

// Checks that a string's length is greater than n.
func StrGt(n int) StringFunc {
	err := fmt.Errorf("must be more than %d characters", n)

	return func(s String) error {
		if utf8.RuneCountInString(s.Value) <= n {
			return err
		}
		return nil
	}
}

// Checks that a string's length is greater than or equal to n.
func StrGte(n int) StringFunc {
	err := fmt.Errorf("must be more than or equal to %d characters", n)

	return func(s String) error {
		if utf8.RuneCountInString(s.Value) <= n {
			return err
		}
		return nil
	}
}

// Checks that a string's length is between n and m.
func StrBtw(n, m int) StringFunc {
	err := fmt.Errorf("must be between %d and %d characters", n, m)

	return func(s String) error {
		if (n == m) || (n > m) {
			return errInvalidConfig
		}
		if strLen := utf8.RuneCountInString(s.Value); strLen < n || strLen > m {
			return err
		}
		return nil
	}
}

// Checks that a string is a member of the given choices.
func StrIn(choices []string) StringFunc {
	return func(s String) error {
		if !slices.Contains(choices, s.Value) {
			return errInvalidChoice
		}
		return nil
	}
}

// Checks that a string matches the given regex.
func StrMatches(rx *regexp.Regexp, errMsg string) StringFunc {
	err := fmt.Errorf(errMsg)

	return func(s String) error {
		if !rx.MatchString(s.Value) {
			return err
		}
		return nil
	}
}

// Checks that a string is an email.
func StrEmail() StringFunc {
	return StrMatches(emailRegex, "must be a valid email")
}

// Checks that a string is a URL.
func StrUrl() StringFunc {
	return StrMatches(urlRegex, "must be a valid url")
}

// ------------------------------------------------------------------
//
//
// Date Check Functions
//
//
// ------------------------------------------------------------------

// Checks that a date is not blank.
func DtRequired() DateFunc {
	return func(d Date) error {
		if d.Blank {
			return errBlankValue
		}
		return nil
	}
}

// Checks that a date is before n (yyyy-mm-dd).
func DtBefore(n string) DateFunc {
	t, terr := NewDate(n)
	err := fmt.Errorf("must be before %s", t)

	return func(d Date) error {
		if terr != nil {
			return errInvalidConfig
		}
		if !d.Value.Before(t.Value) {
			return err
		}
		return nil
	}
}

// Checks that a date is after n (yyyy-mm-dd).
func DtAfter(n string) DateFunc {
	t, terr := NewDate(n)
	err := fmt.Errorf("must be after %s", t)

	return func(d Date) error {
		if terr != nil {
			return errInvalidConfig
		}
		if !d.Value.After(t.Value) {
			return err
		}
		return nil
	}
}

// Checks that a date is in the past.
func DtInPast() DateFunc {
	err := fmt.Errorf("must be a past date")

	return func(d Date) error {
		if !d.Value.Before(time.Now()) {
			return err
		}
		return nil
	}
}

// Checks that a date is in the future.
func DtInFuture() DateFunc {
	err := fmt.Errorf("must be a future date")

	return func(d Date) error {
		if !d.Value.After(time.Now()) {
			return err
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
func BoolRequired() BoolFunc {
	return func(b Bool) error {
		if b.Blank {
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
func DecRequired() DecimalFunc {
	return func(d Decimal) error {
		if d.Blank {
			return errBlankValue
		}
		return nil
	}
}

// Checks that a decimal is less than n.
func DecLt(n string) DecimalFunc {
	nn, nerr := NewDecimal(n)
	err := fmt.Errorf("must be less than %s", nn)

	return func(d Decimal) error {
		if nerr != nil {
			return errInvalidConfig
		}
		// Check if d >= nn.
		if r := d.Value.Cmp(&nn.Value); (r == 0) || (r == 1) {
			return err
		}
		return nil
	}
}

// Checks that a decimal is less than or equal to n.
func DecLte(n string) DecimalFunc {
	nn, nerr := NewDecimal(n)
	err := fmt.Errorf("must be less than or equal to %s", nn)

	return func(d Decimal) error {
		if nerr != nil {
			return errInvalidConfig
		}
		// Check if d > nn.
		if r := d.Value.Cmp(&nn.Value); r == 1 {
			return err
		}
		return nil
	}
}

// Checks that a decimal is more than n.
func DecGt(n string) DecimalFunc {
	nn, nerr := NewDecimal(n)
	err := fmt.Errorf("must be more than %s", nn)

	return func(d Decimal) error {
		if nerr != nil {
			return errInvalidConfig
		}
		// Check if d <= nn.
		if r := d.Value.Cmp(&nn.Value); (r == -1) || (r == 0) {
			return err
		}
		return nil
	}
}

// Checks that a decimal is greater than or equal to n.
func DecGte(n string) DecimalFunc {
	nn, nerr := NewDecimal(n)
	err := fmt.Errorf("must be more than or equal to %s", nn)

	return func(d Decimal) error {
		if nerr != nil {
			return errInvalidConfig
		}
		// Check if d < nn.
		if r := d.Value.Cmp(&nn.Value); r == -1 {
			return err
		}
		return nil
	}
}

// Checks that a decimal is between n and m.
func DecBtw(n, m string) DecimalFunc {
	nn, nerr := NewDecimal(n)
	mm, merr := NewDecimal(m)
	err := fmt.Errorf("must be between %s and %s", nn, mm)

	return func(d Decimal) error {
		if (nerr != nil) || (merr != nil) {
			return errInvalidConfig
		}

		// Check if (n == m) || (n > m)
		r := nn.Value.Cmp(&mm.Value)
		if (r == 0) || (r == 1) {
			return errInvalidConfig
		}

		// Check if (d < n) || (d > m)
		nr := d.Value.Cmp(&nn.Value)
		mr := d.Value.Cmp(&mm.Value)
		if (nr == -1) || (mr == 1) {
			return err
		}
		return nil
	}
}
