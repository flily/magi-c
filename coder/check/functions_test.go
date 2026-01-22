package check

import (
	"testing"

	"strings"
)

func TestCheckFunctionDeclarationNameNotDuplicated(t *testing.T) {
	code := strings.Join([]string{
		"fun add(a int, b int, c int, d int) (int) {",
		"    return a + b + c + d",
		"}",
	}, "\n")

	checkCodeCorrect(t, code)
}

func TestCheckFunctionDeclarationNameDuplicated(t *testing.T) {
	code := strings.Join([]string{
		"fun add(a int, b int, a int) (int) {",
		"    return a + b",
		"}",
	}, "\n")

	expected := strings.Join([]string{
		"   1:   fun add(a int, b int, a int) (int) {",
		"                ^             ^",
		"                duplicate function argument name: 'a'",
	}, "\n")

	checkCodeError(t, code, expected)
}

func TestCheckFunctionReturnValueCountMatched(t *testing.T) {
	code := strings.Join([]string{
		"fun addAndSub(a int, b int) (int, int) {",
		"    return a + b, a - b",
		"}",
	}, "\n")

	checkCodeCorrect(t, code)
}

func TestCheckFunctionReturnValueCountMismatched(t *testing.T) {
	code := strings.Join([]string{
		"fun addAndSub(a int, b int) (int, int) {",
		"    return a + b",
		"}",
	}, "\n")

	expected := strings.Join([]string{
		"   1:   fun addAndSub(a int, b int) (int, int) {",
		"                                     ^^^^ ^^^",
		"   2:       return a + b",
		"            ^^^^^^ ^ ^ ^",
		"            function return value count mismatch, expect 2, got 1",
	}, "\n")

	checkCodeError(t, code, expected)
}
