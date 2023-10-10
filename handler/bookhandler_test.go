package handler_test

import (
	"bytes"
	"companyXyzProject/Model"
	"companyXyzProject/handler"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestAddBook(t *testing.T) {
	// Create a test book
	book := Model.Book{
		Title:           "Test Book",
		ISBN13:          "9781891830853",
		ISBN10:          "1891830856",
		ListPrice:       9.99,
		PublicationYear: 2021,
		ImageURL:        "https://example.com/test.jpg",
		Edition:         "1st",
		Authors: []Model.Author{
			{
				ID:        1,
				FirstName: "Joel",
				LastName:  "Hartse",
			},
		},
		Publisher: []Model.Publisher{
			{
				ID:   1,
				Name: "Test",
			},
		},
	}
	Model.DB.Authors = []Model.Author{
		{
			ID:         1,
			FirstName:  "John",
			LastName:   "Doe",
			MiddleName: "M",
		},
	}
	Model.DB.Publishers = []Model.Publisher{
		{
			ID:   1,
			Name: "Test",
		},
	}
	Model.CsvFileName = "isbn_ean.csv"
	// Marshal the book to JSON
	jsonBook, err := json.Marshal(book)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/books", handler.AddBook).Methods("POST")
	// Create a request with the JSON book data
	req, err := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonBook))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected status code: got %v, expected %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Book added successfully"
	if rr.Body.String() != expected {
		t.Errorf("Unexpected response body: got %v, expected %v", rr.Body.String(), expected)
	}
}

func TestGetBooks(t *testing.T) {

	for i := 1; i <= 20; i++ {
		books := Model.Book{
			ID:              i,
			Title:           "Test Book",
			ISBN13:          "9781891830853",
			ISBN10:          "1891830856",
			ListPrice:       9.99,
			PublicationYear: 2021,
			ImageURL:        "https://example.com/test.jpg",
			Edition:         "1st",
		}
		Model.DB.Books = append(Model.DB.Books, books)
	}

	Model.DB.Authors = []Model.Author{
		{
			ID:         1,
			FirstName:  "John",
			LastName:   "Hartse",
			MiddleName: "M",
		},
	}
	Model.DB.Publishers = []Model.Publisher{
		{
			ID:   1,
			Name: "Test",
		},
	}

	// Create a request to get all books
	r := mux.NewRouter()
	r.HandleFunc("/books", handler.GetBooks).Methods("GET")

	req, err := http.NewRequest("GET", "/books?page=1&limit=10", nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected status code: got %v, expected %v", status, http.StatusOK)
	}

	// Check the response body
	type booksResponse struct {
		Books           []Model.Book
		BooksTotalCount int
	}

	var booksResponses booksResponse

	// json.unmarshal cannot unmarshal the response body
	err = json.Unmarshal([]byte(rr.Body.String()), &booksResponses)
	if err != nil {
		t.Fatalf("Failed to decode response body as JSON: %v", err)
	}

	expectedCount := 10
	if len(booksResponses.Books) != expectedCount {
		t.Errorf("Expected %d publishers, but got %d", expectedCount, len(booksResponses.Books))
	}
}

