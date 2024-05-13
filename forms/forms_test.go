package forms

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"testing"
	"time"
)

// ------------------------------------------------------------------
// ------------------------------------------------------------------
//
//
//
// Form
//
//
//
// ------------------------------------------------------------------
// ------------------------------------------------------------------

func TestForm(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		form := New()
		assert(t, len(form.fields), 0)
		assert(t, len(form.extraerrors), 0)
	})

	testCheck := func(idx *int) CheckFunc {
		return func(v Field) error {
			*idx++
			return nil
		}
	}

	okCheck := func() CheckFunc {
		return func(v Field) error {
			return nil
		}
	}

	failCheck := func() CheckFunc {
		return func(v Field) error {
			return errors.New("boom")
		}
	}

	t.Run("CleanString", func(t *testing.T) {
		// Arrange
		form := New()
		var idx int
		assert(t, idx, 0)

		// Act
		form.CleanString("testStr", "test", testCheck(&idx))

		// Assert
		assert(t, form.CleanedString("testStr"), "test")
		assert(t, idx, 1)
	})

	t.Run("CleanInteger", func(t *testing.T) {
		// Arrange
		form := New()
		var idx int
		assert(t, idx, 0)

		// Act
		form.CleanInteger("testInt", "53", testCheck(&idx))

		// Assert
		assert(t, form.CleanedInteger("testInt"), 53)
		assert(t, idx, 1)
	})

	t.Run("CleanFloat", func(t *testing.T) {
		// Arrange
		form := New()
		var idx int
		assert(t, idx, 0)

		// Act
		form.CleanFloat("testFloat", "53.2", testCheck(&idx))

		// Assert
		assert(t, form.CleanedFloat("testFloat"), 53.2)
		assert(t, idx, 1)
	})

	t.Run("CleanBool", func(t *testing.T) {
		// Arrange
		form := New()
		var idx int
		assert(t, idx, 0)

		// Act
		form.CleanBool("testBool", "true", testCheck(&idx))

		// Assert
		assert(t, form.CleanedBool("testBool"), true)
		assert(t, idx, 1)
	})

	t.Run("CleanDate", func(t *testing.T) {
		// Arrange
		form := New()
		var idx int
		assert(t, idx, 0)

		// Act
		form.CleanDate("testDate", "2020-07-23", testCheck(&idx))

		// Assert
		assert(t, form.CleanedDate("testDate"), time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC))
		assert(t, idx, 1)
	})

	t.Run("CleanDecimal", func(t *testing.T) {
		// Arrange
		form := New()
		var idx int
		assert(t, idx, 0)

		// Act
		form.CleanDecimal("testDecimal", "57.35", testCheck(&idx))

		// Assert
		decimal := form.CleanedDecimal("testDecimal")
		assert(t, decimal.RatString(), "1147/20")
		assert(t, idx, 1)
	})

	t.Run("CleanExtra", func(t *testing.T) {
		// Arrange
		form := New()

		// Act
		form.CleanExtra(true, errors.New("boom"))

		// Assert
		extraErrors := form.ExtraErrors()
		assert(t, len(extraErrors), 1)
		assert(t, extraErrors[0].Error(), "boom")
	})

	t.Run("cleanField-ok", func(t *testing.T) {
		// Arrange
		form := New()

		// Act
		form.cleanField("testInteger", parseInteger("12"), okCheck())

		// Assert
		field := form.MustField("testInteger")
		assert(t, field.Integer, 12)
		assert(t, field.Err(), nil)
		assert(t, len(form.fields), 1)
	})

	t.Run("cleanField-parse-error", func(t *testing.T) {
		// Arrange
		form := New()

		// Act
		form.cleanField("testInteger", parseInteger("bad"), okCheck())

		// Assert
		field := form.MustField("testInteger")
		assert(t, field.Err().Error(), "must be a valid integer")
		assert(t, len(form.fields), 1)
	})

	t.Run("cleanField-validation-error", func(t *testing.T) {
		// Arrange
		form := New()

		// Act
		form.cleanField("testInteger", parseInteger("12"), okCheck(), failCheck())

		// Assert
		field := form.MustField("testInteger")
		assert(t, field.Integer, 12)
		assert(t, field.Err().Error(), "boom")
		assert(t, len(form.fields), 1)
	})

	t.Run("addField", func(t *testing.T) {
		// Arrange
		form := New()
		assert(t, len(form.fields), 0)
		field := Field{val: "a"}

		// Act
		form.addField("test", field)

		// Assert
		assert(t, len(form.fields), 1)
		assert(t, form.fields["test"].val, field.val)
	})

	t.Run("addField-does-not-override-existing-field", func(t *testing.T) {
		// Arrange
		form := New()
		assert(t, len(form.fields), 0)
		fieldA := Field{val: "a"}
		fieldB := Field{val: "b"}

		// Act
		form.addField("test", fieldA)
		form.addField("test", fieldB)

		// Assert
		assert(t, len(form.fields), 1)
		assert(t, form.fields["test"].val, fieldA.val)
	})

	t.Run("Field-exists", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanString("testStr", "test")

		// Act
		field, ok := form.Field("testStr")

		// Assert
		assert(t, ok, true)
		assert(t, field.Value(), "test")
		assert(t, field.String, "test")
	})

	t.Run("Field-does-not-exist", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanString("testStr", "test")

		// Act
		field, ok := form.Field("notFound")

		// Assert
		assert(t, ok, false)
		assert(t, field.Value(), "")
		assert(t, field.String, "")
	})

	t.Run("Field-nil", func(t *testing.T) {
		// Arrange
		form := New()

		// Act
		field, ok := form.Field("notFound")

		// Assert
		assert(t, ok, false)
		assert(t, field.Value(), "")
		assert(t, field.String, "")
	})

	t.Run("Names", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanString("fieldA", "test-a")
		form.CleanString("fieldB", "test-b")
		form.CleanString("fieldC", "test-c")

		// Act
		names := form.Names()

		// Assert
		assert(t, len(names), 3)
		assert(t, slices.Contains(names, "fieldA"), true)
		assert(t, slices.Contains(names, "fieldB"), true)
		assert(t, slices.Contains(names, "fieldC"), true)
	})

	t.Run("MustField-exists", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanString("fieldA", "test-a")

		// Act
		field := form.MustField("fieldA")

		// Assert
		assert(t, field.Value(), "test-a")
	})

	t.Run("MustField-does-not-exist", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanString("fieldA", "test-a")

		// Act / Assert
		assertPanic(t, func() { form.MustField("notFound") })
	})
}

