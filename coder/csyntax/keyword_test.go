package csyntax

import (
	"testing"
)

func TestWriteKeywords(t *testing.T) {
	builder, writer := makeTestWriter(testStyle1)

	writer.Write(0, KeywordIf, KeywordReturn)
	expected := "ifreturn"
	result := builder.String()
	if result != expected {
		t.Fatalf("WriteItems Keywords result wrong, expected '%s', got '%s'", expected, result)
	}
}
