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
		"test.mc:1:23: error: duplicated function argument name: 'a'",
		"    1 | fun add(a int, b int, a int) (int) {",
		"      |                       ^",
		"      |                       duplicated name",
		"test.mc:1:9: note: first declared here",
		"    1 | fun add(a int, b int, a int) (int) {",
		"      |         ^",
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

func TestCheckFunctionReturnValueCountMismatched1(t *testing.T) {
	code := strings.Join([]string{
		"fun addAndSub(a int, b int) (int, int) {",
		"    return a + b",
		"}",
	}, "\n")

	expected := strings.Join([]string{
		"test.mc:2:12: error: function return value count mismatch, expect 2, got 1",
		"    2 |     return a + b",
		"      |            ^ ^ ^",
		"      |            SHALL return 2 values",
		"test.mc:1:30: note: return value types is declared here",
		"    1 | fun addAndSub(a int, b int) (int, int) {",
		"      |                              ^^^^ ^^^",
	}, "\n")

	checkCodeError(t, code, expected)
}

func TestCheckFunctionReturnValueCountMismatched2(t *testing.T) {
	code := strings.Join([]string{
		"fun addAndSub(a int, b int) (int, int) {",
		"    return",
		"}",
	}, "\n")

	expected := strings.Join([]string{
		"test.mc:2:5: error: function return value count mismatch, expect 2, got 0",
		"    2 |     return",
		"      |     ^^^^^^",
		"      |     SHALL return 2 values",
		"test.mc:1:30: note: return value types is declared here",
		"    1 | fun addAndSub(a int, b int) (int, int) {",
		"      |                              ^^^^ ^^^",
	}, "\n")

	checkCodeError(t, code, expected)
}

func TestCheckFunctionMissingReturnStatement1(t *testing.T) {
	code := strings.Join([]string{
		"fun foo(a int, b int) (int) {",
		"}",
	}, "\n")

	expected := strings.Join([]string{
		"test.mc:2:1: error: function missing return statement and reach the end of function",
		"    2 | }",
		"      | ^",
		"test.mc:1:24: note: function return value types is declared here",
		"    1 | fun foo(a int, b int) (int) {",
		"      |                        ^^^",
	}, "\n")

	checkCodeError(t, code, expected)
}

func TestCheckFunctionMissingReturnStatement2(t *testing.T) {
	code := strings.Join([]string{
		"fun foo(a int, b int) (int, int) {",
		"}",
	}, "\n")

	expected := strings.Join([]string{
		"test.mc:2:1: error: function missing return statement and reach the end of function",
		"    2 | }",
		"      | ^",
		"test.mc:1:24: note: function return value types is declared here",
		"    1 | fun foo(a int, b int) (int, int) {",
		"      |                        ^^^^ ^^^",
	}, "\n")

	checkCodeError(t, code, expected)
}

func TestCheckFunctionMainDeclarationWithVoidReturnType(t *testing.T) {
	code := strings.Join([]string{
		"fun main() {",
		"    return 0",
		"}",
	}, "\n")

	checkCodeCorrect(t, code)
}

func TestCheckFunctionMainDeclarationWithSingleIntReturnType(t *testing.T) {
	code := strings.Join([]string{
		"fun main() (int) {",
		"    return 0",
		"}",
	}, "\n")

	checkCodeCorrect(t, code)
}

func TestCheckFunctionMainDeclarationWithMultipleReturnTypes(t *testing.T) {
	code := strings.Join([]string{
		"fun main() (int, int) {",
		"    return 0, 0",
		"}",
	}, "\n")

	expected := strings.Join([]string{
		"test.mc:1:13: error: function 'main' must have return type 'int' or no return type, got 2 return types",
		"    1 | fun main() (int, int) {",
		"      |             ^^^^ ^^^",
		"      |             int or no return type",
	}, "\n")

	checkCodeError(t, code, expected)
}
