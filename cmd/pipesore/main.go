package main

import (
	"fmt"
	"os"

	"github.com/dyson/pipesore/internal/pipesore"
)

func main() {
	s, err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	os.Exit(s)
}

func run() (int, error) {
	if len(os.Args) != 2 {
		return 1, fmt.Errorf("use a single string to define pipeline")
	}

	err := pipesore.Execute(os.Args[1], os.Stdin, os.Stdout)
	if err != nil {
		return 1, fmt.Errorf("error executing pipeline: %w", err)
	}

	return 0, nil
}
