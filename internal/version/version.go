// Package version provides the version of the scripting-tools CLI.
package version

// Version can be overridden at build time with:
// go build -ldflags "-X github.com/arduino/scripting-tools/internal/version.Version=v1.0.0"
var Version = "dev"

// Commit can be overridden at build time with:
// go build -ldflags "-X github.com/arduino/scripting-tools/internal/version.Commit=<commit>"
var Commit = "HEAD"

// Timestamp can be overridden at build time with:
// go build -ldflags "-X github.com/arduino/scripting-tools/internal/version.Timestamp=<timestamp>"
var Timestamp = "unknown"
