package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/yinho999/go-bookings/internal/config"
	"github.com/yinho999/go-bookings/internal/handlers"
	"github.com/yinho999/go-bookings/internal/helpers"
	"github.com/yinho999/go-bookings/internal/models"
	"github.com/yinho999/go-bookings/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

// Global variables
const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application entry point
func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
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

func run() error {
	// what am i going to put in the session
	// register the type of data we want to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	// create a log print out in console window, with INFO prefix, and log.Ldate | log.Ltime
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	// create a log print out in console window, with ERROR prefix, and log.Ldate | log.Ltime | log.Lshortfile
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

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
		return err
	}
	app.TemplateCache = tc
	app.UseCache = false

	// Create a new handler repository
	repo := handlers.NewRepo(&app)
	// Set the handler repository for the handlers
	handlers.NewHandlers(repo)

	// Set the new render app config for the template package
	render.NewTemplates(&app)

	// setup helpers
	helpers.NewHelpers(&app)
	return nil
}
