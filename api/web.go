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
var (
	ErrNotFound    = errors.New("can't find the book with given ID")
	ErrIDImmutable = errors.New("id can't be changed")
)

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
	_, err := w.Write(jsonBody)
	if err != nil {
		err = errors.Wrap(err, "cant write a response to user")
		log.Println(err)
	}
}

func booksHandlerPOST(w http.ResponseWriter, r *http.Request) {
	newBook := storage.Book{}
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
	switch r.Method {
	case http.MethodGet:
		booksHandlerByIDGET(w, r)
	case http.MethodDelete:
		booksHandlerByIDDELETE(w, r)
	case http.MethodPut:
		booksHandlerByIDPUT(w, r)
	case http.MethodPost:
		booksHandlerByIDPOST(w, r)
	}
}

func booksHandlerByIDGET(w http.ResponseWriter, r *http.Request) {
	books := storage.GetBooksData()
	reqBookID := path.Base(r.URL.Path)
	bookIndex, err := indexByID(reqBookID, books)
	if err != nil {
		err = errors.Wrap(err, "book with such id not found")
		http.Error(w, "there is no book with such id", http.StatusNotFound)
		return
	}
	jsonBody, err := json.Marshal(books[bookIndex])
	if err != nil {
		err = errors.Wrap(err, "error converting results to json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonBody)
	if err != nil {
		err = errors.Wrap(err, "cant write a response to user")
		log.Println(err)
	}
}
func booksHandlerByIDDELETE(w http.ResponseWriter, r *http.Request) {
	books := storage.GetBooksData()
	reqBookID := path.Base(r.URL.Path)
	bookIndex, err := indexByID(reqBookID, books)
	if err != nil {
		err = errors.Wrap(err, "book with such id not found")
		http.Error(w, "there is no book with such id", http.StatusNotFound)
		return
	}
	books = append(books[:bookIndex], books[bookIndex+1:]...)
	// books could be changed from other goroutine while executing
	err = storage.SaveBookData(books)
	if err != nil {
		err := errors.Wrap(err, "can't save updated books")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
func booksHandlerByIDPUT(w http.ResponseWriter, r *http.Request) {
	books := storage.GetBooksData()
	reqBookID := path.Base(r.URL.Path)
	bookIndex, err := indexByID(reqBookID, books)
	if err != nil {
		err = errors.Wrap(err, "book with such id not found")
		http.Error(w, "There is no book with such id", http.StatusNotFound)
		return
	}

	book := books[bookIndex]
	origID := book.ID
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		err = errors.Wrap(err, "currapted request body")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if book.ID != origID {
		err = errors.Wrap(ErrIDImmutable, "book id is immutable from ouside")
		log.Println(err)
		http.Error(w, "ID can't be changed", http.StatusBadRequest)
		return
	}

	books[bookIndex] = book
	storage.SaveBookData(books)
	jsonBody, err := json.Marshal(book)
	if err != nil {
		err = errors.Wrap(err, "data cant be jsnify")
		http.Error(w, "Can't jsonify own data", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonBody)
	if err != nil {
		err = errors.Wrap(err, "cant write a response to user")
		log.Println(err)
	}
}

func booksHandlerByIDPOST(w http.ResponseWriter, r *http.Request) {
	books := storage.GetBooksData()
	resBooks := storage.Books{}
	bookFilter := storage.BookFilter{}
	err := json.NewDecoder(r.Body).Decode(&bookFilter)
	if err != nil {
		err = errors.Wrap(err, "currapted request body in filter")
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
		err = errors.Wrap(err, "cant murshal result")
		log.Println(err)
		http.Error(w, "Can't jsonify own data", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonBody)
	if err != nil {
		err = errors.Wrap(err, "cant write a response to user")
		log.Println(err)
	}
}
