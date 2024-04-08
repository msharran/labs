package tmpl

import (
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func UseNewTemplateRenderer(e *echo.Echo, paths ...string) {
	t := template.New("")
	for i := range paths {
		template.Must(t.ParseGlob(paths[i]))
	}
	e.Renderer = newTemplate(t)
	fmt.Println("Defined templates:", t.DefinedTemplates())
}

func newTemplate(templates *template.Template) echo.Renderer {
	return &Template{
		Templates: templates,
	}
}
