package util

import (
	"fmt"
	"os"
)

// OpenFileFromArg opens a file passed as the first argument on the command line.
func OpenFileFromArg() *os.File {
	return OpenFileFromArgAt(1)
}

// OpenFileFromArgAt opens a file passed as the idx argument on the command line.
func OpenFileFromArgAt(idx int) *os.File {
	if len(os.Args) < idx+1 {
		fmt.Fprintf(os.Stderr, "Missing required file\n")
		os.Exit(1)
	}

	filename := os.Args[idx]

	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open %s\n", filename)
		os.Exit(1)
	}
	return f
}
