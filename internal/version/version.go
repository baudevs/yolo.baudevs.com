package version

import (
	"encoding/json"
	"fmt"
	"runtime"
)

var (
	// Version is the current version of YOLO CLI
	Version = "1.0.0-beta.1"

	// GitCommit is the git commit hash, injected at build time
	GitCommit = "unknown"

	// BuildDate is the build date, injected at build time
	BuildDate = "unknown"
)

// Info represents version information
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
	Platform  string `json:"platform"`
}

// Get returns the version info
func Get() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns the string representation of version info
func (i Info) String() string {
	return fmt.Sprintf("YOLO CLI Version: %s\nGit Commit: %s\nBuild Date: %s\nGo Version: %s\nPlatform: %s",
		i.Version,
		i.GitCommit,
		i.BuildDate,
		i.GoVersion,
		i.Platform,
	)
}

// JSON returns the JSON representation of version info
func (i Info) JSON() string {
	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return ""
	}
	return string(data)
}
