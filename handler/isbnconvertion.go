package handler

import (
	"companyXyzProject/Model"
	"companyXyzProject/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func GetISBN13TO10(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	vars := mux.Vars(r)
	isbn13 := vars["isbn13"]
	resultChan := make(chan string)

	valresult := service.Validate13(isbn13)
	if valresult == false {
		http.Error(w, "Invalid ISBN13", http.StatusBadRequest)
		return
	}

	go func() {
		Model.MU.Lock()
		defer Model.MU.Unlock()

		conIsbn10, err := service.To10(isbn13)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resultChan <- conIsbn10
	}()

	// Wait for the goroutine to complete
	result := <-resultChan

	type response struct {
		ISBN10 string `json:"isbn10"`
	}
	Res := response{
		ISBN10: result,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Res)
}
func GetISBN10TO13(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	vars := mux.Vars(r)
	isbn10 := vars["isbn10"]
	resultChan := make(chan string)

	valresult := service.Validate10(isbn10)
	if valresult == false {

		http.Error(w, "Invalid ISBN10", http.StatusBadRequest)
		return
	}

	go func() {
		Model.MU.Lock()
		defer Model.MU.Unlock()

		conIsbn10, err := service.To13(isbn10)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resultChan <- conIsbn10
	}()

	// Wait for the goroutine to complete
	result := <-resultChan

	type response struct {
		ISBN13 string `json:"isbn13"`
	}
	Res := response{
		ISBN13: result,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Res)

}
