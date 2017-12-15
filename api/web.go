package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	// "github.com/gorilla/mux"
	"github.com/RomanDerkach/homework/storage"
)

type Book struct {
	ID     string
	Title  string
	Ganres []string
	Pages  int
	Price  float64
}

func Server() {
	http.HandleFunc("/books", BooksHandler)
	http.HandleFunc("/books/", BooksHandlerByID)
	http.ListenAndServe(":8081", nil)
}

func BooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("AAA")
	} else {
		jsonBody := storage.GetBooksData()
		w.Write(jsonBody)
	}
}

func BooksHandlerByID(w http.ResponseWriter, r *http.Request) {
	data := storage.GetBooksData()
	books := []Book{}
	result := Book{}
	json.Unmarshal(data, &books)
	fmt.Println(books[0].ID)
	for _, book := range books {
		if book.ID == path.Base(r.URL.Path) {
			result = book
		}
	}
	jsonBody, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Error converting results to json",
			http.StatusInternalServerError)
	}
	w.Write(jsonBody)
}
