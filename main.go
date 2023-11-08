package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net/http"
)

var log = logrus.New()

func init() {
	log.Formatter = new(logrus.JSONFormatter)
	log.Level = logrus.DebugLevel
}

func Index(c echo.Context) error {
	if c.QueryParam("format") == "json" {
		j := GetNamespaceWorkloads()
		return c.JSONPretty(http.StatusOK, j, "  ")
	}
	return c.Render(http.StatusOK, "index.html", GetNamespaceWorkloads())
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
		LogLevel:        0,
		DisableStackAll: true,
	}))
	e.GET("/", Index)
	e.Logger.Fatal(e.Start(":8000"))
}
