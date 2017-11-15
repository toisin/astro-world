package util

import (
	"fmt"
	"os"
)

// OpenFileFromArg opens a file passed as the first argument on the command line.
func OpenFileFromArg() *os.File {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Missing required file\n")
		os.Exit(1)
	}

	filename := os.Args[1]

	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open %s\n", filename)
		os.Exit(1)
	}
	return f
}
