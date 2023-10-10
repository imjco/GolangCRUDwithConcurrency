package handler

import (
	"companyXyzProject/Model"
	"companyXyzProject/service"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"
)

// Validate checks if a Book struct meets the validation criteria
func ValidateBookAdd(b Model.Book) error {

	// Check if the book already exists using isbn13
	if len(Model.DB.Books) != 0 {
		for _, book := range Model.DB.Books {
			if book.ISBN13 == b.ISBN13 {
				return fmt.Errorf("Book already exists")
			}
		}
	}

	// Check if title is greater than 3 characters
	if utf8.RuneCountInString(b.Title) < 3 {
		return fmt.Errorf("Title Too Short")
	}
	// Check if ISBN13 is valid
	if !service.Validate(b.ISBN13) {
		return fmt.Errorf("Invalid ISBN13")
	}
	// Check if ISBN10 is valid
	if !service.Validate(b.ISBN10) {
		return fmt.Errorf("Invalid ISBN10")
	}
	// Check if price is greater than 0 and not null
	if b.ListPrice <= 0 {
		return fmt.Errorf("Invalid Price")
	}
	// Check if year is valid and not null and must have 4 digits

	if b.PublicationYear >= 1900 && b.PublicationYear <= time.Now().Year() && b.PublicationYear <= time.Now().Year() {
	} else {
		fmt.Println(b.PublicationYear)
	}

	//// Check if the author exists
	authfound := false
	if len(Model.DB.Authors) != 0 {
		for _, author := range Model.DB.Authors {
			for _, auth := range b.Authors {
				if auth.ID == author.ID {
					authfound = true
					break
				}
			}
		}
	}

	// Check if the publisher exists
	pubfound := false
	if len(Model.DB.Publishers) != 0 {
		for _, publisher := range Model.DB.Publishers {
			for _, pub := range b.Publisher {
				if pub.ID == publisher.ID {
					pubfound = true
					break
				}
			}
		}
	}
	if pubfound == false {
		return fmt.Errorf("Publisher does not exists")
	}
	if authfound == false {
		return fmt.Errorf("Author does not exists")
	}
	return nil
}

func ValidateAuthor(a Model.Author) error {
	if len(Model.DB.Authors) != 0 {
		for _, author := range Model.DB.Authors {
			if author.FirstName == a.FirstName && author.LastName == a.LastName {
				return fmt.Errorf("Author already exists")
			}
		}
	}

	if utf8.RuneCountInString(a.FirstName) < 1 {
		return fmt.Errorf("First Name Too Short")
	}

	if utf8.RuneCountInString(a.LastName) < 1 {
		return fmt.Errorf("Last Name Too Short")
	}

	return nil
}

func ValidatePublisher(p Model.Publisher) error {
	if len(Model.DB.Publishers) != 0 {
		for _, publisher := range Model.DB.Publishers {
			if publisher.Name == p.Name {
				return fmt.Errorf("Publisher already exists")
			}
		}
	}

	if utf8.RuneCountInString(p.Name) < 3 {
		return fmt.Errorf("Publisher Name Too Short")
	}

	return nil
}

func ValidatePageLimit(page string, limit string, structlen int) (int, int, error) {

	// check if structlen is greater than 0
	if structlen <= 0 {
		return 0, 0, fmt.Errorf("Struct Length must be greater than 0")
	}

	// Convert query parameters to integers
	if page == "" || limit == "" {
		return 0, 0, fmt.Errorf("No Page or Limit Found")
	}
	// Convert query parameters to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return 0, 0, fmt.Errorf("No Page or Limit Found")
	}
	// Convert query parameters to integers

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return 0, 0, fmt.Errorf("No Page or Limit Found")
	}

	// Check if page and limit are greater than 0
	if pageInt < 0 || limitInt < 0 {
		return 0, 0, fmt.Errorf("Page and Limit must be greater than 0")
	}
	// Check if page and limit are greater than 0
	if limitInt > 100 {
		return 0, 0, fmt.Errorf("Limit must be less than 100")
	}
	if limitInt > structlen {
		limitInt = structlen
	}

	totalPages := structlen / limitInt
	if structlen%limitInt != 0 {
		totalPages++
	}

	// Check if the requested page is out of range
	if pageInt < 1 || pageInt > totalPages {
		// Return an appropriate response for an out-of-range page
		return 0, 0, fmt.Errorf("Page is out of range")
	}

	// Calculate the start and end indices based on page and limit
	startIndex := (pageInt - 1) * limitInt
	endIndex := startIndex + limitInt

	// Ensure startIndex is within bounds
	if startIndex < 0 {
		startIndex = 0
	}

	return startIndex, endIndex, nil

}
