package handler_test

import (
	"bytes"
	"companyXyzProject/Model"
	"companyXyzProject/handler"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddAuthor(t *testing.T) {
	// Create a test request body
	newAuthor := Model.Author{
		FirstName:  "John",
		LastName:   "Doe",
		MiddleName: "M",
	}
	body, err := json.Marshal(newAuthor)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/authors", handler.AddAuthor).Methods("POST")
	// Create a test request
	req, err := http.NewRequest("POST", "/authors", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expected := "Author added successfully"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("Expected response body %q, but got %q", expected, rr.Body.String())
	}

	// Check that the author was added to the database
	if len(Model.DB.Authors) != 1 {
		t.Errorf("Expected 1 author in the database, but got %d", len(Model.DB.Authors))
	}

	if Model.DB.Authors[0].FirstName != newAuthor.FirstName {
		t.Errorf("Expected author first name %q, but got %q", newAuthor.FirstName, Model.DB.Authors[0].FirstName)
	}

	if Model.DB.Authors[0].LastName != newAuthor.LastName {
		t.Errorf("Expected author last name %q, but got %q", newAuthor.LastName, Model.DB.Authors[0].LastName)
	}

	if Model.DB.Authors[0].MiddleName != newAuthor.MiddleName {
		t.Errorf("Expected author middle name %q, but got %q", newAuthor.MiddleName, Model.DB.Authors[0].MiddleName)
	}
}

func TestGetAuthor(t *testing.T) {

	for i := 0; i < 10; i++ {
		auth := Model.Author{
			ID:         len(Model.DB.Authors) + 1,
			FirstName:  "John",
			LastName:   "Doe",
			MiddleName: "M",
		}
		Model.DB.Authors = append(Model.DB.Authors, auth)
	}

	r := mux.NewRouter()
	r.HandleFunc("/authors", handler.GetAuthor).Methods("GET")
	//// Create a test request with query parameters
	req, err := http.NewRequest("GET", "/authors?page=1&limit=10", nil)

	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	//
	//// Create a test response recorder
	rr := httptest.NewRecorder()

	//
	//// Serve the request using the router
	r.ServeHTTP(rr, req)

	//
	//// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}
	//
	//// Check the response body
	var Authors []Model.Author
	err = json.Unmarshal(rr.Body.Bytes(), &Authors)
	if err != nil {
		t.Fatalf("Failed to decode response body as JSON: %v", err)
	}

	expectedCount := 10
	if len(Authors) != expectedCount {
		t.Errorf("Expected %d publishers, but got %d", expectedCount, len(Authors))
	}
}

func TestGetAuthorByID(t *testing.T) {
	// Add some test authors to the database
	Model.DB.Authors = []Model.Author{
		{
			ID:         1,
			FirstName:  "John",
			LastName:   "Doe",
			MiddleName: "M",
		},
		{
			ID:         2,
			FirstName:  "Jane",
			LastName:   "Doe",
			MiddleName: "N",
		},
	}
	r := mux.NewRouter()
	r.HandleFunc("/authors/{id}", handler.GetAuthorByID).Methods("GET")
	// Create a test request
	req, err := http.NewRequest("GET", "/authors/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	var author Model.Author
	err = json.Unmarshal(rr.Body.Bytes(), &author)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if author.ID != 1 {
		t.Errorf("Expected author ID 1, but got %d", author.ID)
	}

	if author.FirstName != "John" {
		t.Errorf("Expected author first name John, but got %q", author.FirstName)
	}

	if author.LastName != "Doe" {
		t.Errorf("Expected author last name Doe, but got %q", author.LastName)
	}

	if author.MiddleName != "M" {
		t.Errorf("Expected author middle name M, but got %q", author.MiddleName)
	}
}

func TestUpdateAuthor(t *testing.T) {
	// Add some test authors to the database
	Model.DB.Authors = []Model.Author{
		{
			ID:         1,
			FirstName:  "John",
			LastName:   "Doe",
			MiddleName: "M",
		},
		{
			ID:         2,
			FirstName:  "Jane",
			LastName:   "Doe",
			MiddleName: "N",
		},
	}
	r := mux.NewRouter()
	r.HandleFunc("/authors/{id}", handler.UpdateAuthor).Methods("PUT")
	// Create a test request body
	updatedAuthor := Model.Author{
		ID:         1,
		FirstName:  "John",
		LastName:   "Doe",
		MiddleName: "P",
	}

	body, err := json.Marshal(updatedAuthor)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create a test request
	req, err := http.NewRequest("PUT", "/authors/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expected := "Author updated successfully"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("Expected response body %q, but got %q", expected, rr.Body.String())
	}

	// Check that the author was updated in the database
	if Model.DB.Authors[0].MiddleName != updatedAuthor.MiddleName {
		t.Errorf("Expected author middle name %q, but got %q", updatedAuthor.MiddleName, Model.DB.Authors[0].MiddleName)
	}
}

func TestDeleteAuthor(t *testing.T) {
	// Add some test authors to the database
	Model.DB.Authors = []Model.Author{
		{
			ID:         1,
			FirstName:  "John",
			LastName:   "Doe",
			MiddleName: "M",
		},
		{
			ID:         2,
			FirstName:  "Jane",
			LastName:   "Doe",
			MiddleName: "N",
		},
	}
	r := mux.NewRouter()
	r.HandleFunc("/authors/{id}", handler.DeleteAuthor).Methods("DELETE")
	// Create a test request
	req, err := http.NewRequest("DELETE", "/authors/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expected := "Author deleted successfully"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("Expected response body %q, but got %q", expected, rr.Body.String())
	}

	// Check that the author was deleted from the database
	if len(Model.DB.Authors) != 1 {
		t.Errorf("Expected 1 author in the database, but got %d", len(Model.DB.Authors))
	}

	if Model.DB.Authors[0].ID != 2 {
		t.Errorf("Expected author ID 2, but got %d", Model.DB.Authors[0].ID)
	}
}
