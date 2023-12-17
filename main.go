package main

import (
	"encoding/csv"
	"flexera/model"
	"flexera/service"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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
	startTime := time.Now()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 5

	// Read and ignore the CSV header
	if _, err := reader.Read(); err != nil {
		fmt.Println("Error reading header:", err)
		os.Exit(1)
	}

	var rows []model.Installation

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading record:", err)
			os.Exit(1)
		}

		computerID, _ := strconv.Atoi(record[0])
		userID, _ := strconv.Atoi(record[1])
		applicationID, _ := strconv.Atoi(record[2])
		computerType := strings.ToUpper(strings.TrimSpace(record[3]))

		rows = append(rows, model.Installation{
			ComputerID:    computerID,
			UserID:        userID,
			ApplicationID: applicationID,
			ComputerType:  computerType,
		})
	}

	applicationID := 374
	minimumCopies := service.CalculateMinimumCopies(rows, applicationID)
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Total number of copies of application id %d needed by the company: %d\n", applicationID, minimumCopies)
	fmt.Printf("Time taken: %v\n", elapsedTime)

}
