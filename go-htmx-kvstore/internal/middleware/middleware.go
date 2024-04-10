package tmpl

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type Template struct {
	Templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var tpl *template.Template
	var definedName string

	if strings.HasPrefix(name, "pages/") {
		tpl = t.Templates[name]
		definedName = "page.html" // use this to render full page
	} else {
		tpl = t.Templates["partials"]
		definedName = name
	}

	if tpl == nil {
		return fmt.Errorf("Template not found: %s", name)
	}

	return tpl.ExecuteTemplate(w, definedName, data)
}

func MustCompileTemplates(e *echo.Echo) {
	pages, err := filepath.Glob("web/tmpl/pages/*.html")
	if err != nil {
		panic(err)
	}

	tpls := make(map[string]*template.Template, len(pages)+1) // +1 to render partials separately
	for _, p := range pages {
		t, err := template.New("").ParseFiles("web/tmpl/page.html", p)
		if err != nil {
			panic(err)
		}

		t, err = t.ParseGlob("web/tmpl/partials/*.html")
		if err != nil {
			panic(err)
		}

		tpls["pages/"+filepath.Base(p)] = t
	}

	t, err := template.New("").ParseGlob("web/tmpl/partials/*.html")
	if err != nil {
		panic(err)
	}

	tpls["partials"] = t

	e.Renderer = &Template{
		Templates: tpls,
	}

	// t := template.New("")
	// for i := range paths {
	// 	template.Must(t.ParseGlob(paths[i]))
	// }
	// e.Renderer = newTemplate(t)
	for k, v := range tpls {
		fmt.Println("Template:", k)
		fmt.Println("Defined templates:", v.DefinedTemplates())
	}
}
