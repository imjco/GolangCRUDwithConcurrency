package handler

import (
	"companyXyzProject/Model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func AddAuthor(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	// Create a channel to communicate the result
	resultChan := make(chan string)

	// Start a goroutine to handle the addition of the book
	go func() {
		// Parse the JSON data from the request body
		var newAuthor Model.Author
		err := json.NewDecoder(r.Body).Decode(&newAuthor)
		defer r.Body.Close()
		if err != nil {
			resultChan <- fmt.Sprintf("Error: %s", err.Error())
			return
		}

		// Check if the book already exists
		for _, book := range Model.DB.Authors {
			if book.FirstName == newAuthor.FirstName && book.LastName == newAuthor.LastName {
				resultChan <- "Error: Author already exists"
				return
			}
		}

		// Generate a unique ID for the new book
		newAuthor.ID = len(Model.DB.Authors) + 1

		// Lock the mutex to ensure safe access to the database
		Model.MU.Lock()
		defer Model.MU.Unlock()

		// Add the new book to the database
		Model.DB.Authors = append(Model.DB.Authors, newAuthor)

		// Send a success message to the channel
		resultChan <- "Author added successfully"
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
func GetAuthor(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	startIndex, endIndex, err := ValidatePageLimit(pageStr, limitStr, len(Model.DB.Authors))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Use a channel to communicate the list of author
	resultChan := make(chan []Model.Author)

	go func() {
		Model.MU.Lock()
		defer Model.MU.Unlock()

		// Send the list of author to the channel
		resultChan <- Model.DB.Authors[startIndex:endIndex]
	}()

	// Wait for the goroutine to complete
	books := <-resultChan

	// Encode and return the list of author as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}
func GetAuthorByID(w http.ResponseWriter, r *http.Request) {
	//	enableCors(&w)
	vars := mux.Vars(r)
	// Convert the ID into an integer
	findID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Use a channel to communicate the found author
	resultChan := make(chan *Model.Author)

	go func() {
		Model.MU.Lock()
		defer Model.MU.Unlock()

		// Find the book by ID
		var foundBook *Model.Author
		for i := range Model.DB.Authors {
			if Model.DB.Authors[i].ID == findID {
				foundBook = &Model.DB.Authors[i]
				break
			}
		}

		// Send the found book to the channel
		resultChan <- foundBook
	}()

	// Wait for the goroutine to complete
	foundAuthor := <-resultChan

	if foundAuthor == nil {
		http.NotFound(w, r)
		return
	}

	// Encode and return the found book as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(foundAuthor)
}
func UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	vars := mux.Vars(r)
	// Convert the ID into an integer
	findID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new Author based on the request data
	var updatedAuthor Model.Author
	err = json.NewDecoder(r.Body).Decode(&updatedAuthor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Use a channel to communicate the result
	resultChan := make(chan string)
	go func() {
		Model.MU.Lock()
		defer Model.MU.Unlock()

		// Find the Author by ID
		for i := range Model.DB.Authors {
			if Model.DB.Authors[i].ID == findID {
				// Update the publisher with the data from the request
				Model.DB.Authors[i].FirstName = updatedAuthor.FirstName
				Model.DB.Authors[i].LastName = updatedAuthor.LastName
				Model.DB.Authors[i].MiddleName = updatedAuthor.MiddleName

				// Send a success message to the channel
				resultChan <- "Author updated successfully"
				return
			}
		}

		// If the Author is not found, send an error message
		resultChan <- "Author not found"
	}()

	// Wait for the goroutine to complete
	result := <-resultChan

	// Check the result and respond accordingly
	if result == "Author not found" {
		http.NotFound(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}
func DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
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

		// Find the Author by ID
		for i := range Model.DB.Authors {
			if Model.DB.Authors[i].ID == findID {
				// Remove the Author from the database
				Model.DB.Authors = append(Model.DB.Authors[:i], Model.DB.Authors[i+1:]...)

				// Send a success message to the channel
				resultChan <- "Author deleted successfully"
				return
			}
		}

		// If the Author is not found, send an error message
		resultChan <- "Author not found"

	}()

	// Wait for the goroutine to complete
	result := <-resultChan

	// Check the result and respond accordingly
	if result == "Author not found" {
		http.NotFound(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}
