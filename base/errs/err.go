package errs

import (
	"fmt"
	"net/http"
)

var (
	ErrAccessTokenInvalid = AuthError(ErrCodeAccessTokenInvalid, "Access token invalid")
	ErrNoPermission       = AuthError(ErrCodeNoPermission, "No permission")
	ErrParamsInvalid      = ParamsError(ErrCodeParamsInvalid, "Params invalid")
)

func ParamsError(code string, message string, args ...interface{}) ProtocolError {
	return ProtocolError{
		Code:     code,
		Message:  fmt.Sprintf(message, args...),
		HttpCode: http.StatusBadRequest,
	}
}

func AuthError(code string, message string, args ...interface{}) ProtocolError {
	return ProtocolError{
		Code:     code,
		Message:  fmt.Sprintf(message, args...),
		HttpCode: http.StatusUnauthorized,
	}
}
