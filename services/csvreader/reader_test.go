// csvreader/csvreader_test.go
package csvreader_test

import (
	"encoding/csv"
	"flexera/model"
	"flexera/services/csvreader"
	"strings"
	"testing"
)

func TestReadCSV(t *testing.T) {
	csvContent := `ComputerID,UserID,ApplicationID,ComputerType,Comment
1,1,374,LAPTOP,Exported from System A
2,1,374,DESKTOP,Exported from System A
3,1,374,DESKTOP,Exported from System A`

	reader := strings.NewReader(csvContent)

	rows, err := csvreader.ReadCSV(reader)
	if err != nil {
		t.Errorf("Error reading CSV: %v", err)
	}

	expectedRows := []model.Installation{
		{ComputerID: 1, UserID: 1, ApplicationID: 374, ComputerType: "LAPTOP"},
		{ComputerID: 2, UserID: 1, ApplicationID: 374, ComputerType: "DESKTOP"},
		{ComputerID: 3, UserID: 1, ApplicationID: 374, ComputerType: "DESKTOP"},
	}

	if len(rows) != len(expectedRows) {
		t.Errorf("Expected %d rows, got %d", len(expectedRows), len(rows))
		return
	}

	for i, row := range rows {
		if row != expectedRows[i] {
			t.Errorf("Mismatch in row at index %d, expected %+v, got %+v", i, expectedRows[i], row)
		}
	}
}

func TestReadCSVByLine(t *testing.T) {
	csvContent := `1,1,374,LAPTOP,Exported from System A
2,1,374,DESKTOP,Exported from System A
3,1,374,DESKTOP,Exported from System A`

	reader := csv.NewReader(strings.NewReader(csvContent))

	expectedRows := []model.Installation{
		{ComputerID: 1, UserID: 1, ApplicationID: 374, ComputerType: "LAPTOP"},
		{ComputerID: 2, UserID: 1, ApplicationID: 374, ComputerType: "DESKTOP"},
		{ComputerID: 3, UserID: 1, ApplicationID: 374, ComputerType: "DESKTOP"},
	}

	for _, expectedRow := range expectedRows {
		row, err := csvreader.ReadCSVByLine(reader)
		if err != nil {
			t.Errorf("Error reading CSV by line: %v", err)
		}

		if row != expectedRow {
			t.Errorf("Mismatch in row, expected %+v, got %+v", expectedRow, row)
		}
	}
}
