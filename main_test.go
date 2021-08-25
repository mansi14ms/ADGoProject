package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//Testing Home Page
func TestHomePage(t *testing.T) {

	req, err := http.NewRequest("GET", "/homePage", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homePage)

	handler.ServeHTTP(rr, req)
	expected := "Welcome to the HomePage!"

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
