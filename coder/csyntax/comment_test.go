package csyntax

import (
	"strings"
	"testing"
)

func TestCommentEmpty(t *testing.T) {
	comment := NewComment()

	checkInterfaceCodeElement(comment)

	expected := ""
	checkOutputOnStyle(t, testStyle1, expected, comment)
}

func TestCommentSingleLine(t *testing.T) {
	comment := NewComment("the quick brown fox jumps over the lazy dog.")

	checkInterfaceCodeElement(comment)

	expected := "/* the quick brown fox jumps over the lazy dog. */\n"
	checkOutputOnStyle(t, testStyle1, expected, comment)
}

func TestCommentMultiLines(t *testing.T) {
	comment := NewComment(
		"lorem ipsum dolor sit amet,",
		"consectetur adipiscing elit.",
		"sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	)

	checkInterfaceCodeElement(comment)

	expected := strings.Join([]string{
		"/*",
		" * lorem ipsum dolor sit amet,",
		" * consectetur adipiscing elit.",
		" * sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		" */",
	}, "\n") + "\n"
	checkOutputOnStyle(t, testStyle1, expected, comment)
}
