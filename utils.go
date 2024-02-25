package rio

import (
	"encoding/json"
	"net/http"
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
	LogError(err.Error())
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
