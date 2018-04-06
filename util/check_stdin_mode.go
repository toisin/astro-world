package util

import (
	"fmt"
	"os"
)

func CheckStdinMode(cmd string) {
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage:")
		fmt.Printf("  cat yourfile.txt | %s\n", cmd)
		os.Exit(1)
	}
}
