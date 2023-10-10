package service_test

import (
	"companyXyzProject/Model"
	"companyXyzProject/service"
	"encoding/csv"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAppendToCSV(t *testing.T) {
	// Create a temporary file for testing
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "test.csv")

	// Create a test book
	book := Model.Book{
		Title:           "Test Book",
		ISBN13:          "9780596000480",
		ISBN10:          "0596000480",
		ListPrice:       9.99,
		PublicationYear: 2021,
		ImageURL:        "https://example.com/test.jpg",
		Edition:         "1st",
	}

	// Append the book to the CSV file
	err = service.AppendToCSV(tempFile, book)
	if err != nil {
		t.Fatalf("Failed to append book to CSV file: %v", err)
	}

	// Read the CSV file and check its contents
	file, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("Expected 1 record, but got %d", len(records))
	}

	expected := []string{"Test Book", "9780596000480", "0596000480", "9.99", "2021", "https://example.com/test.jpg", "1st"}
	if len(records[0]) != len(expected) {
		t.Fatalf("Expected %d fields, but got %d", len(expected), len(records[0]))
	}

	for i, field := range records[0] {
		if field != expected[i] {
			t.Errorf("Expected field %d to be %q, but got %q", i, expected[i], field)
		}
	}
}
