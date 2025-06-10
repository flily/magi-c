package errors

import (
	"fmt"
)

type Error interface {
	error
	Message() string
	Lower() Error
	Inner() error
	Derive(format string, args ...any) Error
	With(err error) Error
}

type container struct {
	message string
	base    *container
	inner   error
}

func (c *container) Error() string {
	if c == nil {
		return ""
	}

	inner := ""
	if c.inner != nil {
		inner = fmt.Sprintf(" [inner: %s]", c.inner.Error())
	}

	base := ""
	if c.base != nil {
		base = fmt.Sprintf(" < %s", c.base.Error())
	}

	return c.message + base + inner
}

func (c *container) Message() string {
	if c == nil {
		return ""
	}

	return c.message
}

func (c *container) Lower() Error {
	if c == nil {
		return nil
	}

	return c.base
}

func (c *container) Inner() error {
	if c == nil {
		return nil
	}

	return c.inner
}

func (c *container) With(err error) Error {
	c.inner = err
	return c
}

func (c *container) Derive(format string, args ...any) Error {
	return Derive(c, format, args...)
}

func Derive(base *container, format string, args ...any) Error {
	c := &container{
		message: fmt.Sprintf(format, args...),
		base:    base,
	}

	return c
}

func New(format string, args ...any) Error {
	return Derive(nil, format, args...)
}

var (
	Root = New("ErrorRoot")
)