// ------------------------------------------------------------------
// ------------------------------------------------------------------
//
//
//
// Field
//
//
//
// ------------------------------------------------------------------
// ------------------------------------------------------------------

func TestField(t *testing.T) {
	t.Run("getters", func(t *testing.T) {
		var f Field

		// Value()
		f = Field{val: "1.2"}
		assert(t, f.Value(), "1.2")

		// IsBlank()
		f = Field{isBlank: true}
		assert(t, f.IsBlank(), true)

		// Err()
		f = Field{err: errors.New("test error")}
		assert(t, f.Err().Error(), "test error")
	})

	t.Run("add-error", func(t *testing.T) {
		f := Field{}
		assert(t, f.Err(), nil)

		f.addError(errors.New("test error"))
		assert(t, f.Err().Error(), "test error")
	})

	t.Run("add-error-does-not-overwrite-existing-error", func(t *testing.T) {
		f := Field{err: errors.New("test error")}
		assert(t, f.Err().Error(), "test error")

		f.addError(errors.New("another error"))
		assert(t, f.Err().Error(), "test error")
	})
}

// ------------------------------------------------------------------
// ------------------------------------------------------------------
//
//
//
// ParseError
//
//
//
// ------------------------------------------------------------------
// ------------------------------------------------------------------

func TestParseError(t *testing.T) {
	t.Run("string-and-error", func(t *testing.T) {
		err := ParseError{Msg: "test error"}
		assert(t, err.Error(), "test error")
		assert(t, err.String(), "test error")
	})
}

// ------------------------------------------------------------------
// ------------------------------------------------------------------
//
//
//
// Parse functions
//
//
//
// ------------------------------------------------------------------
// ------------------------------------------------------------------

