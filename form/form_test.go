package form

import (
	"fmt"
	"math/big"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func assert(t *testing.T, a interface{}, expected interface{}) {
	if a == expected {
		return
	}
	_, filename, line, _ := runtime.Caller(1)
	msg := "%s:%d expected %v (type %v), got %v (type %v)\n"
	fmt.Printf(msg, filepath.Base(filename), line, expected, reflect.TypeOf(expected), a, reflect.TypeOf(a))
	t.FailNow()
}

func testValidator(key, value string) *Validator {
	data := map[string]string{
		key: value,
	}
	return NewValidator(data)
}

func TestCustomTypes(t *testing.T) {
	t.Run("string len with empty string", func(t *testing.T) {
		s := String{Value: ""}
		assert(t, s.Len(), 0)
	})

	t.Run("string len with english string", func(t *testing.T) {
		s := String{Value: "hello"}
		assert(t, s.Len(), 5)
	})

	t.Run("string len with japanese string", func(t *testing.T) {
		s := String{Value: "ラーメン"}
		assert(t, s.Len(), 4)
	})

	t.Run("date string yyyy-mm-dd", func(t *testing.T) {
		d := Date{Value: time.Date(2024, 1, 15, 7, 47, 28, 0, time.UTC)}
		assert(t, d.String(), "2024-01-15")
	})

	t.Run("decimal string precision rounds value up", func(t *testing.T) {
		br, ok := new(big.Rat).SetString("172.3287")
		assert(t, ok, true)

		d := Decimal{Value: *br}
		assert(t, d.String(), "172.33")
	})

	t.Run("decimal string precision rounds value down", func(t *testing.T) {
		br, ok := new(big.Rat).SetString("172.3237")
		assert(t, ok, true)

		d := Decimal{Value: *br}
		assert(t, d.String(), "172.32")
	})

}

func TestValidator(t *testing.T) {
	t.Run("new validator with nil data", func(t *testing.T) {
		v := NewValidator(nil)
		assert(t, len(v.data), 0)
	})

	t.Run("new validator with valid data", func(t *testing.T) {
		v := testValidator("name", "rio")

		assert(t, len(v.data), 1)
		assert(t, v.data["name"], "rio")
	})

	t.Run("new validator with spaces in data", func(t *testing.T) {
		v := testValidator(" name ", " rio ")

		assert(t, len(v.data), 1)
		// check that spaces are trimmed from the keys and values.
		assert(t, v.data["name"], "rio")
	})

	t.Run("new validator with empty value", func(t *testing.T) {
		v := testValidator(" name ", "")

		assert(t, len(v.data), 0)
		// check that the "name" key does not exist.
		name, ok := v.data["name"]
		assert(t, name, "")
		assert(t, ok, false)
	})
}
