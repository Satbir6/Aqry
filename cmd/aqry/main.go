package main

import (
	"os"

	"aqry/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