func TestParseField(t *testing.T) {
	t.Run("parseString", func(t *testing.T) {
		var f Field

		// blank
		f = parseString("")
		assert(t, f.IsBlank(), true)
		assert(t, f.Value(), "")

		// success
		f = parseString("test")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "test")
		assert(t, f.String, "test")

		// success with trimmed string
		f = parseString(" test ")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), " test ")
		assert(t, f.String, "test")
	})

	t.Run("parseInteger", func(t *testing.T) {
		var f Field

		// blank
		f = parseInteger("")
		assert(t, f.IsBlank(), true)
		assert(t, f.Value(), "")

		// error
		f = parseInteger("bad")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "bad")
		assert(t, f.Err().Error(), "must be a valid integer")

		// success
		f = parseInteger("87")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "87")
		assert(t, f.Integer, 87)

		// success
		f = parseInteger("0")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "0")
		assert(t, f.Integer, 0)
	})

	t.Run("parseFloat", func(t *testing.T) {
		var f Field

		// blank
		f = parseFloat("")
		assert(t, f.IsBlank(), true)
		assert(t, f.Value(), "")

		// error
		f = parseFloat("bad")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "bad")
		assert(t, f.Err().Error(), "must be a valid float")

		// success
		f = parseFloat("87.12")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "87.12")
		assert(t, f.Float, 87.12)

		// success
		f = parseFloat("0.0")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "0.0")
		assert(t, f.Float, 0.0)
	})

	t.Run("parseBool", func(t *testing.T) {
		var f Field

		// blank
		f = parseBool("")
		assert(t, f.IsBlank(), true)
		assert(t, f.Value(), "")

		// error
		f = parseBool("bad")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "bad")
		assert(t, f.Err().Error(), "must be a valid boolean")

		// success
		f = parseBool("true")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "true")
		assert(t, f.Bool, true)

		// success
		f = parseBool("false")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "false")
		assert(t, f.Bool, false)
	})

	t.Run("parseDecimal", func(t *testing.T) {
		var f Field

		// blank
		f = parseDecimal("")
		assert(t, f.IsBlank(), true)
		assert(t, f.Value(), "")

		// error
		f = parseDecimal("bad")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "bad")
		assert(t, f.Err().Error(), "must be a valid decimal")

		// success
		f = parseDecimal("87.13")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "87.13")
		assert(t, f.Decimal.RatString(), "8713/100")

		// success
		f = parseDecimal("0.0")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "0.0")
		assert(t, f.Decimal.RatString(), "0")
	})

	t.Run("parseDate", func(t *testing.T) {
		var f Field

		// blank
		f = parseDate("")
		assert(t, f.IsBlank(), true)
		assert(t, f.Value(), "")

		// error
		f = parseDate("bad")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "bad")
		assert(t, f.Err().Error(), "must be a valid date")

		// --------------------------------------
		// various date layouts
		// --------------------------------------

		// layout "2006-01-02"
		f = parseDate("2020-07-23")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "2020-07-23")
		assert(t, f.Date, time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC))

		// layout "January-2-2006"
		f = parseDate("July-23-2020")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "July-23-2020")
		assert(t, f.Date, time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC))

		// layout "2-January-2006"
		f = parseDate("23-July-2020")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "23-July-2020")
		assert(t, f.Date, time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC))

		// layout "January 2, 2006"
		f = parseDate("July 23, 2020")
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "July 23, 2020")
		assert(t, f.Date, time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC))

		// --------------------------------------
		// lowercase date input
		// --------------------------------------

		// layout "January-2-2006"
		f = parseDate("july-23-2020") // lowercase
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "july-23-2020")
		assert(t, f.Date, time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC))

		// layout "2-January-2006"
		f = parseDate("23-july-2020") // lowercase
		assert(t, f.IsBlank(), false)
		assert(t, f.Value(), "23-july-2020")
		assert(t, f.Date, time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC))
	})
}

// ------------------------------------------------------------------
// ------------------------------------------------------------------
//
//
//
// String CheckFuncs
//
//
//
// ------------------------------------------------------------------
// ------------------------------------------------------------------

func TestStringCheckFuncs(t *testing.T) {
	t.Run("StrRequired-ok", func(t *testing.T) {
		field := parseString("test")
		err := StrRequired()(field)
		assert(t, err, nil)
	})

	t.Run("StrRequired-error", func(t *testing.T) {
		field := parseString("")
		err := StrRequired()(field)
		assert(t, err, errBlankValue)
	})

	t.Run("StrLt-ok", func(t *testing.T) {
		field := parseString("abcd-abcd")
		err := StrLt(10)(field)
		assert(t, err, nil)
	})

	t.Run("StrLt-error", func(t *testing.T) {
		field := parseString("abcd-abcd")
		err := StrLt(5)(field)
		assert(t, err.Error(), "must be less than 5 characters")
	})
}

// ------------------------------------------------------------------
// ------------------------------------------------------------------
//
//
//
// Helpers
//
//
//
// ------------------------------------------------------------------
// ------------------------------------------------------------------

// assert checks that value a is equal to expected.
func assert(t *testing.T, a interface{}, expected interface{}) {
	if a == expected {
		return
	}

	_, filename, line, _ := runtime.Caller(1)
	msg := "%s:%d expected %v (type %v), got %v (type %v)\n"
	fmt.Printf(msg, filepath.Base(filename), line, expected, reflect.TypeOf(expected), a, reflect.TypeOf(a))
	t.FailNow()
}

// assertPanic checks if the given function f panics.
func assertPanic(t *testing.T, f func()) {
	_, filename, line, _ := runtime.Caller(1)

	defer func() {
		if r := recover(); r == nil {
			msg := "%s:%d the code did not panic\n"
			fmt.Printf(msg, filepath.Base(filename), line)
			t.FailNow()
		}
	}()
	f() // function that panics
}
