package main

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"html/template"
	"path/filepath"
	"strconv"
)

var functions = template.FuncMap{
	"formatYearsRangeString": formatYearsRangeString,
	"formatCatalogueName":    formatCatalogueName,
	"formatWorkName":         formatWorkName,
	"formatWorkLength":       formatWorkLength,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.gohtml")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("./ui/html/partials/*.gohtml")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

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
func formatYearsRangeString(startYear int, finishYear pgtype.Int4) string {
	if !isValidYear(startYear) && !finishYear.Valid {
		return ""
	}
	if !finishYear.Valid {
		return fmt.Sprintf("%d–", startYear)
	}
	finishYearInt := int(finishYear.Int32)
	if !isValidYear(startYear) {
		return strconv.Itoa(finishYearInt)
	}
	if centuryEqual(startYear, finishYearInt) {
		return fmt.Sprintf("%d–%s", startYear, sliceYear(finishYearInt))
	}
	return fmt.Sprintf("%d–%d", startYear, finishYearInt)
}

// FormatWorkLength formats minutes into a string with hours and minutes, like "2h 35m"
func formatWorkLength(lengthInMinutes int) string {
	hours := lengthInMinutes / 60
	minutes := lengthInMinutes % 60
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
func formatCatalogueName(catalogueName *string, catalogueNumber *int, cataloguePostfix *string) string {
	if catalogueName == nil || catalogueNumber == nil {
		return ""
	}
	postfix := ""
	if cataloguePostfix != nil {
		postfix = *cataloguePostfix
	}
	return fmt.Sprintf("%s %d%s", *catalogueName, *catalogueNumber, postfix)
}

// FormatWorkName formats music work full name, like "Symphony No. 9 Great".
func formatWorkName(workTitle string, workNo *int, workNickname *string, skipHtml bool) string {
	if workTitle == "" {
		return ""
	}
	workName := workTitle
	if workNo != nil {
		workName = fmt.Sprintf("%s No. %d", workName, *workNo)
	}
	if workNickname != nil {
		if skipHtml {
			workName = fmt.Sprintf("%s %s", workName, *workNickname)
		} else {
			workName = fmt.Sprintf("%s %s", workName, *workNickname)
		}
	}
	return workName
}
