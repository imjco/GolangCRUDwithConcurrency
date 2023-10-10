package main

import (
	"companyXyzProject/handler"
	"companyXyzProject/seeder"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var mu sync.Mutex

func main() {

	// seed data to be used in the application
	//seeder.SeedData()
	seeder.GenerateRandomData()

	// Create a new router using Gorilla Mux
	r := mux.NewRouter()

	// Define routes for authors
	// Sample route localhost:8080/authors?page=1&limit=20
	r.HandleFunc("/authors", handler.GetAuthor).Methods("GET")
	r.HandleFunc("/authors", handler.AddAuthor).Methods("POST")
	r.HandleFunc("/authors/{id}", handler.GetAuthorByID).Methods("GET")
	r.HandleFunc("/authors/{id}", handler.UpdateAuthor).Methods("PUT")
	r.HandleFunc("/authors/{id}", handler.DeleteAuthor).Methods("DELETE")
	// Define routes for publishers
	// Sample route localhost:8080/publishers?page=1&limit=20
	r.HandleFunc("/publishers", handler.GetPublisher).Methods("GET")
	r.HandleFunc("/publishers", handler.AddPublisher).Methods("POST")
	r.HandleFunc("/publishers/{id}", handler.GetPublisherByID).Methods("GET")
	r.HandleFunc("/publishers/{id}", handler.UpdatePublisher).Methods("PUT")
	r.HandleFunc("/publishers/{id}", handler.DeletePublisher).Methods("DELETE")
	// Define routes for books
	// Sample route localhost:8080/books?page=1&limit=20
	r.HandleFunc("/books", handler.GetBooks).Methods("GET")
	r.HandleFunc("/books", handler.AddBook).Methods("POST")
	r.HandleFunc("/books/{isbn13}", handler.GetBookByISBN13).Methods("GET")
	r.HandleFunc("/books/{isbn13}", handler.UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{isbn13}", handler.DeleteBook).Methods("DELETE")

	r.HandleFunc("/convertisbn13to10/{isbn13}", handler.GetISBN13TO10).Methods("GET")
	r.HandleFunc("/convertisbn10to13/{isbn10}", handler.GetISBN10TO13).Methods("GET")

	// Start the HTTP server
	fmt.Println("Server started on :8080")
	http.Handle("/", r)

	go http.ListenAndServe(":8080", handler.CorsMiddleware(r))
	// Block to keep the server running
	select {}
}
