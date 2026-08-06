package shared

import "fmt"

type Code string

const (
	CodeNotFound        Code = "NOT_FOUND"
	CodeAlreadyExists   Code = "ALREADY_EXISTS"
	CodeValidation      Code = "VALIDATION_ERROR"
	CodeUnauthorized    Code = "UNAUTHORIZED"
	CodeForbidden       Code = "FORBIDDEN"
	CodeInternal        Code = "INTERNAL_ERROR"
	CodeBadRequest      Code = "BAD_REQUEST"
	CodeTooManyRequests Code = "TOO_MANY_REQUESTS"
)

type AppError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func NewAppError(code Code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// func NewNotFound(entity string) *AppError {
// 	return &AppError{
// 		Code:    CodeNotFound,
// 		Message: fmt.Sprintf("%s not found", entity),
// 	}
// }

// func NewAlreadyExists(entity string) *AppError {
// 	return &AppError{
// 		Code:    CodeAlreadyExists,
// 		Message: fmt.Sprintf("%s already exists", entity),
// 	}
// }

// func NewValidation(message string) *AppError {
// 	return &AppError{
// 		Code:    CodeValidation,
// 		Message: message,
// 	}
// }

// func NewUnauthorized(message string) *AppError {
// 	return &AppError{
// 		Code:    CodeUnauthorized,
// 		Message: message,
// 	}
// }

// func NewForbidden(message string) *AppError {
// 	return &AppError{
// 		Code:    CodeForbidden,
// 		Message: message,
// 	}
// }

// func NewInternal(err error) *AppError {
// 	return &AppError{
// 		Code:    CodeInternal,
// 		Message: "internal server error",
// 		Err:     err,
// 	}
// }

// func NewBadRequest(message string) *AppError {
// 	return &AppError{
// 		Code:    CodeBadRequest,
// 		Message: message,
// 	}
// }
