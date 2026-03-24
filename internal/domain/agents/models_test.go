package agents

import "testing"

func TestKnownRoles(t *testing.T) {
	roles := KnownRoles()
	if len(roles) < 3 {
		t.Fatalf("unexpected known role count: %d", len(roles))
	}
}
