package main

import (
	"encoding/csv"
	"flexera/services/csvreader"
	"flexera/services/licensemanager"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run program.go <csv_file>")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 5
	startTime := time.Now()

	//minimumCopies := ProcessAllLines(file)

	minimumCopies := ProcessByLine(file)
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Total number of copies of application id %d needed by the company: %d\n", ApplicationID, minimumCopies)
	fmt.Printf("Time taken: %v\n", elapsedTime)

}

func ProcessByLine(file *os.File) int {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 5
	var lm = licensemanager.NewLicenseManager()
	for {
		row, err := csvreader.ReadCSVByLine(reader)
		lm.ProcessRow(row, ApplicationID)
		if err == io.EOF {
			break
		}
	}
	return lm.GetMiniCopiesRequired()

}
func ProcessAllLines(file *os.File) int {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 5
	rows, err := csvreader.ReadCSV(file)
	if err != nil {
		fmt.Println("Error reading csv  file:", err)
		os.Exit(1)
	}
	var lm = licensemanager.NewLicenseManager()
	minimumCopies := lm.CalculateMinimumCopies(rows, ApplicationID)
	return minimumCopies
}
