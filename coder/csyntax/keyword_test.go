package csyntax

import (
	"testing"
)

func TestWriteKeywords(t *testing.T) {
	builder, writer := makeTestWriter(testStyle1)

	err := writer.Write(0, KeywordIf, KeywordReturn)
	if err != nil {
		t.Fatalf("WriteItems Keywords failed: %s", err)
	}

	expected := "ifreturn"
	result := builder.String()
	if result != expected {
		t.Fatalf("WriteItems Keywords result wrong, expected '%s', got '%s'", expected, result)
	}
}
