package rio

import "net/http"

type AppError struct {
	Message string
	Status  int
	IsJson  bool
}

func (a *AppError) Error() string {
	return a.Message
}

func (a *AppError) WriteTo(w http.ResponseWriter) error {
	if a.IsJson {
		return writeJson(w, a.Message, a.Status)
	}

	http.Error(w, a.Message, a.Status)
	return nil
}

func HttpError(message string, status int) *AppError {
	return &AppError{
		Message: message,
		Status:  status,
		IsJson:  false,
	}
}

func JsonError(message string, status int) *AppError {
	return &AppError{
		Message: message,
		Status:  status,
		IsJson:  true,
	}
}
