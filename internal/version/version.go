// version/version.go
package version

import (
	"encoding/json"
	"fmt"
	"runtime"
)

var (
	// These variables are set during build using -ldflags
	version    = "dev"
	commitHash = "unknown"
	buildTime  = "unknown"
	// Protocol version should be bumped when making breaking changes
	protocolVersion = "v1"
)

// Info represents version information
type Info struct {
	Version         string `json:"version"`
	CommitHash      string `json:"commit_hash"`
	BuildTime       string `json:"build_time"`
	ProtocolVersion string `json:"protocol_version"`
	GoVersion       string `json:"go_version"`
	OS              string `json:"os"`
	Arch            string `json:"arch"`
}

// Get returns the current version information
func Get() Info {
	return Info{
		Version:         version,
		CommitHash:      commitHash,
		BuildTime:       buildTime,
		ProtocolVersion: protocolVersion,
		GoVersion:       runtime.Version(),
		OS:              runtime.GOOS,
		Arch:            runtime.GOARCH,
	}
}

// String returns a formatted version string
func (i Info) String() string {
	return fmt.Sprintf("Usercanal GO-SDK %s (Protocol %s)\nCommit: %s\nBuilt: %s\n%s %s/%s",
		i.Version,
		i.ProtocolVersion,
		i.CommitHash,
		i.BuildTime,
		i.GoVersion,
		i.OS,
		i.Arch,
	)
}

// JSON returns version information as JSON
func (i Info) JSON() string {
	b, _ := json.MarshalIndent(i, "", "  ")
	return string(b)
}

// UserAgent returns a user agent string for HTTP headers
func (i Info) UserAgent() string {
	return fmt.Sprintf("usercanal-sdk/%s (%s; %s/%s) %s",
		i.Version,
		i.CommitHash[:7],
		i.OS,
		i.Arch,
		i.GoVersion,
	)
}

// Short returns just the version number
func (i Info) Short() string {
	return i.Version
}

// IsProduction returns true if this is a production build
func (i Info) IsProduction() bool {
	return i.Version != "dev"
}
