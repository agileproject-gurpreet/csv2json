package csv2jsonx

import (
	"encoding/json"
	"io"
	"os"

	"github.com/agileproject-gurpreet/csv2json/internal/parser"
)

func ConvertReader(r io.Reader) ([]byte, error) {
	records, err := parser.ParseCSV(r)
	if err != nil {
		return nil, err
	}
	return json.Marshal(records)
}

func ConvertFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ConvertReader(file)
}
