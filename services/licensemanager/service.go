package licensemanager

import (
	"flexera/model"
	"fmt"
	"strings"
)

type LicenseManager struct {
	LicenseCount     map[int]model.UserCounts
	DuplicateEntries map[string]int
}

func NewLicenseManager() *LicenseManager {
	return &LicenseManager{
		LicenseCount:     make(map[int]model.UserCounts),
		DuplicateEntries: make(map[string]int),
	}
}
func (lm *LicenseManager) CalculateMinimumCopies(rows []model.Installation, applicationID int) int {
	userCounts := make(map[int]model.UserCounts)
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
				userCounts[userID] = model.UserCounts{
					DesktopCount: userCounts[userID].DesktopCount + 1,
					LaptopCount:  userCounts[userID].LaptopCount,
				}
			} else if computerType == "laptop" {
				userCounts[userID] = model.UserCounts{
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

func (lm *LicenseManager) ProcessRow(row model.Installation, applicationID int) {
	if row.ApplicationID == applicationID {
		userID := row.UserID
		computerType := strings.ToLower(row.ComputerType)
		key := fmt.Sprintf("%d-%d-%d-%s", row.ComputerID, row.UserID, row.ApplicationID, strings.ToLower(row.ComputerType))
		if _, ok := lm.DuplicateEntries[key]; ok {
			return // skip duplicate entries
		}
		lm.DuplicateEntries[key]++

		if computerType == "desktop" {

			lm.LicenseCount[userID] = model.UserCounts{
				DesktopCount: lm.LicenseCount[userID].DesktopCount + 1,
				LaptopCount:  lm.LicenseCount[userID].LaptopCount,
			}

		} else if computerType == "laptop" {

			lm.LicenseCount[userID] = model.UserCounts{
				DesktopCount: lm.LicenseCount[userID].DesktopCount,
				LaptopCount:  lm.LicenseCount[userID].LaptopCount + 1,
			}
		}
	}
}

func (lm *LicenseManager) GetMiniCopiesRequired() int {
	minCopiesRequired := 0

	for _, counts := range lm.LicenseCount {
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
func (lm *LicenseManager) Reset() {
	lm.DuplicateEntries = map[string]int{}
	lm.LicenseCount = map[int]model.UserCounts{}
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
