package mongoerr

import "encoding/json"

var (
	ErrNotFound = NewMongoError(nil, "Not found", "Can't find document in Database")
)

type MongoError struct {
	Err              error  `json:"-"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developer_message"`
}

func NewMongoError(err error, message, devMessage string) *MongoError {
	return &MongoError{
		Err:              err,
		Message:          message,
		DeveloperMessage: devMessage,
	}
}

func (e *MongoError) Error() string {
	return e.Message
}

func (e *MongoError) Unwrap() error {
	return e.Err
}

func (e *MongoError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}
