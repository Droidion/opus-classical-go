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
		startYear  pgtype.Int4
		finishYear pgtype.Int4
		expected   string
	}{
		{pgtype.Int4{Int32: 1720, Valid: true}, pgtype.Int4{Int32: 1795, Valid: true}, "1720–95"},
		{pgtype.Int4{Int32: 1720, Valid: true}, pgtype.Int4{Int32: 1805, Valid: true}, "1720–1805"},
		{pgtype.Int4{Int32: 1720, Valid: true}, pgtype.Int4{Valid: false}, "1720–"},
		{pgtype.Int4{Int32: 0, Valid: true}, pgtype.Int4{Valid: false}, ""},
		{pgtype.Int4{Int32: 0, Valid: true}, pgtype.Int4{Int32: 1805, Valid: true}, "1805"},
	}

	for _, test := range tests {
		result := FormatYearsRangeString(test.startYear, test.finishYear)
		if result != test.expected {
			t.Errorf("FormatYearsRangeString(%d, %d) = %s; want %s", test.startYear.Int32, test.finishYear.Int32, result, test.expected)
		}
	}
}

func TestFormatWorkLength(t *testing.T) {
	tests := []struct {
		minutes  pgtype.Int4
		expected string
	}{
		{pgtype.Int4{Int32: 155, Valid: true}, "2h 35m"},
		{pgtype.Int4{Int32: 60, Valid: true}, "1h"},
		{pgtype.Int4{Int32: 30, Valid: true}, "30m"},
		{pgtype.Int4{Int32: 0, Valid: true}, ""},
		{pgtype.Int4{Int32: -10, Valid: true}, ""},
		{pgtype.Int4{Int32: 90, Valid: true}, "1h 30m"},
	}

	for _, test := range tests {
		result := FormatWorkLength(test.minutes)
		if result != test.expected {
			t.Errorf("FormatWorkLength(%d) = %s; want %s", test.minutes.Int32, result, test.expected)
		}
	}
}

func TestFormatCatalogueName(t *testing.T) {
	tests := []struct {
		name     pgtype.Text
		number   pgtype.Int4
		postfix  pgtype.Text
		expected string
	}{
		{pgtype.Text{String: "BWV", Valid: true}, pgtype.Int4{Int32: 12, Valid: true}, pgtype.Text{String: "p", Valid: true}, "BWV 12p"},
		{pgtype.Text{String: "Op.", Valid: true}, pgtype.Int4{Int32: 9, Valid: true}, pgtype.Text{Valid: false}, "Op. 9"},
		{pgtype.Text{Valid: false}, pgtype.Int4{Int32: 9, Valid: true}, pgtype.Text{Valid: false}, ""},
		{pgtype.Text{String: "BWV", Valid: true}, pgtype.Int4{Valid: false}, pgtype.Text{Valid: false}, ""},
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
		no       pgtype.Int4
		nickname pgtype.Text
		expected string
	}{
		{"Symphony", pgtype.Int4{Int32: 9, Valid: true}, pgtype.Text{String: "Great", Valid: true}, "Symphony No. 9 Great"},
		{"Sonata", pgtype.Int4{Int32: 14, Valid: true}, pgtype.Text{Valid: false}, "Sonata No. 14"},
		{"Prelude", pgtype.Int4{Valid: false}, pgtype.Text{String: "Raindrop", Valid: true}, "Prelude Raindrop"},
		{"", pgtype.Int4{Valid: false}, pgtype.Text{Valid: false}, ""},
	}

	for _, test := range tests {
		result := FormatWorkName(test.title, test.no, test.nickname)
		if result != test.expected {
			t.Errorf("FormatWorkName(%s, %v, %v) = %s; want %s", test.title, test.no, test.nickname, result, test.expected)
		}
	}
}
