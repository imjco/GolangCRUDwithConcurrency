package handler_test

import (
	"companyXyzProject/Model"
	"companyXyzProject/handler"
	"fmt"
	"testing"
)

func TestValidateBookAdd(t *testing.T) {

	publisher := Model.Publisher{
		ID:   len(Model.DB.Publishers) + 1,
		Name: "Test",
	}
	Model.DB.Publishers = append(Model.DB.Publishers, publisher)

	author := Model.Author{
		ID:         1,
		FirstName:  "Joel",
		LastName:   "Hartse",
		MiddleName: "",
	}
	Model.DB.Authors = append(Model.DB.Authors, author)

	// Create a test book
	book := Model.Book{
		Title:           "Test Book",
		ISBN13:          "9967615110931",
		ISBN10:          "7615110939",
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

	// Validate the book
	err := handler.ValidateBookAdd(book)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestValidateAuthor(t *testing.T) {
	// Create a test author
	author := Model.Author{
		ID:        len(Model.DB.Authors) + 1,
		FirstName: "John",
		LastName:  "Doe",
	}

	// Validate the author
	err := handler.ValidateAuthor(author)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestValidatePublisher(t *testing.T) {
	// Create a test publisher
	publisher := Model.Publisher{
		ID:   len(Model.DB.Publishers) + 1,
		Name: "Test Publisher" + fmt.Sprint(len(Model.DB.Publishers)),
	}

	// Validate the publisher
	err := handler.ValidatePublisher(publisher)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestValidatePageLimit(t *testing.T) {
	// Create a test slice of length 10
	slice := make([]int, 10)

	// Validate page and limit
	startIndex, endIndex, err := handler.ValidatePageLimit("1", "5", len(slice))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if startIndex != 0 {
		t.Errorf("Expected startIndex to be 0, but got %d", startIndex)
	}

	if endIndex != 5 {
		t.Errorf("Expected endIndex to be 5, but got %d", endIndex)
	}
}
