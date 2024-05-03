package rio

import "net/http"

// AppError is custom error type for Http/Json errors.
type AppError struct {
	Message string
	Status  int
	IsJson  bool
}

// Error satisfies the error interface.
func (a AppError) Error() string {
	return a.Message
}

// WriteTo writes the AppError to the given ResponseWriter.
//
// If the AppError is Json, then a json object containing the message will be written.
// If the AppError is not Json, then a plain text message will be written.
func (a AppError) WriteTo(w http.ResponseWriter) error {
	if a.IsJson {
		return writeJson(w, a.Message, a.Status)
	}

	http.Error(w, a.Message, a.Status)
	return nil
}

// HttpError constructs and returns an Http AppError.
func HttpError(message string, status int) AppError {
	return AppError{
		Message: message,
		Status:  status,
		IsJson:  false,
	}
}

// JsonError constructs and returns a Json AppError.
func JsonError(message string, status int) AppError {
	return AppError{
		Message: message,
		Status:  status,
		IsJson:  true,
	}
}
