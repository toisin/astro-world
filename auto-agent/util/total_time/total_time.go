package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/toisin/astro-world/auto-agent/util"
)

func main() {
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
	defer f.Close()

	r := util.NewCSVReader(f)
	ttpu, err := totalTimePerUser(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	// Put in a slice and sort it.
	slice := make(totalInfoSlice, len(ttpu))
	i := 0
	for username, duration := range ttpu {
		slice[i] = totalInfo{username, duration}
		i++
	}

	sort.Sort(slice)

	w := csv.NewWriter(os.Stdout)
	// Write header
	w.Write([]string{"Username", "Duration"})

	for _, info := range slice {
		w.Write([]string{info.username, strconv.FormatInt(int64(info.duration.Minutes()), 10)})
	}
	w.Flush()
}

type totalInfo struct {
	username string
	duration time.Duration
}

type totalInfoSlice []totalInfo

func (tis totalInfoSlice) Len() int           { return len(tis) }
func (tis totalInfoSlice) Less(i, j int) bool { return tis[i].username < tis[j].username }
func (tis totalInfoSlice) Swap(i, j int)      { tis[i], tis[j] = tis[j], tis[i] }

const (
	usernameIdx         = 0
	dateIdx             = 9
	inactiveTimeMinutes = 30
)

func totalTimePerUser(r *csv.Reader) (map[string]time.Duration, error) {
	// Skip header.
	_, err := r.Read()
	if err != nil {
		return nil, err
	}

	type acc struct {
		lastTime time.Time
		duration time.Duration
	}

	data := map[string]acc{}

	for {
		var row []string
		row, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		username := row[usernameIdx]
		t, err := time.Parse(time.RFC3339, row[dateIdx])
		if err != err {
			return nil, err
		}

		userData, ok := data[username]
		if !ok {
			userData = acc{
				lastTime: t,
			}
		} else {
			if userData.lastTime.IsZero() {
				userData.lastTime = t
			} else {
				diff := t.Sub(userData.lastTime)
				if diff.Minutes() < inactiveTimeMinutes {
					userData.duration += t.Sub(userData.lastTime)
				}
			}
			userData.lastTime = t
		}
		data[username] = userData
	}

	rv := make(map[string]time.Duration, len(data))
	for username, userData := range data {
		rv[username] = userData.duration
	}

	return rv, nil
}
