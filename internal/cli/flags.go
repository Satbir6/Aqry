package cli

type flags struct {
	recordType  string
	all         bool
	json        bool
	server      string
	timeoutSec  int
	interactive bool
	noColor     bool
	version     bool
}
