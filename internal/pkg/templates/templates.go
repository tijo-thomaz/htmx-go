package templates

import (
	"html/template"
	"io"
	"strings"
)

// FuncMap returns the standard template functions
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"multiply": func(a, b int) int { return a * b },
		"slice": func(s string, start, end int) string {
			if end > len(s) {
				end = len(s)
			}
			if start > len(s) {
				return ""
			}
			return s[start:end]
		},
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
	}
}

// Render parses and executes the base layout with a page template
func Render(w io.Writer, page string, data interface{}) error {
	tmpl, err := template.New("base.html").Funcs(FuncMap()).ParseFiles(
		"web/templates/layouts/base.html",
		"web/templates/pages/"+page,
	)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, "base", data)
}

