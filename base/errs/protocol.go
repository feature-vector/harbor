package errs

import (
	"fmt"
)

type ProtocolError struct {
	Code     string
	Message  string
	HttpCode int

	rawError error
}

func (e ProtocolError) Error() string {
	if e.rawError != nil {
		return fmt.Sprintf("Code: %s\nMessage: %s\nCause: %s", e.Code, e.Message, e.rawError.Error())
	}
	return fmt.Sprintf("Code: %s\nMessage: %s", e.Code, e.Message)
}

func (e ProtocolError) Wrap(err error) ProtocolError {
	e.rawError = err
	return e
}
