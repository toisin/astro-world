package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"

	"github.com/toisin/astro-world/auto-agent/util"
)

const usernameIdx = 0

func main() {
	util.CheckStdinMode("random_user")

	f1, err := os.OpenFile(os.Args[1], os.O_RDWR|os.O_CREATE, 0644)
	util.MaybeExit(err)
	defer f1.Close()

	f2, err := os.OpenFile(os.Args[2], os.O_RDWR|os.O_CREATE, 0644)
	util.MaybeExit(err)
	defer f2.Close()

	rnd := rand.New(rand.NewSource(123))
	p := rnd.Perm(15)

	allRes := make([]*regexp.Regexp, len(p))
	for i, n := range p {
		allRes[i] = regexp.MustCompile(fmt.Sprintf(`\.(?:%d+)\.?$`, n+1))
	}

	// First group
	g1 := allRes[:3]

	// Second group
	g2 := allRes[3:6]

	r := util.NewCSVReader(os.Stdin)
	w1 := csv.NewWriter(f1)
	w2 := csv.NewWriter(f2)

	headers, err := r.Read()
	util.MaybeExit(err)

	err = w1.Write(headers)
	util.MaybeExit(err)
	err = w2.Write(headers)
	util.MaybeExit(err)

	rm2Re := regexp.MustCompile("rm2")

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		util.MaybeExit(err)
		if rm2Re.MatchString(row[usernameIdx]) {

			for _, re := range g1 {
				if re.MatchString(row[usernameIdx]) {
					err = w1.Write(row)
					util.MaybeExit(err)
				}
			}
		} else {
			for _, re := range g2 {
				if re.MatchString(row[usernameIdx]) {
					err = w2.Write(row)
					util.MaybeExit(err)
				}
			}
		}
	}

	w1.Flush()
	w2.Flush()
}
