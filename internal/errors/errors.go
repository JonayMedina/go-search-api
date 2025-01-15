package errors

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrInvalidInput = &AppError{Code: "INVALID_INPUT", Message: "Entrada inválida"}
	ErrNotFound     = &AppError{Code: "NOT_FOUND", Message: "Recurso no encontrado"}
	ErrUnauthorized = &AppError{Code: "UNAUTHORIZED", Message: "No autorizado"}
)
