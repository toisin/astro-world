package util

import (
	"fmt"
	"os"
)

// MaybeExit prints to stderr and exits if there was an error.
func MaybeExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
