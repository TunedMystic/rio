package forms

import (
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/tunedmystic/rio/internal/assert"
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
		assert.Equal(t, len(form.fields), 0)
		assert.Equal(t, len(form.extraerrors), 0)
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
		assert.Equal(t, idx, 0)

		// Act
		form.CleanString("testStr", "test", testCheck(&idx))

		// Assert
		assert.Equal(t, form.CleanedString("testStr"), "test")
		assert.Equal(t, idx, 1)
	})

	t.Run("CleanInteger", func(t *testing.T) {
		// Arrange
		form := New()
		var idx int
		assert.Equal(t, idx, 0)

		// Act
		form.CleanInteger("testInt", "53", testCheck(&idx))

		// Assert
		assert.Equal(t, form.CleanedInteger("testInt"), 53)
		assert.Equal(t, idx, 1)
	})

	t.Run("CleanFloat", func(t *testing.T) {
		// Arrange
		form := New()
		var idx int
		assert.Equal(t, idx, 0)

		// Act
		form.CleanFloat("testFloat", "53.2", testCheck(&idx))

		// Assert
		assert.Equal(t, form.CleanedFloat("testFloat"), 53.2)
		assert.Equal(t, idx, 1)
	})

	t.Run("CleanBool", func(t *testing.T) {
		// Arrange
		form := New()
		var idx int
		assert.Equal(t, idx, 0)

		// Act
		form.CleanBool("testBool", "true", testCheck(&idx))

		// Assert
		assert.Equal(t, form.CleanedBool("testBool"), true)
		assert.Equal(t, idx, 1)
	})

	t.Run("CleanDate", func(t *testing.T) {
		// Arrange
		form := New()
		var idx int
		assert.Equal(t, idx, 0)

		// Act
		form.CleanDate("testDate", "2020-07-23", testCheck(&idx))

		// Assert
		assert.Equal(t, form.CleanedDate("testDate"), time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC))
		assert.Equal(t, idx, 1)
	})

	t.Run("CleanDecimal", func(t *testing.T) {
		// Arrange
		form := New()
		var idx int
		assert.Equal(t, idx, 0)

		// Act
		form.CleanDecimal("testDecimal", "57.35", testCheck(&idx))

		// Assert
		decimal := form.CleanedDecimal("testDecimal")
		assert.Equal(t, decimal.RatString(), "1147/20")
		assert.Equal(t, idx, 1)
	})

	t.Run("CleanExtra", func(t *testing.T) {
		// Arrange
		form := New()

		// Act
		form.CleanExtra(true, errors.New("boom"))

		// Assert
		extraErrors := form.ExtraErrors()
		assert.Equal(t, len(extraErrors), 1)
		assert.Equal(t, extraErrors[0].Error(), "boom")
	})

	t.Run("cleanField-ok", func(t *testing.T) {
		// Arrange
		form := New()

		// Act
		form.cleanField("testInteger", parseInteger("12"), okCheck())

		// Assert
		field := form.MustField("testInteger")
		assert.Equal(t, field.Integer, 12)
		assert.Equal(t, field.Err(), nil)
		assert.Equal(t, len(form.fields), 1)
	})

	t.Run("cleanField-parse-error", func(t *testing.T) {
		// Arrange
		form := New()

		// Act
		form.cleanField("testInteger", parseInteger("bad"), okCheck())

		// Assert
		field := form.MustField("testInteger")
		assert.Equal(t, field.Err().Error(), "must be a valid integer")
		assert.Equal(t, len(form.fields), 1)
	})

	t.Run("cleanField-validation-error", func(t *testing.T) {
		// Arrange
		form := New()

		// Act
		form.cleanField("testInteger", parseInteger("12"), okCheck(), failCheck())

		// Assert
		field := form.MustField("testInteger")
		assert.Equal(t, field.Integer, 12)
		assert.Equal(t, field.Err().Error(), "boom")
		assert.Equal(t, len(form.fields), 1)
	})

	t.Run("addField", func(t *testing.T) {
		// Arrange
		form := New()
		assert.Equal(t, len(form.fields), 0)
		field := Field{val: "a"}

		// Act
		form.addField("test", field)

		// Assert
		assert.Equal(t, len(form.fields), 1)
		assert.Equal(t, form.fields["test"].val, field.val)
	})

	t.Run("addField-does-not-override-existing-field", func(t *testing.T) {
		// Arrange
		form := New()
		assert.Equal(t, len(form.fields), 0)
		fieldA := Field{val: "a"}
		fieldB := Field{val: "b"}

		// Act
		form.addField("test", fieldA)
		form.addField("test", fieldB)

		// Assert
		assert.Equal(t, len(form.fields), 1)
		assert.Equal(t, form.fields["test"].val, fieldA.val)
	})

	t.Run("Field-exists", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanString("testStr", "test")

		// Act
		field, ok := form.Field("testStr")

		// Assert
		assert.Equal(t, ok, true)
		assert.Equal(t, field.Value(), "test")
		assert.Equal(t, field.String, "test")
	})

	t.Run("Field-does-not-exist", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanString("testStr", "test")

		// Act
		field, ok := form.Field("notFound")

		// Assert
		assert.Equal(t, ok, false)
		assert.Equal(t, field.Value(), "")
		assert.Equal(t, field.String, "")
	})

	t.Run("Field-nil", func(t *testing.T) {
		// Arrange
		form := New()

		// Act
		field, ok := form.Field("notFound")

		// Assert
		assert.Equal(t, ok, false)
		assert.Equal(t, field.Value(), "")
		assert.Equal(t, field.String, "")
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
		assert.Equal(t, len(names), 3)
		assert.Equal(t, slices.Contains(names, "fieldA"), true)
		assert.Equal(t, slices.Contains(names, "fieldB"), true)
		assert.Equal(t, slices.Contains(names, "fieldC"), true)
	})

	t.Run("MustField-exists", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanString("fieldA", "test-a")

		// Act
		field := form.MustField("fieldA")

		// Assert
		assert.Equal(t, field.Value(), "test-a")
	})

	t.Run("MustField-does-not-exist", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanString("fieldA", "test-a")

		// Act / Assert
		assert.Panic(t, func() { form.MustField("fieldB") })
	})

	t.Run("IsValid", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanString("stringField", "test", StrRequired())
		form.CleanInteger("integerField", "1", IntRequired())
		form.CleanFloat("floatField", "1.23", FltRequired())
		form.CleanBool("boolField", "true", BoolRequired())
		form.CleanDate("dateField", "2020-07-23", DtRequired())
		form.CleanDecimal("decimalField", "11.34", DecRequired())

		// Act / Assert
		assert.Equal(t, form.IsValid(), true)
		for _, name := range form.Names() {
			fmt.Println(name, form.MustField(name).Err())
		}
	})

	t.Run("IsValid-false-fielderror", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanString("stringField", "test", StrRequired())
		form.CleanInteger("integerField", "1", IntRequired())
		form.CleanFloat("floatField", "", FltRequired()) // will cause a field error
		form.CleanBool("boolField", "true", BoolRequired())
		form.CleanDate("dateField", "2020-07-23", DtRequired())
		form.CleanDecimal("decimalField", "11.34", DecRequired())

		// Act / Assert
		assert.Equal(t, form.IsValid(), false)
		assert.Equal(t, len(form.extraerrors), 0)
	})

	t.Run("IsValid-false-non-fielderror", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanString("stringField", "test", StrRequired())
		form.CleanInteger("integerField", "1", IntRequired())
		form.CleanFloat("floatField", "1.23", FltRequired())
		form.CleanBool("boolField", "true", BoolRequired())
		form.CleanDate("dateField", "2020-07-23", DtRequired())
		form.CleanDecimal("decimalField", "11.34", DecRequired())
		form.CleanExtra(true, errors.New("boom")) // will cause a non-field error

		// Act / Assert
		assert.Equal(t, form.IsValid(), false)
		assert.Equal(t, len(form.extraerrors), 1)
	})

	t.Run("HasError-no-error", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanFloat("floatField", "1.23")

		// Act / Assert
		var pErr ParseError
		assert.Equal(t, form.HasError(&pErr), false)
	})

	t.Run("HasError-fielderror", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanFloat("floatField", "x.23")

		// Act / Assert
		var pErr ParseError
		assert.Equal(t, form.HasError(&pErr), true)
	})

	t.Run("HasError-non-fielderror", func(t *testing.T) {
		// Arrange
		form := New()
		form.CleanFloat("floatField", "1.23")
		form.CleanExtra(true, TestError{Msg: "boom"})

		// Act / Assert
		var tErr TestError
		assert.Equal(t, form.HasError(&tErr), true)
	})
}

