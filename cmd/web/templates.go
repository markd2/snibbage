package main

import (
	"html/template"
	"path/filepath"

	"snibbage.borkware.com/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet models.Snippet 
	Snippets []models.Snippet
}


func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// parse the base template file into a template set
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// call parseGlob on this template set to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// call ParseFiles() on this template set to add the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// e.g. "home.tmpl" pointing to its parsed version
		cache[name] = ts
	}

	return cache, nil
}
