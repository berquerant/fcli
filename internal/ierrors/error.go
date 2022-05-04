package ierrors

import (
	"errors"
	"fmt"
)

// Wrapper returns a new error.
type Wrapper func(format string, v ...any) error

type WrapperBuilder struct {
	err error
	msg string
}

func NewWrapperBuilder() *WrapperBuilder {
	return &WrapperBuilder{}
}

// Err sets an error wraps new error.
func (s *WrapperBuilder) Err(err error) *WrapperBuilder {
	s.err = err
	return s
}

// Msg sets a common message added to new error.
func (s *WrapperBuilder) Msg(format string, v ...any) *WrapperBuilder {
	s.msg = fmt.Sprintf(format, v...)
	return s
}

func (s *WrapperBuilder) Build() Wrapper {
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
