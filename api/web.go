package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/RomanDerkach/homework/storage"
	"github.com/satori/go.uuid"
)

// Server is function that starts the main process.
func Server() {
	//run server
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/books/", booksHandlerByID)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var newBook storage.Book
		err := json.NewDecoder(r.Body).Decode(&newBook)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if newBook.Title == "" {
			log.Println("There is no title request")
			http.Error(w, "There is not title for the book in the requst",
				http.StatusBadRequest)
			return
		}
		if len(newBook.Ganres) == 0 {
			log.Println("There is no ganres in request")
			http.Error(w, "There is not ganres for the book in the requst",
				http.StatusBadRequest)
			return
		}
		if newBook.Pages == 0 {
			log.Println("There is no pages in reques")
			http.Error(w, "There is not pages for the book in the requst",
				http.StatusBadRequest)
			return
		}
		if newBook.Price == 0 {
			log.Println("There is no price in request")
			http.Error(w, "There is not price for the book in the requst",
				http.StatusBadRequest)
			return
		}
		newBook.ID = uuid.NewV4().String()
		fmt.Println(newBook)
		storage.SaveNewBook(newBook)
	} else {
		jsonBody := storage.GetDBData()
		w.Write(jsonBody)
	}
}

func booksHandlerByID(w http.ResponseWriter, r *http.Request) {
	books := storage.GetBooksData()
	result := storage.Book{}
	if r.Method == "GET" {
		for _, book := range books {
			if book.ID == path.Base(r.URL.Path) {
				result = book
				break
			}
		}
		jsonBody, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "Error converting results to json",
				http.StatusInternalServerError)
			return
		}
		w.Write(jsonBody)
	}
	if r.Method == "DELETE" {
		notfound := true
		for i, book := range books {
			if book.ID == path.Base(r.URL.Path) {
				books = append(books[:i], books[i+1:]...)
				storage.SaveBookData(books)
				w.WriteHeader(http.StatusAccepted)
				notfound = false
				break
			}
		}
		if notfound {
			http.Error(w, "There is no book with such id", http.StatusNotFound)
		}
	}
	if r.Method == "PUT" {

	}
}
