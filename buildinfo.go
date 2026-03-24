// Buildinfo includes some functions to return information from the built executable,
// such as version, Git commit, architecture, etc.
//
// © 2026 Gwyneth Llewelyn <gwyneth.llewelyn@gwynethllewelyn.net>
// Freely distributed and released under a MIT licence.
package buildinfo

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

// versionInfoType holds the relevant information for this build.
// It is meant to be used as a cache.
type versionInfoType struct {
	version    string    // Runtime version.
	commit     string    // Commit revision number.
	dateString string    // Commit revision time (as a RFC3339 string).
	date       time.Time // Same as before, converted to a time.Time, because that's what the cli package uses.
	builtBy    string    // User who built this (see note).
	goOS       string    // Operating system for this build (from runtime).
	goARCH     string    // Architecture, i.e., CPU type (from runtime).
	goVersion  string    // Go version used to compile this build (from runtime).
	init       bool      // Have we already initialised the cache object?
}

// NOTE: I don't know where the "builtBy" information comes from, so, right now, it gets injected
// during build time, e.g. `go build -ldflags "-X main.TheBuilder=gwyneth"` (gwyneth 20231103)

var (
	versionInfo *versionInfoType // Cached values for this build.
	TheBuilder  string           // To be overwritten via the linker command `go build -ldflags "-X main.TheBuilder=gwyneth"`.
	TheVersion  string           // To be overwritten with `-X main.TheVersion=X.Y.Z`, as above.
)

// Initialises a versionInfo variable.
func initVersionInfo() (vI *versionInfoType, err error) {
	vI = new(versionInfoType)
	if vI.init {
		// already initialised, no need to do anything else!
		return vI, nil
	}
	// get the following entries from the runtime:
	vI.goOS = runtime.GOOS
	vI.goARCH = runtime.GOARCH
	vI.goVersion = runtime.Version()

	// attempt to get some build info as well:
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, fmt.Errorf("no valid build information found")
	}
	// use our supplied version instead of the long, useless, default Go version string.
	if TheVersion == "" {
		vI.version = buildInfo.Main.Version
	} else {
		vI.version = TheVersion
	}

	// Now dig through settings and extract what we can...

	var vcs, rev string // Name of the version control system name (very likely Git) and the revision.
	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs":
			vcs = setting.Value
		case "vcs.revision":
			rev = setting.Value
		case "vcs.time":
			vI.dateString = setting.Value
		}
	}
	vI.commit = "unknown"
	if vcs != "" {
		vI.commit = vcs
	}
	if rev != "" {
		vI.commit += " [" + rev + "]"
	}
	// attempt to parse the date, which comes as a string in RFC3339 format, into a date.Time:
	var parseErr error
	if vI.date, parseErr = time.Parse(vI.dateString, time.RFC3339); parseErr != nil {
		// Note: we can safely ignore the parsing error: either the conversion works, or it doesn't, and we
		// cannot do anything about it... (gwyneth 20231103)
		// However, the AI revision bots dislike this, so we'll assign the current date instead.
		vI.date = time.Now()
	}

	// see comment above
	vI.builtBy = TheBuilder
	// Mark initialisation as complete before returning.
	vI.init = true
	return vI, nil
}

// Returns a pretty-printed version of versionInfo, respecting the String() syntax.
func (vI *versionInfoType) String() string {
	// check if we have a valid builder's name; if yes, add it to the return string.
	// (gwyneth 20251007)
	var builtBy string
	if len(vI.builtBy) > 0 {
		builtBy = " by " + vI.builtBy
	}

	return fmt.Sprintf(
		"\r\t%s\n\t(rev %s)\n\t[%s %s %s]\n\tBuilt on %s%s",
		vI.version,
		vI.commit,
		vI.goOS,
		vI.goARCH,
		vI.goVersion,
		vI.dateString, // Date as string in RFC3339 notation.
		builtBy,
	)
}

// Initialises a global, pre-defined versionInfo variable (we might just need one).
// Panics if allocation failed!
func init() {
	var err error
	if versionInfo, err = initVersionInfo(); err != nil {
		panic(err)
	}
}

// Singleton getters.

// Runtime version.
func Version() string {
	return versionInfo.version
}

// Commit revision number.
func Commit() string {
	return versionInfo.commit
}

// Commit revision time (as a RFC3339 string).
func DateString() string {
	return versionInfo.dateString
}

// Same as before, converted to a time.Time, because that's what the cli package uses.
func Date() time.Time {
	return versionInfo.date
}

// User who built this (see note).
func BuiltBy() string {
	return versionInfo.builtBy
}

// Operating system for this build (from runtime).
func GoOS() string {
	return versionInfo.goOS
}

// Architecture, i.e., CPU type (from runtime).
func GoARCH() string {
	return versionInfo.goARCH
}

// Go version used to compile this build (from runtime).
func GoVersion() string {
	return versionInfo.goVersion
}

// Have we already initialised the cache object?
func IsInit() bool {
	return versionInfo.init
}

// Return everything pretty-printed for the singleton.
func String() string {
	return versionInfo.String()
}
