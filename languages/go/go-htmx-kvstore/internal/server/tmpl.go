package server

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

func newRenderer() (*tmplRenderer, error) {
	pages, err := filepath.Glob("internal/web/tmpl/pages/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	tpls := make(map[string]*template.Template, len(pages)+1) // +1 to render partials separately
	for _, p := range pages {
		t, err := template.New("").ParseFiles("internal/web/tmpl/page.html", p)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", p, err)
		}

		t, err = t.ParseGlob("internal/web/tmpl/partials/*.html")
		if err != nil {
			return nil, fmt.Errorf("failed to parse partials: %w", err)
		}

		tpls["pages/"+filepath.Base(p)] = t
	}

	t, err := template.New("").ParseGlob("internal/web/tmpl/partials/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse partials: %w", err)
	}

	tpls["partials"] = t

	return &tmplRenderer{
		templates: tpls,
	}, nil
}

type tmplRenderer struct {
	templates map[string]*template.Template
}

func (t *tmplRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var tpl *template.Template
	var definedName string

	if strings.HasPrefix(name, "pages/") {
		tpl = t.templates[name]
		definedName = "page.html" // use this to render full page
	} else {
		tpl = t.templates["partials"]
		definedName = name
	}

	if tpl == nil {
		return fmt.Errorf("Template not found: %s", name)
	}

	return tpl.ExecuteTemplate(w, definedName, data)
}
