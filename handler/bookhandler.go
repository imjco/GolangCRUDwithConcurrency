package handler

import (
	"companyXyzProject/Model"
	"companyXyzProject/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

func AddBook(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	// Create a channel to communicate the result
	resultChan := make(chan string)

	// Start a goroutine to handle the addition of the book
	go func() {
		// Parse the JSON data from the request body
		var newBook Model.Book
		err := json.NewDecoder(r.Body).Decode(&newBook)
		defer r.Body.Close()

		if err != nil {
			resultChan <- fmt.Sprintf("Error: %s", err.Error())
			return
		}

		if err := ValidateBookAdd(newBook); err != nil {
			resultChan <- fmt.Sprintf("Error: %s", err)
			return
		}

		// Generate a unique ID for the new book
		newBook.ID = len(Model.DB.Books) + 1

		var bookauthor []Model.BookAuthor
		for _, author := range newBook.Authors {

			bookauthor = append(bookauthor, Model.BookAuthor{
				AuthorID: author.ID,
				BookID:   len(Model.DB.Books) + 1,
			})
		}

		var bookpublisher []Model.BookPublisher
		for _, publisher := range newBook.Publisher {
			bookpublisher = append(bookpublisher, Model.BookPublisher{
				PublisherID: publisher.ID,
				BookID:      len(Model.DB.Books) + 1,
			})
		}

		// Lock the mutex to ensure safe access to the database
		Model.MU.Lock()
		defer Model.MU.Unlock()

		// Add the new book to the database
		Model.DB.Books = append(Model.DB.Books, newBook)
		Model.DB.BookAuthors = append(Model.DB.BookAuthors, bookauthor...)
		Model.DB.BookPublishers = append(Model.DB.BookPublishers, bookpublisher...)
		// Call the function to append to the CSV file

		if err := service.AppendToCSV(Model.CsvFileName, newBook); err != nil {
			log.Fatal("Error appending to CSV file:", err)
		}

		// Send a success message to the channel
		resultChan <- "Book added successfully"
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
func GetBooks(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)

	//enableCors(&w)

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	startIndex, endIndex, err := ValidatePageLimit(pageStr, limitStr, len(Model.DB.Books))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authorChan := make(chan []Model.Author)
	bookChan := make(chan []Model.Book)
	publisherChan := make(chan []Model.Publisher)

	booksTotalCount := len(Model.DB.Books)

	go func() {
		Model.MU.Lock()
		defer Model.MU.Unlock()

		// Fetch Authors
		authorChan <- Model.DB.Authors
		// Fetch Books
		bookChan <- Model.DB.Books[startIndex:endIndex]
		// Fetch Publishers
		publisherChan <- Model.DB.Publishers
	}()

	authors := <-authorChan
	books := <-bookChan
	publishers := <-publisherChan

	for i := range books {
		authorIDs := []int{}
		for _, bookAuthor := range Model.DB.BookAuthors {
			if bookAuthor.BookID == books[i].ID {
				authorIDs = append(authorIDs, bookAuthor.AuthorID)
			}
		}
		relatedAuthors := []Model.Author{}
		for _, authorID := range authorIDs {
			for _, author := range authors {
				if author.ID == authorID {
					relatedAuthors = append(relatedAuthors, author)
				}
			}
		}

		books[i].Authors = relatedAuthors
	}

	// Fetch related publishers for each book using BookPublisher
	for i := range books {
		publisherIDs := []int{}
		for _, bookPublisher := range Model.DB.BookPublishers {
			if bookPublisher.BookID == books[i].ID {
				publisherIDs = append(publisherIDs, bookPublisher.PublisherID)
			}
		}

		relatedPublishers := []Model.Publisher{}
		for _, publisherID := range publisherIDs {
			for _, publisher := range publishers {
				if publisher.ID == publisherID {
					relatedPublishers = append(relatedPublishers, publisher)
				}
			}
		}

		books[i].Publisher = relatedPublishers
	}

	response := Model.BooksResponse{
		Books:           books,
		BooksTotalCount: booksTotalCount,
	}

	// Encode and return the list of books as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func GetBookByISBN13(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	vars := mux.Vars(r)

	isbn13 := vars["isbn13"]

	// Use a channel to communicate the found book
	authorChan := make(chan []Model.Author)
	bookChan := make(chan []Model.Book)
	publisherChan := make(chan []Model.Publisher)

	go func() {
		Model.MU.Lock()
		defer Model.MU.Unlock()

		// Find the book by ISBN13
		var foundBook *Model.Book
		for i := range Model.DB.Books {
			if Model.DB.Books[i].ISBN13 == isbn13 {
				foundBook = &Model.DB.Books[i]
				break
			}
		}
		if foundBook == nil {
			bookChan <- nil // Only send the found book if it's not nil
		} else {
			bookSlice := []Model.Book{*foundBook}
			bookChan <- bookSlice
		}

		authorChan <- Model.DB.Authors
		publisherChan <- Model.DB.Publishers

	}()

	// Wait for the goroutine to complete
	books := <-bookChan
	authors := <-authorChan
	publishers := <-publisherChan
	booksTotalCount := len(books)
	// Check if foundBook is nil
	if len(books) == 0 {
		// Book not found, return an error response
		http.NotFound(w, r)
		return
	}
	for i := range books {
		authorIDs := []int{}
		for _, bookAuthor := range Model.DB.BookAuthors {
			if bookAuthor.BookID == books[i].ID {
				authorIDs = append(authorIDs, bookAuthor.AuthorID)
			}
		}
		relatedAuthors := []Model.Author{}
		for _, authorID := range authorIDs {
			for _, author := range authors {
				if author.ID == authorID {
					relatedAuthors = append(relatedAuthors, author)
				}
			}
		}

		books[i].Authors = relatedAuthors
	}

	// Fetch related publishers for each book using BookPublisher
	for i := range books {
		publisherIDs := []int{}
		for _, bookPublisher := range Model.DB.BookPublishers {
			if bookPublisher.BookID == books[i].ID {
				publisherIDs = append(publisherIDs, bookPublisher.PublisherID)
			}
		}

		relatedPublishers := []Model.Publisher{}
		for _, publisherID := range publisherIDs {
			for _, publisher := range publishers {
				if publisher.ID == publisherID {
					relatedPublishers = append(relatedPublishers, publisher)
				}
			}
		}

		books[i].Publisher = relatedPublishers
	}

	if books == nil {
		http.NotFound(w, r)
		return
	}

	response := Model.BooksResponse{
		Books:           books,
		BooksTotalCount: booksTotalCount,
	}

	// Encode and return the found book as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	vars := mux.Vars(r)
	isbn13 := vars["isbn13"]

	// Create a new book based on the request data
	var updatedBook Model.Book
	err := json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Use a channel to communicate the result
	resultChan := make(chan string)

	idsToRemove := []int{}
	for _, publisher := range Model.DB.BookPublishers {
		if publisher.BookID == updatedBook.ID {
			idsToRemove = append(idsToRemove, publisher.PublisherID)
		}
	}
	idsToRemoveAuthor := []int{}
	for _, authors := range Model.DB.BookAuthors {
		if authors.BookID == updatedBook.ID {
			idsToRemoveAuthor = append(idsToRemoveAuthor, authors.AuthorID)
		}
	}

	// Remove the book publishers with the given IDs from the slice.
	newBookPublishers := service.RemoveBookPublishersWithIDs(Model.DB.BookPublishers, idsToRemove)
	// Remove the book authors with the given IDs from the slice.
	newBookAuthors := service.RemoveBookPublishersWithIDsAuthors(Model.DB.BookAuthors, idsToRemoveAuthor)

	go func() {
		Model.MU.Lock()
		defer Model.MU.Unlock()

		// Find the book by ISBN13
		for i := range Model.DB.Books {

			if Model.DB.Books[i].ISBN13 == isbn13 && Model.DB.Books[i].ID == updatedBook.ID {

				// Update the book with the data from the request
				Model.DB.Books[i].Title = updatedBook.Title
				Model.DB.Books[i].ISBN10 = updatedBook.ISBN10
				Model.DB.Books[i].ListPrice = updatedBook.ListPrice
				Model.DB.Books[i].PublicationYear = updatedBook.PublicationYear
				Model.DB.Books[i].ImageURL = updatedBook.ImageURL
				Model.DB.Books[i].Edition = updatedBook.Edition

				bookauthor := []Model.BookAuthor{}
				for x := range updatedBook.Authors {
					bookauthor = append(bookauthor, Model.BookAuthor{
						AuthorID: updatedBook.Authors[x].ID,
						BookID:   updatedBook.ID,
					})
				}
				Model.DB.BookAuthors = append(newBookAuthors, bookauthor...)

				bookpublisher2 := []Model.BookPublisher{}
				for y := range updatedBook.Publisher {
					bookpublisher2 = append(bookpublisher2, Model.BookPublisher{
						PublisherID: updatedBook.Publisher[y].ID,
						BookID:      updatedBook.ID,
					})
				}
				Model.DB.BookPublishers = append(newBookPublishers, bookpublisher2...)

				// Send a success message to the channel
				resultChan <- "Book updated successfully"
				return
			}
		}

		// If the book is not found, send an error message
		resultChan <- "Book not found"
	}()

	// Wait for the goroutine to complete
	result := <-resultChan

	// Check the result and respond accordingly
	if result == "Book not found" {
		http.NotFound(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)

	vars := mux.Vars(r)
	isbn13 := vars["isbn13"]

	// Use a channel to communicate the result
	resultChan := make(chan string)

	go func() {
		Model.MU.Lock()
		defer Model.MU.Unlock()

		// Find the book by ISBN13
		for i := range Model.DB.Books {
			if Model.DB.Books[i].ISBN13 == isbn13 {
				// Remove the book from the database
				Model.DB.Books = append(Model.DB.Books[:i], Model.DB.Books[i+1:]...)

				// Send a success message to the channel
				resultChan <- "Book deleted successfully"
				return
			}
		}

		// If the book is not found, send an error message
		resultChan <- "Book not found"
	}()

	// Wait for the goroutine to complete
	result := <-resultChan

	// Check the result and respond accordingly
	if result == "Book not found" {
		http.NotFound(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}
