package rio

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"
)

// ------------------------------------------------------------------
//
//
// Http Status Responses
//
//
// ------------------------------------------------------------------

func Http200(w http.ResponseWriter) {
	status := http.StatusOK
	http.Error(w, http.StatusText(status), status)
}

func Http301(w http.ResponseWriter) {
	status := http.StatusMovedPermanently
	http.Error(w, http.StatusText(status), status)
}

func Http400(w http.ResponseWriter) {
	status := http.StatusBadRequest
	http.Error(w, http.StatusText(status), status)
}

func Http401(w http.ResponseWriter) {
	status := http.StatusUnauthorized
	http.Error(w, http.StatusText(status), status)
}

func Http403(w http.ResponseWriter) {
	status := http.StatusForbidden
	http.Error(w, http.StatusText(status), status)
}

func Http404(w http.ResponseWriter) {
	status := http.StatusNotFound
	http.Error(w, http.StatusText(status), status)
}

func Http500(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	http.Error(w, http.StatusText(status), status)
}

// ------------------------------------------------------------------
//
//
// Json Responses
//
//
// ------------------------------------------------------------------

func Json200(w http.ResponseWriter, data any) {
	writeJson(w, http.StatusOK, data)
}

func Json201(w http.ResponseWriter, data any) {
	writeJson(w, http.StatusCreated, data)
}

func Json204(w http.ResponseWriter, data any) {
	writeJson(w, http.StatusNoContent, data)
}

func Json301(w http.ResponseWriter, data any) {
	writeJson(w, http.StatusMovedPermanently, data)
}

func Json400(w http.ResponseWriter, data any) {
	writeJson(w, http.StatusBadRequest, data)
}

func Json401(w http.ResponseWriter, data any) {
	writeJson(w, http.StatusUnauthorized, data)
}

func Json403(w http.ResponseWriter, data any) {
	writeJson(w, http.StatusForbidden, data)
}

func Json404(w http.ResponseWriter, data any) {
	writeJson(w, http.StatusNotFound, data)
}

func Json500(w http.ResponseWriter, data any) {
	writeJson(w, http.StatusInternalServerError, data)
}

func writeJson(w http.ResponseWriter, status int, data any) {
	if data == nil {
		data = struct {
			Message string `json:"message"`
		}{
			Message: http.StatusText(status),
		}
	}

	js, err := json.Marshal(data)
	if err != nil {
		Http500(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

// ------------------------------------------------------------------
//
//
// Template Functions and Function Helpers
//
//
// ------------------------------------------------------------------

func SafeHTML(content string) template.HTML {
	return template.HTML(content)
}

func TimeDisplay(d time.Time) string {
	return d.Format("3:04 PM")
}

func DateDisplay(d time.Time) string {
	return d.Format("January 02, 2006")
}

func DateTimeDisplay(d time.Time) string {
	return d.Format("January 02, 2006, 3:04 PM")
}

func WrapString(val string) func() string {
	return func() string {
		return val
	}
}

func WrapBool(val bool) func() bool {
	return func() bool {
		return val
	}
}

func WrapInt(val int) func() int {
	return func() int {
		return val
	}
}

func WrapFloat(val float64) func() float64 {
	return func() float64 {
		return val
	}
}

func WrapTime(val time.Time) func() time.Time {
	return func() time.Time {
		return val
	}
}
