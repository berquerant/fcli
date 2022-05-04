package fcli

import (
	"errors"
	"fmt"
)

// ErrorWrapper returns a new error.
type ErrorWrapper func(format string, v ...any) error

type ErrorWrapperBuilder struct {
	err error
	msg string
}

func NewErrorWrapperBuilder() *ErrorWrapperBuilder {
	return &ErrorWrapperBuilder{}
}

// Err sets an error wraps new error.
func (s *ErrorWrapperBuilder) Err(err error) *ErrorWrapperBuilder {
	s.err = err
	return s
}

// Msg sets a common message added to new error.
func (s *ErrorWrapperBuilder) Msg(format string, v ...any) *ErrorWrapperBuilder {
	s.msg = fmt.Sprintf(format, v...)
	return s
}

func (s *ErrorWrapperBuilder) Build() ErrorWrapper {
	return func(format string, v ...any) error {
		x := fmt.Sprintf(format, v...)
		if s.msg != "" {
			x = fmt.Sprintf("%s %s", s.msg, x)
		}
		if s.err != nil {
			return fmt.Errorf("%w %s", s.err, x)
		}
		return errors.New(x)
	}
}
