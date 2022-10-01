package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
)

// Form creates a custom form struct which embeds an url.Values object.
type Form struct {
	// url.Values is a map[string][]string
	url.Values
	// Errors is a map[string]string
	Errors errors
}

// New initializes a custom form struct.
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if the form has a particular field in it.
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

// Valid returns true if there are no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Required checks for required fields and returns an error if not present.
// This is a method on the custom form struct.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MinLength checks for a minimum length of a field and returns an error if
// the field is too short.
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if len(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

// IsEmail checks for a valid email address and returns an error if the
// email address is not valid.
func (f *Form) IsEmail(field string) {
	value := f.Get(field)
	if !govalidator.IsEmail(value) {
		f.Errors.Add(field, "Invalid email address")
	}
}
