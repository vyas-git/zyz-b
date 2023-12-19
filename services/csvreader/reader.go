package csvreader

import (
	"encoding/csv"
	"flexera/model"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ReadCSV(file io.Reader) ([]model.Installation, error) {

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 5
	var rows []model.Installation

	// Read and ignore the CSV header
	if _, err := reader.Read(); err != nil {
		fmt.Println("Error reading header:", err)
		return rows, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading record:", err)
			return rows, err
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
	return rows, nil
}

func ReadCSVByLine(reader *csv.Reader) (model.Installation, error) {

	var row model.Installation

	record, err := reader.Read()

	if err != nil {
		return row, err
	}

	computerID, _ := strconv.Atoi(record[0])
	userID, _ := strconv.Atoi(record[1])
	applicationID, _ := strconv.Atoi(record[2])
	computerType := strings.ToUpper(strings.TrimSpace(record[3]))

	row = model.Installation{
		ComputerID:    computerID,
		UserID:        userID,
		ApplicationID: applicationID,
		ComputerType:  computerType,
	}

	return row, nil
}
