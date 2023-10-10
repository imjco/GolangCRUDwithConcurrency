package handler

import (
	"companyXyzProject/Model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func AddPublisher(w http.ResponseWriter, r *http.Request) {
	//	enableCors(&w)
	// Create a channel to communicate the result
	resultChan := make(chan string)
	go func() {
		// Parse the JSON data from the request body
		var newPublisher Model.Publisher
		err := json.NewDecoder(r.Body).Decode(&newPublisher)
		defer r.Body.Close()
		if err != nil {
			resultChan <- fmt.Sprintf("Error: %s", err.Error())
			return
		}

		// check if publisher is already exists
		for _, publisher := range Model.DB.Publishers {
			if publisher.Name == newPublisher.Name {
				resultChan <- "Error: Publisher already exists"
				return
			}
		}

		// Generate a unique ID for the new publisher
		newPublisher.ID = len(Model.DB.Publishers) + 1
		// Lock the mutex to ensure safe access to the database
		Model.MU.Lock()
		defer Model.MU.Unlock()

		// Add the new book to the database
		Model.DB.Publishers = append(Model.DB.Publishers, newPublisher)

		// Send a success message to the channel
		resultChan <- "Publishers added successfully"
	}()

	// Wait for the goroutine to complete and get the result from the channel
	result := <-resultChan

	// Check the result and respond accordingly
	if strings.HasPrefix(result, "Error:") {
		http.Error(w, result, http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}

}
func GetPublisher(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	startIndex, endIndex, err := ValidatePageLimit(pageStr, limitStr, len(Model.DB.Publishers))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Use a channel to communicate the list of books
	resultChan := make(chan []Model.Publisher)

	go func() {
		Model.MU.Lock()
		defer Model.MU.Unlock()

		// Send the list of books to the channel
		resultChan <- Model.DB.Publishers[startIndex:endIndex]

	}()

	// Wait for the goroutine to complete
	publisher := <-resultChan

	// Encode and return the list of books as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(publisher)
}
func GetPublisherByID(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	vars := mux.Vars(r)

	// Convert the ID into an integer
	findID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Use a channel to communicate the found publisher
	resultChan := make(chan *Model.Publisher)

	go func() {
		Model.MU.Lock()
		defer Model.MU.Unlock()

		// Find the Publisher by ID
		var foundPublisher *Model.Publisher
		for i := range Model.DB.Publishers {
			if Model.DB.Publishers[i].ID == findID {
				foundPublisher = &Model.DB.Publishers[i]
				break
			}
		}

		// Send the found book to the channel
		resultChan <- foundPublisher
	}()

	// Wait for the goroutine to complete
	foundBook := <-resultChan

	if foundBook == nil {
		http.NotFound(w, r)
		return
	}

	// Encode and return the found book as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(foundBook)
}
func UpdatePublisher(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	vars := mux.Vars(r)
	// Convert the ID into an integer

	findID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new publisher based on the request data
	var updatedPublisher Model.Publisher
	err = json.NewDecoder(r.Body).Decode(&updatedPublisher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Use a channel to communicate the result
	resultChan := make(chan string)
	go func() {
		Model.MU.Lock()
		defer Model.MU.Unlock()

		// Find the book by ISBN13
		for i := range Model.DB.Publishers {
			if Model.DB.Publishers[i].ID == findID {
				// Update the publisher with the data from the request
				Model.DB.Publishers[i].Name = updatedPublisher.Name

				// Send a success message to the channel
				resultChan <- "Publisher updated successfully"
				return
			}
		}

		// If the book is not found, send an error message
		resultChan <- "Publisher not found"
	}()

	// Wait for the goroutine to complete
	result := <-resultChan

	// Check the result and respond accordingly
	if result == "Publisher not found" {
		http.NotFound(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}
func DeletePublisher(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// Convert the ID into an integer
	findID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Use a channel to communicate the result
	resultChan := make(chan string)
	go func() {
		Model.MU.Lock()
		defer Model.MU.Unlock()

		// Find the book by ISBN13
		for i := range Model.DB.Publishers {
			if Model.DB.Publishers[i].ID == findID {
				// Remove the book from the database
				Model.DB.Publishers = append(Model.DB.Publishers[:i], Model.DB.Publishers[i+1:]...)

				// Send a success message to the channel
				resultChan <- "Publisher deleted successfully"
				return
			}
		}

		// If the Publisher is not found, send an error message
		resultChan <- "Publisher not found"
	}()

	// Wait for the goroutine to complete
	result := <-resultChan

	// Check the result and respond accordingly
	if result == "Publisher not found" {
		http.NotFound(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}
