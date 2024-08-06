package helpers

import (
	"github.com/jackc/pgx/v5/pgtype"
	"testing"
)

func TestIsValidYear(t *testing.T) {
	tests := []struct {
		year     int
		expected bool
	}{
		{1234, true},
		{0, false},
		{-123, false},
		{123, true},
		{12345, false},
		{9999, true},
	}

	for _, test := range tests {
		result := isValidYear(test.year)
		if result != test.expected {
			t.Errorf("isValidYear(%d) = %v; want %v", test.year, result, test.expected)
		}
	}
}

func TestSliceYear(t *testing.T) {
	tests := []struct {
		year     int
		expected string
	}{
		{1985, "85"},
		{2000, "00"},
		{123, "123"},
		{12, "12"},
	}

	for _, test := range tests {
		result := sliceYear(test.year)
		if result != test.expected {
			t.Errorf("sliceYear(%d) = %s; want %s", test.year, result, test.expected)
		}
	}
}

func TestCenturyEqual(t *testing.T) {
	tests := []struct {
		year1    int
		year2    int
		expected bool
	}{
		{1320, 1399, true},
		{1300, 1400, false},
		{2000, 2099, true},
		{1999, 2000, false},
		{0, 100, false}, // Invalid years
		{100, 0, false}, // Invalid years
	}

	for _, test := range tests {
		result := centuryEqual(test.year1, test.year2)
		if result != test.expected {
			t.Errorf("centuryEqual(%d, %d) = %v; want %v", test.year1, test.year2, result, test.expected)
		}
	}
}

func TestFormatYearsRangeString(t *testing.T) {
	tests := []struct {
		startYear  int
		finishYear pgtype.Int4
		expected   string
	}{
		{1720, pgtype.Int4{Int32: 1795, Valid: true}, "1720–95"},
		{1720, pgtype.Int4{Int32: 1805, Valid: true}, "1720–1805"},
		{1720, pgtype.Int4{Valid: false}, "1720–"},
		{0, pgtype.Int4{Valid: false}, ""},
		{0, pgtype.Int4{Int32: 1805, Valid: true}, "1805"},
	}

	for _, test := range tests {
		result := FormatYearsRangeString(test.startYear, test.finishYear)
		if result != test.expected {
			t.Errorf("FormatYearsRangeString(%d, %v) = %s; want %s", test.startYear, test.finishYear, result, test.expected)
		}
	}
}

func TestFormatWorkLength(t *testing.T) {
	tests := []struct {
		minutes  int
		expected string
	}{
		{155, "2h 35m"},
		{60, "1h"},
		{30, "30m"},
		{0, ""},
		{-10, ""},
		{90, "1h 30m"},
	}

	for _, test := range tests {
		result := FormatWorkLength(test.minutes)
		if result != test.expected {
			t.Errorf("FormatWorkLength(%d) = %s; want %s", test.minutes, result, test.expected)
		}
	}
}

func TestFormatCatalogueName(t *testing.T) {
	tests := []struct {
		name     *string
		number   *int
		postfix  *string
		expected string
	}{
		{strPtr("BWV"), intPtr(12), strPtr("p"), "BWV 12p"},
		{strPtr("Op."), intPtr(9), nil, "Op. 9"},
		{nil, intPtr(9), nil, ""},
		{strPtr("BWV"), nil, nil, ""},
	}

	for _, test := range tests {
		result := FormatCatalogueName(test.name, test.number, test.postfix)
		if result != test.expected {
			t.Errorf("FormatCatalogueName(%v, %v, %v) = %s; want %s", test.name, test.number, test.postfix, result, test.expected)
		}
	}
}

func TestFormatWorkName(t *testing.T) {
	tests := []struct {
		title    string
		no       *int
		nickname *string
		skipHtml bool
		expected string
	}{
		{"Symphony", intPtr(9), strPtr("Great"), false, "Symphony No. 9 Great"},
		{"Sonata", intPtr(14), nil, false, "Sonata No. 14"},
		{"Prelude", nil, strPtr("Raindrop"), true, "Prelude Raindrop"},
		{"", nil, nil, false, ""},
	}

	for _, test := range tests {
		result := FormatWorkName(test.title, test.no, test.nickname, test.skipHtml)
		if result != test.expected {
			t.Errorf("FormatWorkName(%s, %v, %v, %v) = %s; want %s", test.title, test.no, test.nickname, test.skipHtml, result, test.expected)
		}
	}
}

// Helper functions for creating pointers
func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
