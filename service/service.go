package service

import (
	"flexera/model"
	"fmt"
	"strings"
)

type UserCounts struct {
	DesktopCount int
	LaptopCount  int
}

func CalculateMinimumCopies(rows []model.Installation, applicationID int) int {
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
