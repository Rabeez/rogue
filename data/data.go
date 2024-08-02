package data

import (
	"embed"
	"encoding/csv"
	"log"

	"github.com/dimchansky/utfbom"
)

//go:embed *
var data embed.FS

// NOTE: The above embed will show error if there is an empty subdirectory in 'assets'. Just make a dummy file to get rid of it

var TestLevel = mustLoadLevel("levels/test_level.csv")

func mustLoadLevel(name string) [][]string {
	file, err := data.Open(name)
	if err != nil {
		log.Fatal("Error opening CSV: ", err)
	}
	defer file.Close()

	// https://github.com/golang/go/issues/33887#issuecomment-644862879
	sr, _ := utfbom.Skip(file)
	reader := csv.NewReader(sr)
	reader.Comma = ','
	records, err := reader.ReadAll()

	if err != nil {
		log.Fatal("Error reading CSV: ", err)
	}

	return records
}
