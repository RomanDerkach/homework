package api

import (
	"encoding/json"
	"log"
	"net/http"
	"path"

	"github.com/RomanDerkach/homework/storage"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var (
	ErrIDImmutable = errors.New("id can't be changed")
)

//Storage describe a storage interface for handler
type Storage interface {
	GetBooks() (storage.Books, error)
	GetBook(bookID string) (storage.Book, error)
	SaveBook(book storage.Book) error
	DeleteBook(bookID string) error
	UpdateBook(bookID string, updBook storage.Book) (storage.Book, error)
	FilterBooks(filter storage.BookFilter) (storage.Books, error)
}

type handler struct {
	storage Storage
}

//NewHandler constructor for handlers
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

func (h *handler) booksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.booksHandlerPOST(w, r)
	case http.MethodGet:
		h.booksHandlerGET(w, r)
	}
}

func (h *handler) booksHandlerGET(w http.ResponseWriter, r *http.Request) {
	books, err := h.storage.GetBooks()
	if err != nil {
		err = errors.Wrap(err, "cant get data from storage")
		log.Println(err)
		http.Error(w, "Cant get data from storage", http.StatusInternalServerError)
		return
	}
	jsonBody, err := json.Marshal(books)
	if err != nil {
		err = errors.Wrap(err, "Cant marshal books")
		log.Println(err)
		http.Error(w, "Can't form a response", http.StatusInternalServerError)
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
	err = h.storage.SaveBook(newBook)
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
	reqBookID := path.Base(r.URL.Path)
	book, err := h.storage.GetBook(reqBookID)
	if err != nil {
		if err == storage.ErrNotFound {
			err = errors.Wrap(err, "there is no book with such id in storage")
			log.Println(err)
			http.Error(w, "There is no book with such id", http.StatusNotFound)
			return
		}
		err = errors.Wrap(err, "cant get data from storage")
		log.Println(err)
		http.Error(w, "Can't get data from the storage", http.StatusInternalServerError)
		return
	}
	jsonBody, err := json.Marshal(book)
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
	reqBookID := path.Base(r.URL.Path)
	err := h.storage.DeleteBook(reqBookID)
	if err != nil {
		if err == storage.ErrNotFound {
			err = errors.Wrap(err, "there is no book with such id in storage")
			log.Println(err)
			http.Error(w, "There is no book with such id", http.StatusNotFound)
			return
		}
		err = errors.Wrap(err, "can't save updated books")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

//follow here
func (h *handler) booksHandlerByIDPUT(w http.ResponseWriter, r *http.Request) {

	reqBookID := path.Base(r.URL.Path)
	book := storage.Book{}

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		err = errors.Wrap(err, "currapted request body")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err = h.storage.UpdateBook(reqBookID, book)
	if err != nil {
		if err == storage.ErrNotFound {
			err = errors.Wrap(err, "there is no book with such id in storage")
			log.Println(err)
			http.Error(w, "There is no book with such id", http.StatusNotFound)
			return
		}
		err = errors.Wrap(err, "can't update a book")
		log.Println(err)
		http.Error(w, "Can't upgrade a book", http.StatusInternalServerError)
	}

	jsonBody, err := json.Marshal(book)
	if err != nil {
		err = errors.Wrap(err, "data cant be jsonify")
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

	books, err := h.storage.FilterBooks(bookFilter)
	if err != nil {
		err = errors.Wrap(err, "cant filter books")
		log.Println(err)
		http.Error(w, "Cant filter books", http.StatusInternalServerError)
	}

	jsonBody, err := json.Marshal(books)
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
