package scenario

import "testing"

func TestKnownTypes(t *testing.T) {
	types := KnownTypes()
	if len(types) == 0 {
		t.Fatal("expected known scenario types")
	}
}
