package form

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
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
	t.Run("test String len with empty string", func(t *testing.T) {
		s := String{Value: ""}
		assert(t, s.Len(), 0)
	})

	t.Run("test String len with english string", func(t *testing.T) {
		s := String{Value: "hello"}
		assert(t, s.Len(), 5)
	})

	t.Run("test String len with japanese string", func(t *testing.T) {
		s := String{Value: "ラーメン"}
		assert(t, s.Len(), 4)
	})
}

func TestValidator(t *testing.T) {
	// t.Run("", func(t *testing.T) {})

	t.Run("new validator with nil data", func(t *testing.T) {
		v := NewValidator(nil)
		assert(t, len(v.data), 0)
		assert(t, len(v.errors), 0)
	})

	t.Run("new validator with valid data", func(t *testing.T) {
		v := testValidator("name", "rio")
		assert(t, len(v.data), 1)
		assert(t, len(v.errors), 0)
		assert(t, v.data["name"], "rio")
	})

	t.Run("new validator with spaces in data", func(t *testing.T) {
		v := testValidator(" name ", " rio ")
		assert(t, len(v.data), 1)
		assert(t, len(v.errors), 0)
		assert(t, v.data["name"], "rio")
	})
}
