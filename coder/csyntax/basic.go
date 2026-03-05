package csyntax

import (
	"fmt"
)

type CStandard int

const (
	C89 CStandard = iota
	C99

	EOLCR   = "\r"
	EOLLF   = "\n"
	EOLCRLF = "\r\n"
)

type StyleBoolean bool

func (b StyleBoolean) Not() StyleBoolean {
	return !b
}

func (b StyleBoolean) Select(value CodeElement) CodeElement {
	if b {
		return value
	}

	return DelimiterNone
}

type StringElement string

func NewIntegerStringElement(i int) StringElement {
	return StringElement(fmt.Sprintf("%d", i))
}

func FormatStringElement(format string, args ...any) StringElement {
	s := fmt.Sprintf(format, args...)
	return StringElement(s)
}

func (e StringElement) codeElement() {}

func (e StringElement) String() string {
	return string(e)
}

func (e StringElement) ItemString() string {
	return e.String()
}

type DelimiterCharacter string

func NewDelimiter(c string) DelimiterCharacter {
	return DelimiterCharacter(c)
}

func (d DelimiterCharacter) codeElement() {}

func (d DelimiterCharacter) String() string {
	return string(d)
}

func (d DelimiterCharacter) ItemString() string {
	return d.String()
}

const (
	DelimiterNone    StringElement      = ""
	DelimiterSpace   DelimiterCharacter = " "
	DefaultDelimiter DelimiterCharacter = " "
)
