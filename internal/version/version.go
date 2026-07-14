package version

import "fmt"

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func String() string {
	if Commit == "unknown" && Date == "unknown" {
		return Version
	}

	return fmt.Sprintf("%s (commit %s, built %s)", Version, Commit, Date)
}
