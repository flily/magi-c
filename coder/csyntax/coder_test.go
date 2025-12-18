package csyntax

import (
	"testing"

	"strings"
)

func makeTestWriter(sylte *CodeStyle) (*strings.Builder, *StyleWriter) {
	var b strings.Builder
	writer := sylte.MakeWriter(&b)
	return &b, writer
}

func TestCodeStyleClone(t *testing.T) {
	newStyle := KRStyle.Clone()

	newStyle.Indent = "\t"

	if KRStyle.Indent == newStyle.Indent {
		t.Fatalf("CodeStyle Clone failed: modifying clone affected original")
	}
}
