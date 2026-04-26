package utils

import "net/http"

type ErrorCode string

const (
	ErrorCodeShiftNotActive       ErrorCode = "SHIFT_NOT_ACTIVE"
	ErrorCodeShiftAlreadyOpened   ErrorCode = "SHIFT_ALREADY_OPENED"
	ErrorCodeInvalidCredentials   ErrorCode = "INVALID_CREDENTIALS"
	ErrorCodeValidation           ErrorCode = "VALIDATION_ERROR"
	ErrorCodeNotFound             ErrorCode = "NOT_FOUND"
	ErrorCodeDuplicateName        ErrorCode = "DUPLICATE_NAME"
	ErrorCodeInsufficientStock    ErrorCode = "INSUFFICIENT_STOCK"
	ErrorCodeIngredientUsedRecipe ErrorCode = "INGREDIENT_USED_IN_RECIPE"
	ErrorCodeInactiveMenuItem     ErrorCode = "INACTIVE_MENU_ITEM"
	ErrorCodeInactiveIngredient   ErrorCode = "INACTIVE_INGREDIENT"
	ErrorCodeConflict             ErrorCode = "CONFLICT"
	ErrorCodeInternal             ErrorCode = "INTERNAL_ERROR"
	ErrorCodeUnsupportedOperation ErrorCode = "UNSUPPORTED_OPERATION"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Status  int       `json:"-"`
	Details any       `json:"details,omitempty"`
	Err     error     `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func ValidationError(message string, details any) *AppError {
	return &AppError{Code: ErrorCodeValidation, Message: message, Status: http.StatusBadRequest, Details: details}
}

func UnauthorizedError(message string) *AppError {
	return &AppError{Code: ErrorCodeInvalidCredentials, Message: message, Status: http.StatusUnauthorized}
}

func NotFoundError(message string) *AppError {
	return &AppError{Code: ErrorCodeNotFound, Message: message, Status: http.StatusNotFound}
}

func ConflictError(code ErrorCode, message string, details any) *AppError {
	return &AppError{Code: code, Message: message, Status: http.StatusConflict, Details: details}
}

func BusinessError(code ErrorCode, message string, details any) *AppError {
	return &AppError{Code: code, Message: message, Status: http.StatusBadRequest, Details: details}
}

func InternalError(message string, err error) *AppError {
	return &AppError{Code: ErrorCodeInternal, Message: message, Status: http.StatusInternalServerError, Err: err}
}
