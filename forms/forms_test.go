package forms

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

func TestParseError(t *testing.T) {
	t.Run("ParseError string", func(t *testing.T) {
		err := ParseError{Msg: "test error"}
		assert(t, err.Error(), "test error")
		assert(t, err.String(), "test error")
	})
}
