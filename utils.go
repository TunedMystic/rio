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

func Http200(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusOK)
}

func Http301(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusMovedPermanently)
}

func Http400(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusBadRequest)
}

func Http401(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusUnauthorized)
}

func Http403(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusForbidden)
}

func Http404(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusNotFound)
}

func Http500(w http.ResponseWriter) {
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

func Json200(w http.ResponseWriter, data any) error {
	return writeJson(w, data, http.StatusOK)
}

func Json201(w http.ResponseWriter, data any) error {
	return writeJson(w, data, http.StatusCreated)
}

func Json204(w http.ResponseWriter, data any) error {
	return writeJson(w, data, http.StatusNoContent)
}

func Json301(w http.ResponseWriter, data any) error {
	return writeJson(w, data, http.StatusMovedPermanently)
}

func Json400(w http.ResponseWriter, data any) error {
	return writeJson(w, data, http.StatusBadRequest)
}

func Json401(w http.ResponseWriter, data any) error {
	return writeJson(w, data, http.StatusUnauthorized)
}

func Json403(w http.ResponseWriter, data any) error {
	return writeJson(w, data, http.StatusForbidden)
}

func Json404(w http.ResponseWriter, data any) error {
	return writeJson(w, data, http.StatusNotFound)
}

func Json500(w http.ResponseWriter) error {
	return writeJson(w, nil, http.StatusInternalServerError)
}

type defaultJsonMessage struct {
	Message string `json:"message"`
}

func writeJson(w http.ResponseWriter, data any, status int) error {
	if data == nil {
		data = defaultJsonMessage{
			Message: http.StatusText(status),
		}
	}

	if dataStr, ok := data.(string); ok {
		data = defaultJsonMessage{
			Message: dataStr,
		}
	}

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// ------------------------------------------------------------------
//
//
// Template Functions and Function Helpers
//
//
// ------------------------------------------------------------------

func DisplaySafeHTML(content string) template.HTML {
	return template.HTML(content)
}

func DisplayTime(d time.Time) string {
	return d.Format("3:04 PM")
}

func DisplayDate(d time.Time) string {
	return d.Format("January 02, 2006")
}

func DisplayDateTime(d time.Time) string {
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
