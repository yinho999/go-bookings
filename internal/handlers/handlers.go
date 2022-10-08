package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/yinho999/go-bookings/internal/config"
	"github.com/yinho999/go-bookings/internal/forms"
	"github.com/yinho999/go-bookings/internal/helpers"
	"github.com/yinho999/go-bookings/internal/models"
	"github.com/yinho999/go-bookings/internal/render"
	"net/http"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	// Store remote Ip address
	//remoteIP := r.RemoteAddr
	// Store into the session
	//m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{}, r)
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//// perform some logic
	//stringMap := make(map[string]string)
	//stringMap["test"] = "Hello, again."
	//
	//// retrieve data from the session
	//remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	//stringMap["remote_ip"] = remoteIP

	// send the data to the template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		//StringMap: stringMap,
	}, r)
}

// Reservation renders the make a reservation page and display form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	// Create reservation and send to template
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	}, r)
}

// PostReservation handles the submission of reservation form
func (m *Repository) PostReservation(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		//log.Println(err)
		helpers.ServerError(writer, err)
		return
	}
	reservation := models.Reservation{
		FirstName: request.Form.Get("first_name"),
		LastName:  request.Form.Get("last_name"),
		Email:     request.Form.Get("email"),
		Phone:     request.Form.Get("phone"),
	}
	// form validation
	form := forms.New(request.PostForm)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	// If errors, redisplay the reservation page
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(writer, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, request)
		return
	}
	// pass the reservation data to session
	m.App.Session.Put(request.Context(), "reservation", reservation)
	http.Redirect(writer, request, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "generals.page.tmpl", &models.TemplateData{}, r)
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "majors.page.tmpl", &models.TemplateData{}, r)
}

// Availability renders the search availability page and display form
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.tmpl", &models.TemplateData{}, r)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.tmpl", &models.TemplateData{}, r)
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// get the start date
	start := r.Form.Get("start")
	// get the end date
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		Ok:      true,
		Message: "Available!",
	}
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// get the reservation data from the session
	// .models.Reservation is the type of the data
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get item from session")
		// redirect to make reservation page
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// remove the reservation data from the session after we get it
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	}, r)
}
