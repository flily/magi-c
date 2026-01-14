package csyntax

import (
	"testing"
)

func TestWriteKeywords(t *testing.T) {
	expected := "ifreturn"
	checkOutputOnStyle(t, testStyle1, expected, KeywordIf, KeywordReturn)
}