type TestError struct {
	Msg string
}

func (t TestError) Error() string {
	return t.Msg
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
		assert.Equal(t, f.Value(), "1.2")

		// IsBlank()
		f = Field{isBlank: true}
		assert.Equal(t, f.IsBlank(), true)

		// Err()
		f = Field{err: errors.New("test error")}
		assert.Equal(t, f.Err().Error(), "test error")
	})

	t.Run("add-error", func(t *testing.T) {
		f := Field{}
		assert.Equal(t, f.Err(), nil)

		f.addError(errors.New("test error"))
		assert.Equal(t, f.Err().Error(), "test error")
	})

	t.Run("add-error-does-not-overwrite-existing-error", func(t *testing.T) {
		f := Field{err: errors.New("test error")}
		assert.Equal(t, f.Err().Error(), "test error")

		f.addError(errors.New("another error"))
		assert.Equal(t, f.Err().Error(), "test error")
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
		assert.Equal(t, err.Error(), "test error")
		assert.Equal(t, err.String(), "test error")
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
		assert.Equal(t, f.IsBlank(), true)
		assert.Equal(t, f.Value(), "")

		// success
		f = parseString("test")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "test")
		assert.Equal(t, f.String, "test")

		// success with trimmed string
		f = parseString(" test ")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), " test ")
		assert.Equal(t, f.String, "test")
	})

	t.Run("parseInteger", func(t *testing.T) {
		var f Field

		// blank
		f = parseInteger("")
		assert.Equal(t, f.IsBlank(), true)
		assert.Equal(t, f.Value(), "")

		// error
		f = parseInteger("bad")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "bad")
		assert.Equal(t, f.Err().Error(), "must be a valid integer")

		// success
		f = parseInteger("87")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "87")
		assert.Equal(t, f.Integer, 87)

		// success
		f = parseInteger("0")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "0")
		assert.Equal(t, f.Integer, 0)
	})

	t.Run("parseFloat", func(t *testing.T) {
		var f Field

		// blank
		f = parseFloat("")
		assert.Equal(t, f.IsBlank(), true)
		assert.Equal(t, f.Value(), "")

		// error
		f = parseFloat("bad")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "bad")
		assert.Equal(t, f.Err().Error(), "must be a valid float")

		// success
		f = parseFloat("87.12")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "87.12")
		assert.Equal(t, f.Float, 87.12)

		// success
		f = parseFloat("0.0")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "0.0")
		assert.Equal(t, f.Float, 0.0)
	})

	t.Run("parseBool", func(t *testing.T) {
		var f Field

		// blank
		f = parseBool("")
		assert.Equal(t, f.IsBlank(), true)
		assert.Equal(t, f.Value(), "")

		// error
		f = parseBool("bad")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "bad")
		assert.Equal(t, f.Err().Error(), "must be a valid boolean")

		// success
		f = parseBool("true")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "true")
		assert.Equal(t, f.Bool, true)

		// success
		f = parseBool("false")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "false")
		assert.Equal(t, f.Bool, false)
	})

	t.Run("parseDecimal", func(t *testing.T) {
		var f Field

		// blank
		f = parseDecimal("")
		assert.Equal(t, f.IsBlank(), true)
		assert.Equal(t, f.Value(), "")

		// error
		f = parseDecimal("bad")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "bad")
		assert.Equal(t, f.Err().Error(), "must be a valid decimal")

		// success
		f = parseDecimal("87.13")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "87.13")
		assert.Equal(t, f.Decimal.RatString(), "8713/100")

		// success
		f = parseDecimal("0.0")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "0.0")
		assert.Equal(t, f.Decimal.RatString(), "0")
	})

	t.Run("parseDate", func(t *testing.T) {
		var f Field

		// blank
		f = parseDate("")
		assert.Equal(t, f.IsBlank(), true)
		assert.Equal(t, f.Value(), "")

		// error
		f = parseDate("bad")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "bad")
		assert.Equal(t, f.Err().Error(), "must be a valid date")

		// --------------------------------------
		// various date layouts
		// --------------------------------------

		// layout "2006-01-02"
		f = parseDate("2020-07-23")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "2020-07-23")
		assert.Equal(t, f.Date, time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC))

		// layout "January-2-2006"
		f = parseDate("July-23-2020")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "July-23-2020")
		assert.Equal(t, f.Date, time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC))

		// layout "2-January-2006"
		f = parseDate("23-July-2020")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "23-July-2020")
		assert.Equal(t, f.Date, time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC))

		// layout "January 2, 2006"
		f = parseDate("July 23, 2020")
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "July 23, 2020")
		assert.Equal(t, f.Date, time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC))

		// --------------------------------------
		// lowercase date input
		// --------------------------------------

		// layout "January-2-2006"
		f = parseDate("july-23-2020") // lowercase
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "july-23-2020")
		assert.Equal(t, f.Date, time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC))

		// layout "2-January-2006"
		f = parseDate("23-july-2020") // lowercase
		assert.Equal(t, f.IsBlank(), false)
		assert.Equal(t, f.Value(), "23-july-2020")
		assert.Equal(t, f.Date, time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC))
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
	t.Run("StrRequired", func(t *testing.T) {
		// ok
		field1 := parseString("test")
		err1 := StrRequired()(field1)
		assert.Equal(t, err1, nil)

		// error
		field2 := parseString("")
		err2 := StrRequired()(field2)
		assert.Equal(t, err2, errBlankValue)
	})

	t.Run("StrLt", func(t *testing.T) {
		// ok
		field1 := parseString("abcd")
		err1 := StrLt(5)(field1)
		assert.Equal(t, err1, nil)

		// error
		field2 := parseString("abcde")
		err2 := StrLt(5)(field2)
		assert.Equal(t, err2.Error(), "must be less than 5 characters")
	})

	t.Run("StrLte", func(t *testing.T) {
		// ok
		field1 := parseString("abcde")
		err1 := StrLte(5)(field1)
		assert.Equal(t, err1, nil)

		// error
		field2 := parseString("abcdef")
		err2 := StrLte(5)(field2)
		assert.Equal(t, err2.Error(), "must be less than or equal to 5 characters")
	})

	t.Run("StrGt", func(t *testing.T) {
		// ok
		field1 := parseString("abcdef")
		err1 := StrGt(5)(field1)
		assert.Equal(t, err1, nil)

		// error
		field2 := parseString("abcde")
		err2 := StrGt(5)(field2)
		assert.Equal(t, err2.Error(), "must be more than 5 characters")
	})

	t.Run("StrGte", func(t *testing.T) {
		// ok
		field1 := parseString("abcdef")
		err1 := StrGte(5)(field1)
		assert.Equal(t, err1, nil)

		// ok
		field2 := parseString("abcde")
		err2 := StrGte(5)(field2)
		assert.Equal(t, err2, nil)

		// error
		field3 := parseString("abcd")
		err3 := StrGte(5)(field3)
		assert.Equal(t, err3.Error(), "must be more than or equal to 5 characters")
	})
}
