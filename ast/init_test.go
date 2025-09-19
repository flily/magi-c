package ast

import (
	"testing"
)

func TestOperatorListOrderInStringLength(t *testing.T) {
	for i := 1; i < len(OperatorList); i++ {
		if len(OperatorList[i-1]) < len(OperatorList[i]) {
			t.Errorf("operatorList is not in descending order of string length: %q (len=%d) < %q (len=%d)",
				OperatorList[i-1], len(OperatorList[i-1]),
				OperatorList[i], len(OperatorList[i]))
		}
	}
}
