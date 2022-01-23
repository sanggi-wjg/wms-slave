package e

import (
	"fmt"
)

const (
	Success       = 200
	NotFound      = 404
	InternalError = 500

	InvalidRequestHeader = 1000

	ErrorStockNotFound = 10000
)

var msg = map[int]string{
	Success:       "success",
	NotFound:      "not found",
	InternalError: "error occurred",

	InvalidRequestHeader: "invalid request header",

	ErrorStockNotFound: "can't not found stock",
}

type wmsException struct {
	Code    int
	Message string
}

func NewException(code int) *wmsException {
	e := wmsException{Code: code, Message: msg[code]}
	return &e
}

func NewExceptionAddMsg(code int, message string) *wmsException {
	e := wmsException{Code: code, Message: fmt.Sprintf("%s :%s", msg[code], message)}
	return &e
}

func (e *wmsException) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}
