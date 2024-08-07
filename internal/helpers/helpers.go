package helpers

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
)

// isValidYear checks if given number is a 4 digits number, like 1234 (not -123, 123, or 12345).
func isValidYear(num int) bool {
	return num > 1 && num < 10000
}

// sliceYear returns slice of the full year, like 85 from 1985.
func sliceYear(year int) string {
	yearStr := strconv.Itoa(year)
	if len(yearStr) < 4 {
		return yearStr
	}
	return yearStr[2:4]
}

// centuryEqual checks if two given years are of the same century, like 1320 and 1399.
func centuryEqual(year1, year2 int) bool {
	if !isValidYear(year1) || !isValidYear(year2) {
		return false
	}
	getCentury := func(year int) string {
		yearStr := strconv.Itoa(year)
		return yearStr[:2]
	}
	return getCentury(year1) == getCentury(year2)
}

// FormatYearsRangeString formats the range of two years into the string, e.g. "1720–95", or "1720–1805", or "1720–".
// Start year and dash are always present.
// It's supposed to be used for lifespans, meaning we always have birth, but may not have death.
func FormatYearsRangeString(startYear pgtype.Int4, finishYear pgtype.Int4) string {
	if !startYear.Valid && !finishYear.Valid {
		return ""
	}
	if !isValidYear(int(startYear.Int32)) && !isValidYear(int(finishYear.Int32)) {
		return ""
	}
	if startYear.Valid && isValidYear(int(startYear.Int32)) && !finishYear.Valid {
		return fmt.Sprintf("%d–", startYear.Int32)
	}
	finishYearInt := int(finishYear.Int32)
	if !startYear.Valid || !isValidYear(int(startYear.Int32)) {
		return strconv.Itoa(finishYearInt)
	}
	if centuryEqual(int(startYear.Int32), finishYearInt) {
		return fmt.Sprintf("%d–%s", startYear.Int32, sliceYear(finishYearInt))
	}
	return fmt.Sprintf("%d–%d", startYear.Int32, finishYearInt)
}

// FormatWorkLength formats minutes into a string with hours and minutes, like "2h 35m"
func FormatWorkLength(lengthInMinutes pgtype.Int4) string {
	if !lengthInMinutes.Valid {
		return ""
	}
	hours := lengthInMinutes.Int32 / 60
	minutes := lengthInMinutes.Int32 % 60
	if hours == 0 && minutes == 0 {
		return ""
	}
	if hours < 0 || minutes < 0 {
		return ""
	}
	if hours == 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	if minutes == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

// FormatCatalogueName formats catalogue name of the musical work, like "BWV 12p".
func FormatCatalogueName(catalogueName pgtype.Text, catalogueNumber pgtype.Int4, cataloguePostfix pgtype.Text) string {
	if !catalogueName.Valid || !catalogueNumber.Valid {
		return ""
	}
	postfix := ""
	if cataloguePostfix.Valid {
		postfix = cataloguePostfix.String
	}
	return fmt.Sprintf("%s %d%s", catalogueName.String, catalogueNumber.Int32, postfix)
}

// FormatWorkName formats music work full name, like "Symphony No. 9 Great".
func FormatWorkName(workTitle string, workNo pgtype.Int4, workNickname pgtype.Text) string {
	if workTitle == "" {
		return ""
	}
	workName := workTitle
	if workNo.Valid {
		workName = fmt.Sprintf("%s No. %d", workName, workNo.Int32)
	}
	if workNickname.Valid {
		workName = fmt.Sprintf("%s %s", workName, workNickname.String)
	}
	return workName
}