func TestGetBookByISBN13(t *testing.T) {
	books := Model.Book{
		ID:              1,
		Title:           "Test Book",
		ISBN13:          "9781891830853",
		ISBN10:          "1891830856",
		ListPrice:       9.99,
		PublicationYear: 2021,
		ImageURL:        "https://example.com/test.jpg",
		Edition:         "1st",
	}
	Model.DB.Books = append(Model.DB.Books, books)
	r := mux.NewRouter()
	r.HandleFunc("/books/{isbn13}", handler.GetBookByISBN13).Methods("GET")
	// Create a request to get a book by ISBN13
	req, err := http.NewRequest("GET", "/books/9781891830853", nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected status code: got %v, expected %v", status, http.StatusOK)
	}

	// Check the response body
	type booksResponse struct {
		Books           []Model.Book
		BooksTotalCount int
	}

	var booksResponses booksResponse

	// json.unmarshal cannot unmarshal the response body
	err = json.Unmarshal([]byte(rr.Body.String()), &booksResponses)
	if err != nil {
		t.Fatalf("Failed to decode response body as JSON: %v", err)
	}

	if len(booksResponses.Books) != 1 {
		t.Errorf("Unexpected number of books: got %v, expected %v", len(booksResponses.Books), 1)
	}

	if booksResponses.Books[0].ISBN13 != "9781891830853" {
		t.Errorf("Unexpected ISBN13: got %v, expected %v", booksResponses.Books[0].ISBN13, "9967615110931")
	}
}
func TestUpdateBook(t *testing.T) {
	// Create a test book
	book := Model.Book{
		ID:              1,
		Title:           "Test Book",
		ISBN13:          "9781891830853",
		ISBN10:          "1891830856",
		ListPrice:       9.99,
		PublicationYear: 2021,
		ImageURL:        "https://example.com/test.jpg",
		Edition:         "1st",
		Authors: []Model.Author{
			{
				ID:        1,
				FirstName: "Joel",
				LastName:  "Hartse",
			},
		},
		Publisher: []Model.Publisher{
			{
				ID:   1,
				Name: "Test",
			},
		},
	}
	Model.DB.Books = []Model.Book{book}
	Model.DB.BookAuthors = []Model.BookAuthor{
		{
			AuthorID: 1,
			BookID:   1,
		},
	}
	Model.DB.BookPublishers = []Model.BookPublisher{
		{
			PublisherID: 1,
			BookID:      1,
		},
	}
	// Marshal the book to JSON
	jsonBook, err := json.Marshal(book)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/books/{isbn13}", handler.UpdateBook).Methods("PUT")
	// Create a request with the JSON book data
	req, err := http.NewRequest("PUT", "/books/9781891830853", bytes.NewBuffer(jsonBook))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected status code: got %v, expected %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Book updated successfully"
	if rr.Body.String() != expected {
		t.Errorf("Unexpected response body: got %v, expected %v", rr.Body.String(), expected)
	}

	// Check that the book was updated in the database
	if Model.DB.Books[0].Title != book.Title {
		t.Errorf("Unexpected book title: got %v, expected %v", Model.DB.Books[0].Title, book.Title)
	}
	if Model.DB.Books[0].ISBN10 != book.ISBN10 {
		t.Errorf("Unexpected ISBN10: got %v, expected %v", Model.DB.Books[0].ISBN10, book.ISBN10)
	}
	if Model.DB.Books[0].ListPrice != book.ListPrice {
		t.Errorf("Unexpected list price: got %v, expected %v", Model.DB.Books[0].ListPrice, book.ListPrice)
	}
	if Model.DB.Books[0].PublicationYear != book.PublicationYear {
		t.Errorf("Unexpected publication year: got %v, expected %v", Model.DB.Books[0].PublicationYear, book.PublicationYear)
	}
	if Model.DB.Books[0].ImageURL != book.ImageURL {
		t.Errorf("Unexpected image URL: got %v, expected %v", Model.DB.Books[0].ImageURL, book.ImageURL)
	}
	if Model.DB.Books[0].Edition != book.Edition {
		t.Errorf("Unexpected edition: got %v, expected %v", Model.DB.Books[0].Edition, book.Edition)
	}
	if len(Model.DB.BookAuthors) != len(book.Authors) {
		t.Errorf("Unexpected number of book authors: got %v, expected %v", len(Model.DB.BookAuthors), len(book.Authors))
	}
	if len(Model.DB.BookPublishers) != len(book.Publisher) {
		t.Errorf("Unexpected number of book publishers: got %v, expected %v", len(Model.DB.BookPublishers), len(book.Publisher))
	}
}
func TestDeleteBook(t *testing.T) {
	// Create a test book
	book := Model.Book{
		Title:           "Test Book",
		ISBN13:          "9781891830853",
		ISBN10:          "1891830856",
		ListPrice:       9.99,
		PublicationYear: 2021,
		ImageURL:        "https://example.com/test.jpg",
		Edition:         "1st",
		Authors: []Model.Author{
			{
				ID:        1,
				FirstName: "Joel",
				LastName:  "Hartse",
			},
		},
		Publisher: []Model.Publisher{
			{
				ID:   1,
				Name: "Test",
			},
		},
	}

	// Add the test book to the database
	Model.DB.Books = append(Model.DB.Books, book)

	// Create a request to delete the book
	r := mux.NewRouter()
	r.HandleFunc("/books/{isbn13}", handler.DeleteBook).Methods("DELETE")
	req, err := http.NewRequest("DELETE", "/books/9781891830853", nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected status code: got %v, expected %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Book deleted successfully"
	if rr.Body.String() != expected {
		t.Errorf("Unexpected response body: got %v, expected %v", rr.Body.String(), expected)
	}

	//// Check that the book was actually deleted from the database
	//for _, b := range Model.DB.Books {
	//	if b.ISBN13 == "9781891830853" {
	//		t.Errorf("Book was not deleted from the database")
	//	}
	//}
}

func TestDeleteNonexistentBook(t *testing.T) {
	// Create a request to delete a nonexistent book
	r := mux.NewRouter()
	r.HandleFunc("/books/{isbn13}", handler.DeleteBook).Methods("DELETE")
	req, err := http.NewRequest("DELETE", "/books/9781891830851", nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Unexpected status code: got %v, expected %v", status, http.StatusNotFound)
	}

}
