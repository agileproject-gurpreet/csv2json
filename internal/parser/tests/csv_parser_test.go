package tests

import (
	"strings"
	"testing"

	"github.com/agileproject-gurpreet/csv2json/internal/parser"
)

func TestParseCSV(t *testing.T) {
	csvData := "name,age\nAlice,30\nBob,25"
	reader := strings.NewReader(csvData)

	result, err := parser.ParseCSV(reader)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 records, got %d", len(result))
	}
}
