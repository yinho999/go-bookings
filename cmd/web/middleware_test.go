package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(myH)
	switch v := h.(type) {
	case http.Handler:
		// do nothing - test passed
	default:
		// print out the type of v
		t.Error(fmt.Sprintf("type is not http.Handler, but %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	s := SessionLoad(myH)
	switch v := s.(type) {
	case http.Handler:
		// do nothing - test passed
	default:
		// print out the type of v
		t.Error(fmt.Sprintf("type is not http.Handler, but %T", v))
	}
}
