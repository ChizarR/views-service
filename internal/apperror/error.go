package apperror

import "encoding/json"

var (
	UndefinedError = NewAppError(nil, "UndefinedError", "")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developer_message"`
}

func NewAppError(err error, message, devMessage string) *AppError {
	return &AppError{
		Err:              err,
		Message:          message,
		DeveloperMessage: devMessage,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}

func systemError(err error) *AppError {
	return NewAppError(err, "system error", err.Error())
}
