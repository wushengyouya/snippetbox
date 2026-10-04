package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.alexedwards.net/internal/models"
	"snippetbox.alexedwards.net/ui"
)

type templateData struct {
	CurrentYear     int
	Form            any
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

var templateFunctions = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.DateTime)
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{"html/base.tmpl", "html/partials/*.tmpl", page}
		ts, err := template.New(name).Funcs(templateFunctions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
