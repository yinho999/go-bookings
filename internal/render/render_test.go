package render

import (
	"github.com/yinho999/go-bookings/internal/models"
	"net/http"
	"testing"
)

func getSession() (*http.Request, error) {
	// Create a http request
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	// Create a session
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	// Save the session
	r = r.WithContext(ctx)
	return r, nil
}

func TestAddDefaultData(t *testing.T) {
	// Create a template data struct
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	// Add the session to the context
	session.Put(r.Context(), "flash", "123")
	// Call the function
	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("Failed: Flash value 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	// create template cache
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc
	// Create a new request
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	// Create ResponseWriter
	ww := &myWriter{}

	// Should be able to render
	err = RenderTemplate(ww, "home.page.tmpl", &models.TemplateData{}, r)
	if err != nil {
		t.Error("Error writing template to browser")
	}
	// Should not be able to render
	err = RenderTemplate(ww, "non-existent.page.tmpl", &models.TemplateData{}, r)
	if err == nil {
		t.Error("Rendered template that does not exist")
	}
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}
func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
