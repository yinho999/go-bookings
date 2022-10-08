package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/yinho999/go-bookings/internal/config"
	"github.com/yinho999/go-bookings/internal/models"
	"github.com/yinho999/go-bookings/internal/render"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func getRoutes() http.Handler {
	// what am i going to put in the session
	// register the type of data we want to put in the session
	gob.Register(models.Reservation{})
	// create a log print out in console window, with INFO prefix, and log.Ldate | log.Ltime
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	// create a log print out in console window, with ERROR prefix, and log.Ldate | log.Ltime | log.Lshortfile
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// change this to true when in production
	app.InProduction = false

	// Initialize the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour // 24 hours
	// Setting the session cookie
	// Persist is set to true so that the cookie is stored in the browser
	// even the browser is closed
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	// Only send the cookie over HTTPS, we dont need that in development
	session.Cookie.Secure = app.InProduction
	app.Session = session
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = true

	// Create a new handler repository
	repo := NewRepo(&app)
	// Set the handler repository for the handlers
	NewHandlers(repo)
	// Set the new render app config for the template package
	render.NewTemplates(&app)
	// Chi
	mux := chi.NewRouter()

	// Use middleware
	mux.Use(middleware.Recoverer)
	// Ignore post without csrf token
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)
	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	// File server
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := make(map[string]*template.Template)

	// Get all *page.tmpl files from the templates folder
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
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
		ls, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		// add the layout files to the template
		if len(ls) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
