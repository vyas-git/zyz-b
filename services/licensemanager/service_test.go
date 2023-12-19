// services/licensemanager/licensemanager_test.go
package licensemanager_test

import (
	"flexera/model"
	"flexera/services/licensemanager"
	"testing"
)

func TestLicenseManager_ProcessRow(t *testing.T) {
	tests := []struct {
		name           string
		installation   model.Installation
		applicationID  int
		expectedCounts model.UserCounts
	}{
		{
			name: "Process Desktop Installation",
			installation: model.Installation{
				ComputerID:    1,
				UserID:        1,
				ApplicationID: 374,
				ComputerType:  "DESKTOP",
			},
			applicationID:  374,
			expectedCounts: model.UserCounts{DesktopCount: 1, LaptopCount: 0},
		},
		{
			name: "Process Laptop Installation",
			installation: model.Installation{
				ComputerID:    2,
				UserID:        2,
				ApplicationID: 374,
				ComputerType:  "LAPTOP",
			},
			applicationID:  374,
			expectedCounts: model.UserCounts{DesktopCount: 0, LaptopCount: 1},
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lm := licensemanager.NewLicenseManager()
			lm.ProcessRow(test.installation, test.applicationID)

			// Check if the counts are as expected
			counts, ok := lm.LicenseCount[test.installation.UserID]
			if !ok {
				t.Fatal("UserID not found in LicenseCount")
			}

			if counts != test.expectedCounts {
				t.Errorf("Expected counts: %v, Got counts: %v", test.expectedCounts, counts)
			}
		})
	}
}

func TestLicenseManager_GetMiniCopiesRequired(t *testing.T) {
	tests := []struct {
		name          string
		licenseCounts map[int]model.UserCounts
		expected      int
	}{
		{
			name: "Get Mini Copies Required",
			licenseCounts: map[int]model.UserCounts{
				1: {DesktopCount: 2, LaptopCount: 1},
				2: {DesktopCount: 1, LaptopCount: 3},
				3: {DesktopCount: 0, LaptopCount: 2},
			},
			expected: 7,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lm := licensemanager.NewLicenseManager()
			lm.LicenseCount = test.licenseCounts

			result := lm.GetMiniCopiesRequired()
			if result != test.expected {
				t.Errorf("Expected: %d, Got: %d", test.expected, result)
			}
		})
	}
}
