// Basic testing, just to see if all was correctly initialised.
package buildinfo

import (
	"testing"
)

// Main testing function, just to confirm we have all data in the singleton.
func TestAllOK(t *testing.T) {
	if !IsInit() {
		t.Errorf("singleton not initialised")
	}
	t.Log("\nBuilt-in variables\n==================",
		"\nTheBuilder:", TheBuilder,
		"\nTheVersion:", TheVersion,
		"\n\nBuilt-in methods\n==================",
		"\nVersion:", Version(),
		"\nCommit:", Commit(),
		"\nDateString:", DateString(),
		"\nDate:", Date(),
		"\nBuiltBy:", BuiltBy(),
		"\nGoOS:", GoOS(),
		"\nGoARCH:", GoARCH(),
		"\nGoVersion:", GoVersion(),
	)
}
