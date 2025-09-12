package token

import (
	"testing"
)

func TestOperatorListOrderInStringLength(t *testing.T) {
	for i := 1; i < len(operatorList); i++ {
		if len(operatorList[i-1]) < len(operatorList[i]) {
			t.Errorf("operatorList is not in descending order of string length: %q (len=%d) < %q (len=%d)",
				operatorList[i-1], len(operatorList[i-1]),
				operatorList[i], len(operatorList[i]))
		}
	}
}
