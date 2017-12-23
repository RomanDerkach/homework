package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/RomanDerkach/homework/storage"
	"github.com/satori/go.uuid"
)

// ErrNotFound for the error when we havn't found something
var ErrNotFound = errors.New("can't find the book with given ID")

// Server is function that starts the main process.
func Server() {
	//run server
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/books/", booksHandlerByID)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func indexByID(id string, books []storage.Book) (int, error) {
	for index, book := range books {
		if book.ID == id {
			return index, nil
		}
	}
	return 0, ErrNotFound
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
	reqBookID := path.Base(r.URL.Path)
	resBook := storage.Book{}
	if r.Method == "GET" {
		for _, book := range books {
			if book.ID == reqBookID {
				resBook = book
				break
			}
		}
		jsonBody, err := json.Marshal(resBook)
		if err != nil {
			http.Error(w, "Error converting results to json",
				http.StatusInternalServerError)
			return
		}
		w.Write(jsonBody)
	}
	if r.Method == "DELETE" {
		bookIndex, err := indexByID(reqBookID, books)
		if err != nil {
			http.Error(w, "There is no book with such id", http.StatusNotFound)
			return
		}
		books = append(books[:bookIndex], books[bookIndex+1:]...)
		storage.SaveBookData(books)
		w.WriteHeader(http.StatusAccepted)

	}
	if r.Method == "PUT" {
		bookIndex, err := indexByID(reqBookID, books)
		if err != nil {
			http.Error(w, "There is no book with such id", http.StatusNotFound)
			return
		}
		err = json.NewDecoder(r.Body).Decode(&resBook)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		book := books[bookIndex]
		if resBook.Title != "" {
			book.Title = resBook.Title
		}
		if len(resBook.Ganres) != 0 {
			book.Ganres = resBook.Ganres
		}
		if resBook.Pages != 0 {
			book.Pages = resBook.Pages
		}
		if resBook.Price != 0 {
			book.Price = resBook.Price
		}
		books[bookIndex] = book
		storage.SaveBookData(books)
		jsonBody, err := json.Marshal(book)
		if err != nil {
			http.Error(w, "Can't jsonify own data", http.StatusInternalServerError)
			return
		}
		w.Write(jsonBody)
	}
}
