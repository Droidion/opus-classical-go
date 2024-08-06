package main

import (
	"html/template"
	"opus-classical-go/internal/helpers"
	"path/filepath"
)

var functions = template.FuncMap{
	"formatYearsRangeString": helpers.FormatYearsRangeString,
	"formatCatalogueName":    helpers.FormatCatalogueName,
	"formatWorkName":         helpers.FormatWorkName,
	"formatWorkLength":       helpers.FormatWorkLength,
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
