// 错误处理
package errors

import (
	"fmt"

	"github.com/perpower/goframe/utils/pcode"
)

// Error is custom error for additional features.
type Error struct {
	error error      // Wrapped error.
	text  string     // Custom Error text when Error is created, might be empty when its code is not nil.
	code  pcode.Code // Error code if necessary.
}

// ICode is the interface for Code feature.
type ICode interface {
	Error() string
	Code() pcode.Code
}

// IUnwrap is the interface for Unwrap feature.
type IUnwrap interface {
	Error() string
	Unwrap() error
}

// Error implements the interface of Error, it returns all the error as string.
func (err *Error) Error() string {
	if err == nil {
		return ""
	}
	errStr := err.text
	if errStr == "" && err.code != nil {
		errStr = err.code.Message()
	}
	if err.error != nil {
		if errStr != "" {
			errStr += ": "
		}
		errStr += err.error.Error()
	}
	return errStr
}

// New creates and returns an error which is formatted from given text.
func New(text string) error {
	return &Error{
		text: text,
		code: pcode.CodeNil,
	}
}

// Newf returns an error that formats as the given format and args.
func Newf(format string, args ...interface{}) error {
	return &Error{
		text: fmt.Sprintf(format, args...),
		code: pcode.CodeNil,
	}
}

// Wrap wraps error with text. It returns nil if given err is nil.
// Note that it does not lose the error code of wrapped error, as it inherits the error code from it.
func Wrap(err error, text string) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		text:  text,
		code:  Code(err),
	}
}

// Wrapf returns an error annotating err with a stack trace at the point Wrapf is called, and the format specifier.
// It returns nil if given `err` is nil.
// Note that it does not lose the error code of wrapped error, as it inherits the error code from it.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		text:  fmt.Sprintf(format, args...),
		code:  Code(err),
	}
}
