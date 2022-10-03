package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Before running the tests, do something then exit
	os.Exit(m.Run())
}

type myHandler struct{}

func (m myHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

}
