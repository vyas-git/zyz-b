/*
Concurrency added while reading csv of 2 billion lines
Concurrency added while processing and calculation of minimum copies

Please note :
But concurrency in reading local CSV involves i/o operations
adding more goroutines may not significantly improve performance
because most of the time is spent waiting for I/O operations to complete, So time spent is less in one go routine

If its network or external call which makes us to wait more time then adding go routines will helps in performance
*/
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Installation struct {
	ComputerID    int
	UserID        int
	ApplicationID int
	ComputerType  string
}

var rows []Installation

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
	rows, err := readCSV(file)
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		os.Exit(1)
	}
	applicationID := 374
	minimumCopies := calculateMinimumCopies(rows, applicationID)
	endTime := time.Now()

	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Total number of copies of application id %d needed by the company: %d\n", applicationID, minimumCopies)
	fmt.Printf("Time taken: %v\n", elapsedTime)
}

func readCSV(file *os.File) ([]Installation, error) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 5

	// Read and ignore the CSV header
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("Error reading header: %v", err)
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	const numWorkers = 4
	lines := make(chan []string, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go readCSVWorker(reader, lines, &wg, &mu)
	}

	go func() {
		defer close(lines)
		for {
			record, err := reader.Read()
			if err != nil {
				break
			}
			lines <- record
		}
	}()

	wg.Wait()

	return rows, nil
}

func readCSVWorker(reader *csv.Reader, lines chan []string, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	for line := range lines {
		computerID, _ := strconv.Atoi(line[0])
		userID, _ := strconv.Atoi(line[1])
		applicationID, _ := strconv.Atoi(line[2])
		computerType := strings.ToUpper(strings.TrimSpace(line[3]))

		mu.Lock()
		rows = append(rows, Installation{
			ComputerID:    computerID,
			UserID:        userID,
			ApplicationID: applicationID,
			ComputerType:  computerType,
		})
		mu.Unlock()
	}
}

type UserCounts struct {
	DesktopCount int
	LaptopCount  int
}

// Add your calculateMinimumCopies and other functions here

func calculateMinimumCopies(rows []Installation, applicationID int) int {
	userCounts := make(map[int]UserCounts)
	duplicateEntries := make(map[string]int)

	for _, row := range rows {
		if row.ApplicationID == applicationID {
			userID := row.UserID
			computerType := strings.ToLower(row.ComputerType)

			key := fmt.Sprintf("%d-%d-%d-%s", row.ComputerID, row.UserID, row.ApplicationID, strings.ToLower(row.ComputerType))
			if _, ok := duplicateEntries[key]; ok {
				continue // skip duplicate entries
			}
			duplicateEntries[key]++

			if computerType == "desktop" {
				userCounts[userID] = UserCounts{
					DesktopCount: userCounts[userID].DesktopCount + 1,
					LaptopCount:  userCounts[userID].LaptopCount,
				}
			} else if computerType == "laptop" {
				userCounts[userID] = UserCounts{
					DesktopCount: userCounts[userID].DesktopCount,
					LaptopCount:  userCounts[userID].LaptopCount + 1,
				}
			}
		}
	}

	minCopiesRequired := 0

	for _, counts := range userCounts {
		deskCnt := counts.DesktopCount
		lapCnt := counts.LaptopCount
		if lapCnt > 0 && lapCnt != deskCnt-lapCnt {
			minCopiesRequired += max(lapCnt, deskCnt-lapCnt)
		} else {
			minCopiesRequired += deskCnt
		}
	}

	return minCopiesRequired
}

func calculateMinimumCopiesConcurrent(rows []Installation, applicationID int) int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	userCounts := make(map[int]UserCounts)
	duplicateEntries := make(map[string]int)

	chunkSize := 100
	numChunks := len(rows) / chunkSize
	if chunkSize > len(rows) {
		numChunks = 1
	}
	for i := 0; i < numChunks; i++ {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				row := rows[j]
				if row.ApplicationID == applicationID {
					userID := row.UserID
					computerType := strings.ToLower(row.ComputerType)

					key := fmt.Sprintf("%d-%d-%d-%s", row.ComputerID, row.UserID, row.ApplicationID, strings.ToLower(row.ComputerType))
					mu.Lock()
					if _, ok := duplicateEntries[key]; ok {
						mu.Unlock()
						continue // skip duplicate entries
					}
					duplicateEntries[key]++
					mu.Unlock()

					if computerType == "desktop" {
						mu.Lock()
						userCounts[userID] = UserCounts{
							DesktopCount: userCounts[userID].DesktopCount + 1,
							LaptopCount:  userCounts[userID].LaptopCount,
						}
						mu.Unlock()
					} else if computerType == "laptop" {
						mu.Lock()
						userCounts[userID] = UserCounts{
							DesktopCount: userCounts[userID].DesktopCount,
							LaptopCount:  userCounts[userID].LaptopCount + 1,
						}
						mu.Unlock()
					}
				}
			}
		}(i*chunkSize, (i+1)*chunkSize)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	minCopiesRequired := 0

	for _, counts := range userCounts {
		deskCnt := counts.DesktopCount
		lapCnt := counts.LaptopCount
		if lapCnt > 0 && lapCnt != deskCnt-lapCnt {
			minCopiesRequired += max(lapCnt, deskCnt-lapCnt)
		} else {
			minCopiesRequired += deskCnt
		}
	}

	return minCopiesRequired
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
