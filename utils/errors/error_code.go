package errors

import (
	"fmt"
	"strings"

	"github.com/perpower/goframe/utils/pcode"
)

const commaSeparatorSpace = ", "

// NewCode creates and returns an error that has error code and given text.
func NewCode(code pcode.Code, text ...string) error {
	return &Error{
		text: strings.Join(text, commaSeparatorSpace),
		code: code,
	}
}

// NewCodef returns an error that has error code and formats as the given format and args.
func NewCodef(code pcode.Code, format string, args ...interface{}) error {
	return &Error{
		text: fmt.Sprintf(format, args...),
		code: code,
	}
}

// WrapCode wraps error with code and text.
// It returns nil if given err is nil.
func WrapCode(code pcode.Code, err error, text ...string) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		text:  strings.Join(text, commaSeparatorSpace),
		code:  code,
	}
}

// WrapCodef wraps error with code and format specifier.
// It returns nil if given `err` is nil.
func WrapCodef(code pcode.Code, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		text:  fmt.Sprintf(format, args...),
		code:  code,
	}
}

// Code returns the error code of current error.
// It returns `CodeNil` if it has no error code neither it does not implement interface Code.
func Code(err error) pcode.Code {
	if err == nil {
		return pcode.CodeNil
	}
	if e, ok := err.(ICode); ok {
		return e.Code()
	}
	if e, ok := err.(IUnwrap); ok {
		return Code(e.Unwrap())
	}
	return pcode.CodeNil
}

// HasCode checks and reports whether `err` has `code` in its chaining errors.
func HasCode(err error, code pcode.Code) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(ICode); ok {
		return code == e.Code()
	}
	if e, ok := err.(IUnwrap); ok {
		return HasCode(e.Unwrap(), code)
	}
	return false
}
