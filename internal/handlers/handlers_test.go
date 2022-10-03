package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// struct of test cases
type postData struct {
	key   string
	value string
}

// Test table of each handler
var theTest = []struct {
	name         string
	url          string
	method       string
	params       []postData
	expectedCode int
}{
	{
		"home", "/", "GET", []postData{}, http.StatusOK,
	}, {
		"about", "/about", "GET", []postData{}, http.StatusOK,
	}, {
		"contact", "/contact", "GET", []postData{}, http.StatusOK,
	}, {
		"general quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK,
	}, {
		"major suite", "/majors-suite", "GET", []postData{}, http.StatusOK,
	}, {
		"search availability", "/search-availability", "GET", []postData{}, http.StatusOK,
	}, {
		"make reservation", "/make-reservation", "GET", []postData{}, http.StatusOK,
	}, {
		"post-search-availability", "/search-availability", "POST", []postData{
			{"start", "2021-01-01"},
			{"end", "2021-01-02"},
		}, http.StatusOK,
	},
	{
		"post-search-availability-json", "/search-availability-json", "POST", []postData{
			{"start", "2021-01-01"},
			{"end", "2021-01-02"},
		}, http.StatusOK,
	}, {
		"post-make-reservation", "/make-reservation", "POST", []postData{
			{"first_name", "John"},
			{"last_name", "Smith"},
			{"email", "me@gmail.com"},
			{"phone", "123456789"},
		}, http.StatusOK,
	},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	// create a test server
	ts := httptest.NewTLSServer(routes)
	// close the server when the test finishes
	defer ts.Close()

	for _, e := range theTest {
		if e.method == "GET" {
			// create a get request
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Error(err)
			}
			if resp.StatusCode != e.expectedCode {
				t.Errorf("Expected %d, got %d", e.expectedCode, resp.StatusCode)
			}
		} else if e.method == "POST" {
			// create a url post form data from the params
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			// create a post request
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Error(err)
			}
			if resp.StatusCode != e.expectedCode {
				t.Errorf("Expected %d, got %d", e.expectedCode, resp.StatusCode)
			}

		}
	}
}
