package main

import (
	"flexera/model"
	"flexera/services/licensemanager"
	"testing"
)

func TestCalculateMinimumCopies(t *testing.T) {

	tests := []struct {
		name     string
		rows     []model.Installation
		appID    int
		expected int
	}{
		{
			name: "Example One ",
			rows: []model.Installation{
				{1, 1, 374, "LAPTOP"},
				{2, 1, 374, "DESKTOP"},
			},
			appID:    374,
			expected: 1,
		}, {
			name: "Example Two ",
			rows: []model.Installation{
				{1, 1, 374, "LAPTOP"},
				{2, 1, 374, "DESKTOP"},
				{3, 2, 374, "DESKTOP"},
				{4, 2, 374, "DESKTOP"},
			},
			appID:    374,
			expected: 3,
		}, {
			name: "Example Three ",
			rows: []model.Installation{
				{1, 1, 374, "LAPTOP"},
				{2, 2, 374, "DESKTOP"},
				{2, 2, 374, "desktop"},
			},
			appID:    374,
			expected: 2,
		}, {
			name: "Basic scenario",
			rows: []model.Installation{
				{1, 1, 374, "LAPTOP"},
				{2, 1, 374, "DESKTOP"},
				{3, 1, 374, "DESKTOP"},
				{4, 2, 374, "DESKTOP"},
				{5, 2, 374, "DESKTOP"},
				{6, 2, 374, "DESKTOP"},
			},
			appID:    374,
			expected: 5,
		},
		{
			name: "No desktops, only laptops",
			rows: []model.Installation{
				{1, 1, 374, "LAPTOP"},
				{2, 1, 374, "LAPTOP"},
				{3, 2, 374, "LAPTOP"},
			},
			appID:    374,
			expected: 3,
		},
		{
			name: "Equal number of laptops and desktops",
			rows: []model.Installation{
				{1, 1, 374, "LAPTOP"},
				{2, 1, 374, "DESKTOP"},
				{3, 2, 374, "DESKTOP"},
			},
			appID:    374,
			expected: 2,
		},
		{
			name:     "No installations",
			rows:     []model.Installation{},
			appID:    374,
			expected: 0,
		},
		{
			name: "Only desktops, no laptops",
			rows: []model.Installation{
				{1, 1, 374, "DESKTOP"},
				{2, 2, 374, "DESKTOP"},
				{3, 3, 374, "DESKTOP"},
			},
			appID:    374,
			expected: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lm := licensemanager.NewLicenseManager()
			for _, row := range test.rows {
				lm.ProccessRow(row, test.appID)
			}
			result := lm.GetMiniCopiesRequired()
			if result != test.expected {
				t.Errorf("Expected: %d, Got: %d", test.expected, result)
			}
		})
	}
}
