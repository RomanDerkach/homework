package api

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/RomanDerkach/homework/storage"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

// ErrNotFound for the error when we havn't found something
var ErrNotFound = errors.New("can't find the book with given ID")

// Server is function that starts the main process.
func Server() error {
	//run server
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/books/", booksHandlerByID)
	return http.ListenAndServe(":8081", nil)
}

func indexByID(id string, books storage.Books) (int, error) {
	for index, book := range books {
		if book.ID == id {
			return index, nil
		}
	}
	return 0, ErrNotFound
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		booksHandlerPOST(w, r)
	case http.MethodGet:
		booksHandlerGET(w, r)
	}
}

func booksHandlerGET(w http.ResponseWriter, r *http.Request) {
	jsonBody := storage.GetDBData()
	w.Write(jsonBody)
}

func booksHandlerPOST(w http.ResponseWriter, r *http.Request) {
	var newBook storage.Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		err := errors.Wrap(err, "corrupted request body")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := newBook.Validate(); err != nil {
		err := errors.Wrap(err, "not valid book")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newBook.ID = uuid.NewV4().String()
	//fmt.Println(newBook)
	err = storage.SaveNewBook(newBook)
	if err != nil {
		err := errors.Wrap(err, "can't save book")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func booksHandlerByID(w http.ResponseWriter, r *http.Request) {
	books := storage.GetBooksData()
	reqBookID := path.Base(r.URL.Path)
	resBook := storage.Book{}
	resBooks := storage.Books{}

	if r.Method == "GET" {
		// use indexByID(reqBookID)
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
		_, err = w.Write(jsonBody)
		if err != nil {
			panic(err)
		}
	}

	if r.Method == "DELETE" {
		bookIndex, err := indexByID(reqBookID, books)
		if err != nil {
			if err == ErrNotFound {
				http.Error(w, "There is no book with such id", http.StatusNotFound)
				return
			}
			http.Error(w, "There is no book with such id", http.StatusNotFound)
			return
		}
		books = append(books[:bookIndex], books[bookIndex+1:]...)
		// books could be changed from other goroutine while executing
		storage.SaveBookData(books)
		w.WriteHeader(http.StatusAccepted)
	}

	if r.Method == "PUT" {
		bookIndex, err := indexByID(reqBookID, books)
		if err != nil {
			http.Error(w, "There is no book with such id", http.StatusNotFound)
			return
		}

		book := books[bookIndex]
		origID := book.ID
		err = json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if book.ID != origID {
			// error ID can't be changed
		}

		books[bookIndex] = book
		storage.SaveBookData(books)
		jsonBody, err := json.Marshal(book)
		if err != nil {
			http.Error(w, "Can't jsonify own data", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(jsonBody)
	}

	if r.Method == "POST" {
		bookFilter := storage.BookFilter{}
		err := json.NewDecoder(r.Body).Decode(&bookFilter)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if bookFilter == (storage.BookFilter{}) {
			http.Error(w, "Currently we support only {'title': 'Tom'}", http.StatusBadRequest)
			return
		}

		//fmt.Println(bookFilter)
		for _, book := range books {
			if strings.Contains(strings.ToLower(book.Title), strings.ToLower(bookFilter.Title)) {
				resBooks = append(resBooks, book)
			}
		}

		jsonBody, err := json.Marshal(resBooks)
		if err != nil {
			http.Error(w, "Can't jsonify own data", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(jsonBody)
		if err != nil {
			panic(err)
		}
	}
}
