package utils

type HttpError struct {
	Code    int    `json:"code"`
	Key     string `json:"error"`
	Message string `json:"message"`
}

func NewHTTPError(code int, key string, msg string) *HttpError {
	return &HttpError{
		Code:    code,
		Key:     key,
		Message: msg,
	}
}

// Error makes it compatible with `error` interface.
func (e *HttpError) Error() string {
	return e.Key + ": " + e.Message
}
