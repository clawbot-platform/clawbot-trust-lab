package version

import (
	"runtime/debug"
	"strings"
)

var (
	Value     = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
}

type buildInfoSnapshot struct {
	MainVersion string
	Settings    map[string]string
}

func Current() Info {
	return currentWithBuildInfo(readBuildInfo())
}

func currentWithBuildInfo(snapshot buildInfoSnapshot) Info {
	info := Info{
		Version:   strings.TrimSpace(Value),
		Commit:    strings.TrimSpace(Commit),
		BuildDate: strings.TrimSpace(BuildDate),
	}

	if isFallbackVersion(info.Version) {
		switch {
		case strings.TrimSpace(snapshot.MainVersion) != "" && snapshot.MainVersion != "(devel)":
			info.Version = strings.TrimSpace(snapshot.MainVersion)
		case vcsRevision(snapshot) != "":
			info.Version = shortRevision(vcsRevision(snapshot))
		}
	}

	if isFallbackCommit(info.Commit) && vcsRevision(snapshot) != "" {
		info.Commit = shortRevision(vcsRevision(snapshot))
	}

	if isFallbackBuildDate(info.BuildDate) {
		if vcsTime := strings.TrimSpace(snapshot.Settings["vcs.time"]); vcsTime != "" {
			info.BuildDate = vcsTime
		}
	}

	return info
}

func readBuildInfo() buildInfoSnapshot {
	info, ok := debug.ReadBuildInfo()
	if !ok || info == nil {
		return buildInfoSnapshot{}
	}

	settings := make(map[string]string, len(info.Settings))
	for _, item := range info.Settings {
		settings[item.Key] = item.Value
	}

	return buildInfoSnapshot{
		MainVersion: info.Main.Version,
		Settings:    settings,
	}
}

func vcsRevision(snapshot buildInfoSnapshot) string {
	return strings.TrimSpace(snapshot.Settings["vcs.revision"])
}

func shortRevision(value string) string {
	value = strings.TrimSpace(value)
	if len(value) > 12 {
		return value[:12]
	}
	return value
}

func isFallbackVersion(value string) bool {
	value = strings.TrimSpace(value)
	return value == "" || value == "dev"
}

func isFallbackCommit(value string) bool {
	value = strings.TrimSpace(value)
	return value == "" || value == "unknown"
}

func isFallbackBuildDate(value string) bool {
	value = strings.TrimSpace(value)
	return value == "" || value == "unknown"
}
