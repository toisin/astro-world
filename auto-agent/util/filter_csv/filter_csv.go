package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	util ".."
)

func main() {
	// Usage: ./filter_csv filter.csv source.csv

	filterFile := util.OpenFileFromArgAt(1)
	defer filterFile.Close()
	filterReader := util.NewCSVReader(filterFile)

	sourceFile := util.OpenFileFromArgAt(2)
	defer sourceFile.Close()

	r := util.NewCSVReader(sourceFile)
	w := csv.NewWriter(os.Stdout)

	// Header
	headers, err := r.Read()
	util.MaybeExit(err)

	fieldName := filterCSVFieldName(filterReader)
	idx := findIndex(fieldName, headers)
	if idx == -1 {
		log.Fatalf("Could not find a field with name %s", fieldName)
		os.Exit(1)
	}

	err = w.Write(headers)
	util.MaybeExit(err)

	set := filterCSVToSet(filterReader)

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		util.MaybeExit(err)

		val := row[idx]
		if _, ok := set[val]; ok {
			err = w.Write(row)
			util.MaybeExit(err)
		}
	}

	w.Flush()
}

func findIndex(fieldName string, headers []string) int {
	for i, header := range headers {
		if header == fieldName {
			return i
		}
	}
	return -1
}

func filterCSVFieldName(r *csv.Reader) string {
	row, err := r.Read()
	util.MaybeExit(err)

	if len(row) != 1 {
		log.Panicln("The filter csv file should only have one column", len(row))
	}
	return row[0]
}

func filterCSVToSet(r *csv.Reader) map[string]struct{} {
	set := map[string]struct{}{}

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		util.MaybeExit(err)

		if len(row) != 1 {
			log.Panicln("The filter csv file should only have one column", len(row))
		}

		set[row[0]] = struct{}{}
	}
	return set
}
