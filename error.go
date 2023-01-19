package jsonrpc2

import (
	"fmt"
	"net/http"
)

func statusCode(err *Error) int {
	switch err.Code {
	case -32600:
		return http.StatusBadRequest
	case 32601:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// When a rpc call encounters an error, the Response Object MUST contain the error member.
type Error struct {
	// A Number that indicates the error type that occurred.
	// This MUST be an integer.
	Code int `json:"code"`
	// A String providing a short description of the error.
	// The message SHOULD be limited to a concise single sentence.
	Message string `json:"message"`
	// A Primitive or Structured value that contains additional information about the error.
	// This may be omitted.
	// The value of this member is defined by the Server (e.g. detailed error information, nested errors etc.).
	Data any `json:"data,omitempty"`
}

// Internal JSON-RPC error.
func NewInternalError() *Error {
	return &Error{
		Code:    -32603,
		Message: "Internal JSON-RPC error",
	}
}

// Invalid method parameter(s).
func NewInvalidParamsError() *Error {
	return &Error{
		Code:    -32602,
		Message: "Invalid params",
	}
}

// The JSON sent is not a valid Request object.
func NewInvalidRequestError() *Error {
	return &Error{
		Code:    -32600,
		Message: "Invalid Request",
	}
}

// The method does not exist / is not available.
func NewMethodNotFoundError() *Error {
	return &Error{
		Code:    -32601,
		Message: "Method not found",
	}
}

// Invalid JSON was received by the server.
// An error occurred on the server while parsing the JSON text.
func NewParseError() *Error {
	return &Error{
		Code:    -32700,
		Message: "Parse error",
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
