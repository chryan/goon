package goon

import (
	"testing"
	"strings"
)

func TestSerialise(t *testing.T) {
	return
	m := map[string]interface{}{
		"unit1":   compareUnit,
	}
	if bytes, err := Marshal(m, "goon"); err == nil {
		if lhs, rhs := strings.Trim(string(bytes), "\n"), strings.Trim(string(complexTypeTest), "\n"); lhs != rhs {
			t.Fatalf("Failed to serialise. Strings do not match:\n%+v\n%+v", lhs, rhs)
		}
	}
}
