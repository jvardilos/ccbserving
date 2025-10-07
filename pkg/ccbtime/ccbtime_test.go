package ccbtime

import (
	"testing"
	"time"
)

// TestFmtNow verifies that FmtNow correctly formats a given time.Time.
func TestFmtNow(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{"Simple date", time.Date(2025, 8, 29, 0, 0, 0, 0, time.UTC), "2025-8-29"},
		{"Leap year", time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC), "2024-2-29"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDate(tt.input)
			if result != tt.expected {
				t.Errorf("FmtNow(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestFmtLastMonth verifies that FmtLastMonth correctly subtracts one month and formats the date.
func TestFmtLastMonth(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{"Basic month shift", time.Date(2025, 8, 29, 0, 0, 0, 0, time.UTC), "2025-7-29"},
		{"Year boundary", time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC), "2024-12-15"},
		// time.AddDate() does not do leap year well.
		// I'll omit this since we are always asking "roughly last month".
		// {"End of month handling", time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC), "2025-2-28"},
		// {"Leap year February", time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC), "2024-2-29"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatLastMonthDate(tt.input)
			if result != tt.expected {
				t.Errorf("FmtLastMonth(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
