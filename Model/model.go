package Model

import (
	"sync"
)

type Author struct {
	ID         int
	FirstName  string
	LastName   string
	MiddleName string
	Books      []Book // Slice to represent many books by the author
}

type Book struct {
	ID              int
	Title           string
	ISBN13          string `maxLen:"13"`
	ISBN10          string
	ListPrice       float64
	PublicationYear int
	Publisher       []Publisher
	ImageURL        string
	Edition         string
	Authors         []Author // Slice to represent many authors of the book
}

type BooksResponse struct {
	Books           []Book
	BooksTotalCount int
}

type Publisher struct {
	ID    int
	Name  string
	Books []Book // Slice to represent many books by the publisher
}

type BookAuthor struct {
	BookID   int
	AuthorID int
}

type BookPublisher struct {
	PublisherID int
	BookID      int
}

type Database struct {
	Books          []Book
	Authors        []Author
	Publishers     []Publisher
	BookAuthors    []BookAuthor
	BookPublishers []BookPublisher
}

// declare a variable of type Database
var (
	DB          Database
	MU          sync.Mutex
	CsvFileName string
)
