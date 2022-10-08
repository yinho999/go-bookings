package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/yinho999/go-bookings/internal/config"
	"github.com/yinho999/go-bookings/internal/models"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	/* From main.go run() */
	// what am i going to put in the session
	// register the type of data we want to put in the session
	gob.Register(models.Reservation{})
	// create a log print out in console window, with INFO prefix, and log.Ldate | log.Ltime
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	// create a log print out in console window, with ERROR prefix, and log.Ldate | log.Ltime | log.Lshortfile
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	// change this to true when in production
	testApp.InProduction = false

	// Initialize the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour // 24 hours
	// Setting the session cookie
	// Persist is set to true so that the cookie is stored in the browser
	// even the browser is closed
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false // Only send the cookie over HTTPS, we dont need that in development

	testApp.Session = session

	// in render.go
	app = &testApp

	// just before it exit, run our test
	os.Exit(m.Run())
}

// Implement the http.ResponseWriter interface into myWriter
type myWriter struct{}

func (m *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (m *myWriter) Write(bytes []byte) (int, error) {
	length := len(bytes)
	return length, nil
}

func (m *myWriter) WriteHeader(statusCode int) {

}
