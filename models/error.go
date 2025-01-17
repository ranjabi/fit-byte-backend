package models

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewError(code int, message string) error {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
