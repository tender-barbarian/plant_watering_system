package errors

import (
	"log/slog"
	"net/http"
)

type ErrorsWrapper struct {
	logger *slog.Logger
}

func NewErrorsWrapper(logger *slog.Logger) ErrorsWrapper {
	return ErrorsWrapper{logger: logger}
}

// The serverError helper writes a log entry at Error level (including the request
// method and URI as attributes), then sends a generic 500 Internal Server Error
// response to the user.
func (e *ErrorsWrapper) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	e.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description // to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (e *ErrorsWrapper) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
