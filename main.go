package main

import (
  "html/template"
  "io"
  "github.com/labstack/echo/middleware"
	"github.com/labstack/echo"
	"os"
	"net/http"
)

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", GetNamespaceDeployments())
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
  renderer := &TemplateRenderer{
      templates: template.Must(template.ParseGlob("*.html")),
  }
  e := echo.New()
  e.Renderer = renderer
	e.Use(middleware.Logger())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
    LogLevel:  0,
		DisableStackAll: true,
	}))


}
