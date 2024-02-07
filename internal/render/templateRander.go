package render

import (
	"path/filepath"
	"text/template"
	"time"
)

type TemplatesHTML map[string]*template.Template

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format("2/1/2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateHTML(dir string) (TemplatesHTML, error) {
	tmlp := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "/page", "*page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "/partials", "*partial.html"))
		if err != nil {
			return nil, err
		}

		tmlp[name] = ts

	}

	return tmlp, nil
}
