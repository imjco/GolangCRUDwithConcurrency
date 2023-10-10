package handler_test

import (
	"bytes"
	"companyXyzProject/Model"
	"companyXyzProject/handler"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestAddPublisher(t *testing.T) {
	// Create a test publisher
	newPublisher := Model.Publisher{
		Name: "Test PublisherWWWW",
	}

	// Encode the publisher as JSON
	jsonData, err := json.Marshal(newPublisher)
	if err != nil {
		t.Fatalf("Failed to encode publisher as JSON: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/publishers", handler.AddPublisher).Methods("POST")

	// Create a test request with the JSON data
	req, err := http.NewRequest("POST", "/publishers", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}
	fmt.Print(rr.Code)
	// Check the response body
	expected := "Publishers added successfully"
	if rr.Body.String() != expected {
		t.Errorf("Expected response body %q, but got %q", expected, rr.Body.String())
	}
}

func TestGetPublisher(t *testing.T) {

	for i := 0; i < 10; i++ {
		publisher := Model.Publisher{
			ID:   len(Model.DB.Publishers) + 1,
			Name: "Test" + fmt.Sprint(i),
		}
		Model.DB.Publishers = append(Model.DB.Publishers, publisher)
	}

	// Create a test request with query parameters
	req, err := http.NewRequest("GET", "/publishers?page=1&limit=10", nil)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Call the GetPublisher handler function
	handler.GetPublisher(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	var publishers []Model.Publisher
	err = json.Unmarshal(rr.Body.Bytes(), &publishers)
	if err != nil {
		t.Fatalf("Failed to decode response body as JSON: %v", err)
	}

	expectedCount := 10
	if len(publishers) != expectedCount {
		t.Errorf("Expected %d publishers, but got %d", expectedCount, len(publishers))
	}
}

func TestGetPublisherByID(t *testing.T) {
	newPublisher := Model.Publisher{
		ID:   len(Model.DB.Publishers) + 1,
		Name: "Test Publisher",
	}
	// Add the publisher to the database
	Model.DB.Publishers = append(Model.DB.Publishers, newPublisher)

	r := mux.NewRouter()
	r.HandleFunc("/publishers/{id}", handler.GetPublisherByID).Methods("GET")

	// Create a test request with a specific ID
	req, err := http.NewRequest("GET", "/publishers/"+strconv.Itoa(newPublisher.ID), nil) // Replace "123" with the desired ID
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	var publisher Model.Publisher
	err = json.Unmarshal(rr.Body.Bytes(), &publisher)

	if err != nil {
		t.Fatalf("Failed to decode response body as JSON: %v", err)
	}
	if publisher.ID != newPublisher.ID {
		t.Errorf("Expected publisher ID %d, but got %d", newPublisher.ID, publisher.ID)
	}

}

func TestUpdatePublisher(t *testing.T) {
	// Create a test publisher
	newPublisher := Model.Publisher{
		ID:   len(Model.DB.Publishers) + 1,
		Name: "Test Publisher",
	}
	// Add the publisher to the database
	Model.DB.Publishers = append(Model.DB.Publishers, newPublisher)

	// Create a test request with the publisher ID and updated data
	updatedPublisher := Model.Publisher{
		ID:   Model.DB.Publishers[0].ID,
		Name: "Updated Publisher",
	}
	jsonData, err := json.Marshal(updatedPublisher)
	if err != nil {
		t.Fatalf("Failed to encode publisher as JSON: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/publishers/{id}", handler.UpdatePublisher).Methods("PUT")

	req, err := http.NewRequest("PUT", "/publishers/"+strconv.Itoa(updatedPublisher.ID), bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	expected := "Publisher updated successfully"
	if rr.Body.String() != expected {
		t.Errorf("Expected response body %q, but got %q", expected, rr.Body.String())
	}

}

func TestDeletePublisher(t *testing.T) {
	// Create a test publisher
	newPublisher := Model.Publisher{
		ID:   len(Model.DB.Publishers) + 1,
		Name: "Test Publisher",
	}

	// Add the publisher to the database
	Model.DB.Publishers = append(Model.DB.Publishers, newPublisher)

	r := mux.NewRouter()
	r.HandleFunc("/publishers/{id}", handler.DeletePublisher).Methods("DELETE")

	// Create a test request with the publisher ID
	req, err := http.NewRequest("DELETE", "/publishers/"+strconv.Itoa(newPublisher.ID), nil)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
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
	expected := "Publisher deleted successfully"
	if rr.Body.String() != expected {
		t.Errorf("Expected response body %q, but got %q", expected, rr.Body.String())
	}

	// Check that the publisher was deleted from the database
	foundPublisher := false
	for _, publisher := range Model.DB.Publishers {
		if publisher.ID == newPublisher.ID {
			foundPublisher = true
			break
		}
	}
	if foundPublisher {
		t.Errorf("Expected publisher %v to be deleted, but it was found in the database", newPublisher)
	}
}
