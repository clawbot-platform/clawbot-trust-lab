package version

import "testing"

func TestCurrentReturnsBuildInfo(t *testing.T) {
	originalValue := Value
	originalCommit := Commit
	originalBuildDate := BuildDate
	t.Cleanup(func() {
		Value = originalValue
		Commit = originalCommit
		BuildDate = originalBuildDate
	})

	Value = "1.2.3"
	Commit = "abc123"
	BuildDate = "2026-03-25"

	info := Current()
	if info.Version != "1.2.3" || info.Commit != "abc123" || info.BuildDate != "2026-03-25" {
		t.Fatalf("unexpected version info %#v", info)
	}
}
