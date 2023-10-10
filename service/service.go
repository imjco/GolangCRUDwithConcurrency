package service

import (
	"companyXyzProject/Model"
	"encoding/csv"
	"fmt"
	"os"
)

// function to remove a publisher from a slice of publishers
func RemoveBookPublishersWithIDs(bookPublishers []Model.BookPublisher, idsToRemove []int) []Model.BookPublisher {
	newBookPublishers := []Model.BookPublisher{}
	for _, bookPublisher := range bookPublishers {
		if !contains(idsToRemove, bookPublisher.PublisherID) {
			newBookPublishers = append(newBookPublishers, bookPublisher)
		}
	}
	return newBookPublishers
}

// function to remove a author from a slice of authors
func RemoveBookPublishersWithIDsAuthors(bookAuthors []Model.BookAuthor, idsToRemove []int) []Model.BookAuthor {
	newBookAuthors := []Model.BookAuthor{}
	for _, bookAuthor := range bookAuthors {
		if !contains(idsToRemove, bookAuthor.AuthorID) {
			newBookAuthors = append(newBookAuthors, bookAuthor)
		}
	}
	return newBookAuthors
}

// function to check if a slice contains a given item
func contains(slice []int, item int) bool {
	for _, ele := range slice {
		if ele == item {
			return true
		}
	}
	return false
}

// Function to append a new ISBN/EAN to a CSV file
func AppendToCSV(fileName string, isbnEan Model.Book) error {
	// Check if the file exists

	_, err := os.Stat(fileName)

	// If the file doesn't exist, create it
	if os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		file.Close()
	}

	// Open the CSV file for appending
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)

	// Write the new ISBN/EAN to the CSV file
	err = writer.Write([]string{isbnEan.Title, isbnEan.ISBN13, isbnEan.ISBN10, fmt.Sprint(isbnEan.ListPrice), fmt.Sprint(isbnEan.PublicationYear), isbnEan.ImageURL, isbnEan.Edition})
	if err != nil {
		return err
	}

	// Flush any buffered data to the file
	writer.Flush()

	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}
