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

// ErrNotFound for the error when we haven't found something
var (
	ErrNotFound    = errors.New("can't find the book with given ID")
	ErrIDImmutable = errors.New("id can't be changed")
)

//Storage describe a storage interface for handler
type Storage interface {
	GetDBData() ([]byte, error)
	GetBooksData() (storage.Books, error)
	SaveNewBook(book storage.Book) error
	SaveBookData(books storage.Books) error
}

type handler struct {
	storage Storage
}

func NewHandler(storage Storage) *handler {
	return &handler{
		storage: storage,
	}
}

// Server is function that starts the main process.
func Server(handler *handler) error {
	//run server
	http.HandleFunc("/books", handler.booksHandler)
	http.HandleFunc("/books/", handler.booksHandlerByID)
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

func (h *handler) booksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.booksHandlerPOST(w, r)
	case http.MethodGet:
		h.booksHandlerGET(w, r)
	}
}

func (h *handler) booksHandlerGET(w http.ResponseWriter, r *http.Request) {
	jsonBody, err := h.storage.GetDBData()
	if err != nil {
		err = errors.Wrap(err, "cant get data from storage")
		log.Println(err)
		http.Error(w, "Cant get data from storage", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonBody)
	if err != nil {
		err = errors.Wrap(err, "cant write a response to user")
		log.Println(err)
	}
}

func (h *handler) booksHandlerPOST(w http.ResponseWriter, r *http.Request) {
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
	err = h.storage.SaveNewBook(newBook)
	if err != nil {
		err := errors.Wrap(err, "can't save book")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) booksHandlerByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.booksHandlerByIDGET(w, r)
	case http.MethodDelete:
		h.booksHandlerByIDDELETE(w, r)
	case http.MethodPut:
		h.booksHandlerByIDPUT(w, r)
	case http.MethodPost:
		h.booksHandlerByIDPOST(w, r)
	}
}

func (h *handler) booksHandlerByIDGET(w http.ResponseWriter, r *http.Request) {
	books, err := h.storage.GetBooksData()
	if err != nil {
		err = errors.Wrap(err, "cant get data from storage")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reqBookID := path.Base(r.URL.Path)
	bookIndex, err := indexByID(reqBookID, books)
	if err != nil {
		err = errors.Wrap(err, "book with such id not found")
		log.Println(err)
		http.Error(w, "there is no book with such id", http.StatusNotFound)
		return
	}
	jsonBody, err := json.Marshal(books[bookIndex])
	if err != nil {
		err = errors.Wrap(err, "error converting results to json")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonBody)
	if err != nil {
		err = errors.Wrap(err, "cant write a response to user")
		log.Println(err)
	}
}
func (h *handler) booksHandlerByIDDELETE(w http.ResponseWriter, r *http.Request) {
	books, err := h.storage.GetBooksData()
	if err != nil {
		err = errors.Wrap(err, "cant get data from storage")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reqBookID := path.Base(r.URL.Path)
	bookIndex, err := indexByID(reqBookID, books)
	if err != nil {
		err = errors.Wrap(err, "book with such id not found")
		log.Println(err)
		http.Error(w, "there is no book with such id", http.StatusNotFound)
		return
	}
	books = append(books[:bookIndex], books[bookIndex+1:]...)
	// books could be changed from other goroutine while executing
	err = h.storage.SaveBookData(books)
	if err != nil {
		err := errors.Wrap(err, "can't save updated books")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
func (h *handler) booksHandlerByIDPUT(w http.ResponseWriter, r *http.Request) {
	books, err := h.storage.GetBooksData()
	if err != nil {
		err = errors.Wrap(err, "cant get data from storage")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reqBookID := path.Base(r.URL.Path)
	bookIndex, err := indexByID(reqBookID, books)
	if err != nil {
		err = errors.Wrap(err, "book with such id not found")
		log.Println(err)
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
	err = h.storage.SaveBookData(books)
	if err != nil {
		errors.Wrap(err, "cant save updated books")
		log.Println(err)
		http.Error(w, "Cant save updated data", http.StatusInternalServerError)
	}
	jsonBody, err := json.Marshal(book)
	if err != nil {
		err = errors.Wrap(err, "data cant be jsnify")
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

func (h *handler) booksHandlerByIDPOST(w http.ResponseWriter, r *http.Request) {
	books, err := h.storage.GetBooksData()
	if err != nil {
		err = errors.Wrap(err, "cant get data from storage")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resBooks := storage.Books{}
	bookFilter := storage.BookFilter{}
	err = json.NewDecoder(r.Body).Decode(&bookFilter)
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
