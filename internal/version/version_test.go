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

func TestCurrentWithBuildInfoFallsBackToEmbeddedMetadata(t *testing.T) {
	originalValue := Value
	originalCommit := Commit
	originalBuildDate := BuildDate
	t.Cleanup(func() {
		Value = originalValue
		Commit = originalCommit
		BuildDate = originalBuildDate
	})

	Value = "dev"
	Commit = "unknown"
	BuildDate = "unknown"

	info := currentWithBuildInfo(buildInfoSnapshot{
		MainVersion: "(devel)",
		Settings: map[string]string{
			"vcs.revision": "1234567890abcdef",
			"vcs.time":     "2026-03-27T12:34:56Z",
		},
	})

	if info.Version != "1234567890ab" {
		t.Fatalf("expected version from vcs revision, got %#v", info)
	}
	if info.Commit != "1234567890ab" {
		t.Fatalf("expected commit from vcs revision, got %#v", info)
	}
	if info.BuildDate != "2026-03-27T12:34:56Z" {
		t.Fatalf("expected build date from vcs time, got %#v", info)
	}
}
