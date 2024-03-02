package main

import (
	"fmt"
	"os"

	"github.com/dyson/pipesore/internal/pipesore"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	s, err := pipesore.Run(version, commit, date)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	os.Exit(s)
}
