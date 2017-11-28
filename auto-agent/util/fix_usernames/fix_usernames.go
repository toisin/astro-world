package main

import (
	"encoding/csv"
	"io"
	"os"

	util ".."
)

const (
	usernameIdx = 0
)

var mappings = map[string]string{
	"village.rm.7":     "village.rm10.7",
	"village.rm10-.4":  "village.rm10.4",
	"village.rm10.8.":  "village.rm10.8",
	"village.rm.10.14": "village.rm10.14",
}

func main() {
	util.CheckStdinMode("order_rows")

	r := util.NewCSVReader(os.Stdin)
	w := csv.NewWriter(os.Stdout)

	headers, err := r.Read()
	util.MaybeExit(err)

	err = w.Write(headers)
	util.MaybeExit(err)

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		util.MaybeExit(err)

		username := row[usernameIdx]
		if un, ok := mappings[username]; ok {
			row[usernameIdx] = un
		}

		w.Write(row)
	}

	w.Flush()
}
