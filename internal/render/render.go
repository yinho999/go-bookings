package render

import (
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/yinho999/go-bookings/internal/config"
	"github.com/yinho999/go-bookings/internal/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, request *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(request)
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData, request *http.Request) {
	var tc map[string]*template.Template
	// If we are in dev mode not prod mode, dont use template cache,
	// instead create template cache on each request

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// Get the template by its name from the cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	// Create buffer
	buf := new(bytes.Buffer)

	// Add the default data
	td = AddDefaultData(td, request)

	// test the template in buffer
	err := t.Execute(buf, td)
	if err != nil {
		fmt.Println("Error executing template", err)
	}

	// Render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := make(map[string]*template.Template)

	// Get all *page.tmpl files from the templates folder
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// Loop through the pages, create a template and add it to the cache
	for _, page := range pages {
		name := filepath.Base(page)

		// parse the template file and stored the template called name
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// find all the layout files
		ls, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		// add the layout files to the template
		if len(ls) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
