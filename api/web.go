package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	// "github.com/gorilla/mux"
	"github.com/RomanDerkach/homework/storage"
	"github.com/satori/go.uuid"
)

func Server() {
	//run server
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/books/", booksHandlerByID)
	http.ListenAndServe(":8081", nil)
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var newBook storage.Book
		err := decoder.Decode(&newBook)
		if err != nil {
			log.Println(err)
		}
		if newBook.Title == "" {
			log.Println("There is no title")
		}
		if len(newBook.Ganres) == 0 {
			log.Println("There is no ganres")
		}
		if newBook.Pages == 0 {
			log.Println("There is no pages")
		}
		if newBook.Price == 0 {
			log.Println("There is no price")
		}
		newBook.ID = uuid.NewV4().String()
		fmt.Println(newBook)
		storage.SaveBookData(newBook)
	} else {
		jsonBody := storage.GetDBData()
		w.Write(jsonBody)
	}
}

func booksHandlerByID(w http.ResponseWriter, r *http.Request) {
	books := storage.GetBooksData()
	result := storage.Book{}
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
