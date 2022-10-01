package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/yinho999/go-bookings/internal/config"
	"github.com/yinho999/go-bookings/internal/handlers"
	"github.com/yinho999/go-bookings/internal/models"
	"github.com/yinho999/go-bookings/internal/render"
	"log"
	"net/http"
	"time"
)

// Global variables
const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application entry point
func main() {
	// what am i going to put in the session
	// register the type of data we want to put in the session
	gob.Register(models.Reservation{})

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
	session.Cookie.Secure = app.InProduction // Only send the cookie over HTTPS, we dont need that in development

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false

	// Create a new handler repository
	repo := handlers.NewRepo(&app)
	// Set the handler repository for the handlers
	handlers.NewHandlers(repo)

	// Set the new render app config for the template package
	render.NewTemplates(&app)

	fmt.Printf(fmt.Sprintf("Starting application on port %s", portNumber))

	// create a new serve mux and register the routes
	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	// start the server
	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
