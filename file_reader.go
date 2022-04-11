package main

import (
	"encoding/csv"
	"log"
	"os"
)

type fileContent [][]string

func readFile(filename string) fileContent {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error reading input file: %s", err)
	}
	defer f.Close()

	// read CSV values (using csv.Reader)
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading input file: %s", err)
	}

	return data
}
