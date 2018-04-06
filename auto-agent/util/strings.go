package util

import (
	"regexp"
)

func ContainsWord(sentence string, word string) int {
	rp := regexp.MustCompile(`\b` + word + `\b`)
	foundStrings := rp.FindAllString(sentence, -1)
	return len(foundStrings)
}
