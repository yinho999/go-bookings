package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/do", nil)
	form := New(r.PostForm)

	if !form.Valid() {
		t.Error("got Valid() false, want true")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/do", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("got Valid() true, want false")
	}

	if form.Errors.Get("a") != "This field cannot be blank" {
		t.Error("expected a required error for 'a'")
	}

	if form.Errors.Get("b") != "This field cannot be blank" {
		t.Error("expected a required error for 'b'")
	}

	if form.Errors.Get("c") != "This field cannot be blank" {
		t.Error("expected a required error for 'c'")
	}

	postData := url.Values{}
	r, _ = http.NewRequest("POST", "/do", nil)
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")
	r.PostForm = postData
	form = New(r.PostForm)

	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("got Valid() false, want true")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/do", nil)
	form := New(r.PostForm)

	has := form.Has("a")
	if has {
		t.Error("got Has() true, want false")
	}

	postData := url.Values{}
	r, _ = http.NewRequest("POST", "/do", nil)
	postData.Add("a", "a")
	r.PostForm = postData
	form = New(r.PostForm)

	has = form.Has("a")
	if !has {
		t.Error("got Has() false, want true")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/do", nil)
	form := New(r.PostForm)

	form.MinLength("a", 10)

	if form.Valid() {
		t.Error("got Valid() true, want false")
	}

	if form.Errors.Get("a") != "This field is too short (minimum is 10 characters)" {
		t.Error("expected a min length error for 'a'")
	}

	r, _ = http.NewRequest("POST", "/do", nil)
	postData := url.Values{}
	r.PostForm = postData
	form = New(r.PostForm)
	postData.Add("a", "1234567890")
	form.MinLength("a", 1000)

	if form.Valid() {
		t.Error("got Valid() true, want false")
	}

	if form.Errors.Get("a") != "This field is too short (minimum is 1000 characters)" {
		t.Error("expected a min length error for 'a'")
	}
	postData = url.Values{}
	r.PostForm = postData
	form = New(r.PostForm)
	postData.Add("abc", "1234567890")
	form.MinLength("abc", 10)

	if !form.Valid() {
		t.Error("got Valid() false, want true")
	}

	if form.Errors.Get("abc") != "" {
		t.Error("got error for valid field")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	form.IsEmail("a")

	if form.Valid() {
		t.Error("form shows valid on non existent field")
	}

	if form.Errors.Get("a") != "Invalid email address" {
		t.Error("expected an email error for 'a'")
	}

	postData = url.Values{}

	postData.Add("a", "asdf@gmail.com")

	form = New(postData)

	form.IsEmail("a")

	if !form.Valid() {
		t.Error("got invalid email for valid email")
	}

	postData = url.Values{}
	postData.Add("a", "asdf")
	form = New(postData)

	form.IsEmail("a")

	if form.Valid() {
		t.Error("got valid for invalid email address")
	}
}
