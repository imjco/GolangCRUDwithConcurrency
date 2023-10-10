package handler_test

import (
	"companyXyzProject/handler"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetISBN13TO10(t *testing.T) {

	r := mux.NewRouter()
	r.HandleFunc("/convertisbn13to10/{isbn13}", handler.GetISBN13TO10).Methods("GET")

	// Create a request with an ISBN13 parameter
	req, err := http.NewRequest("GET", "/convertisbn13to10/9967615110931", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "7615110939"

	// decode the response body to struct
	type response struct {
		Isbn10 string `json:"isbn10"`
	}

	var isbn10 response
	err = json.NewDecoder(rr.Body).Decode(&isbn10)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if isbn10.Isbn10 != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func TestGetISBN10TO13(t *testing.T) {

	r := mux.NewRouter()
	r.HandleFunc("/convertisbn10to13/{isbn10}", handler.GetISBN10TO13).Methods("GET")

	// Create a request with an ISBN10 parameter
	req, err := http.NewRequest("GET", "/convertisbn10to13/1891830856", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "9781891830853"

	type response struct {
		Isbn13 string `json:"isbn13"`
	}

	var res response
	err = json.NewDecoder(rr.Body).Decode(&res)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if res.Isbn13 != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
