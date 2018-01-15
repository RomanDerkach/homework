package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

//JSONStorage this a struct for describing json storage
type JSONStorage struct {
	storage      *os.File
	storageMutex sync.RWMutex
}

//NewJSONStorage create a new object of JSON storage
func NewJSONStorage(pathToStorage string) (*JSONStorage, error) {
	storagePath, err := filepath.Abs(pathToStorage)
	if err != nil {
		return nil, err
	}

	if _, err = os.Stat(storagePath); err != nil {
		// It would be better to know if something wrong before application start
		return nil, err
	}
	file, err := os.OpenFile(storagePath, os.O_RDWR|os.O_CREATE, 0660)
	//We handle this closing in main ???
	// defer file.Close()
	if err != nil {
		return nil, err
	}
	return &JSONStorage{storage: file}, nil
}

func indexByID(id string, books Books) (int, error) {
	for index, book := range books {
		if book.ID == id {
			return index, nil
		}
	}
	return 0, ErrNotFound
}

//Close opened file.
func (j *JSONStorage) Close() error {
	return nil
}

func (j *JSONStorage) writeToFile(books Books) error {
	raw, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		return errors.Wrap(err, "can't marshal json")
	}
	err = j.storage.Truncate(0)
	if err != nil {
		return err
	}
	_, err = j.storage.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = j.storage.Write(raw)
	return errors.Wrap(err, "can't write to file")
}

//GetBooks return all books as a list of book struct
func (j *JSONStorage) GetBooks() (Books, error) {
	j.storageMutex.RLock()
	defer j.storageMutex.RUnlock()
	return j.getBooks()
}

func (j *JSONStorage) getBooks() (Books, error) {
	var books Books

	_, err := j.storage.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(j.storage)
	if err != nil {
		err = errors.Wrap(err, "cant get books")
		log.Println(err)
		return nil, err
	}
	err = json.Unmarshal(raw, &books)
	if err != nil {
		err = errors.Wrap(err, "cant unmarshal books")
		log.Println(err)
		return nil, err
	}
	return books, err
}

//SaveBook save only one new book to storage
func (j *JSONStorage) SaveBook(book Book) error {
	j.storageMutex.Lock()
	defer j.storageMutex.Unlock()
	books, err := j.getBooks()
	if err != nil {
		err = errors.Wrap(err, "can't get book data")
		log.Println(err)
		return err
	}
	books = append(books, book)
	return errors.Wrap(j.writeToFile(books), "can't save book data")
}

//GetBook ...
func (j *JSONStorage) GetBook(bookID string) (Book, error) {
	book := Book{}
	books, err := j.getBooks()
	if err != nil {
		err = errors.Wrap(err, "cant get books from storage")
		log.Println(err)
		return book, err
	}
	bookIndex, err := indexByID(bookID, books)
	if err != nil {
		err = errors.Wrap(err, "book with such id not found")
		log.Println(err)
		return book, err
	}
	book = books[bookIndex]
	return book, nil
}

//DeleteBook ...
func (j *JSONStorage) DeleteBook(bookID string) error {
	j.storageMutex.Lock()
	defer j.storageMutex.Unlock()
	books, err := j.getBooks()
	if err != nil {
		err = errors.Wrap(err, "cant get books from storage")
		log.Println(err)
		return err
	}
	bookIndex, err := indexByID(bookID, books)
	if err != nil {
		err = errors.Wrap(err, "book with such id not found")
		log.Println(err)
		return err
	}
	fmt.Println(bookIndex)
	books = append(books[:bookIndex], books[bookIndex+1:]...)
	err = j.writeToFile(books)
	if err != nil {
		err = errors.Wrap(err, "can't update storage")
		log.Println(err)
		return err
	}
	return nil
}

//UpdateBook ...
func (j *JSONStorage) UpdateBook(bookID string, updBook Book) (Book, error) {
	j.storageMutex.Lock()
	defer j.storageMutex.Unlock()
	books, err := j.getBooks()
	if err != nil {
		err = errors.Wrap(err, "cant get books from storage")
		log.Println(err)
		return Book{}, err
	}
	bookIndex, err := indexByID(bookID, books)
	if err != nil {
		err = errors.Wrap(err, "cant get book with such id from storage")
		log.Println(err)
		return Book{}, err
	}
	fmt.Println(updBook)

	book := &books[bookIndex]
	book.Price = updBook.Price
	book.Title = updBook.Title
	book.Pages = updBook.Pages
	book.Genres = updBook.Genres

	return *book, nil
}

//FilterBooks filter books that is needed
func (j *JSONStorage) FilterBooks(filter BookFilter) (Books, error) {
	resBooks := Books{}
	books, err := j.getBooks()
	if err != nil {
		err = errors.Wrap(err, "cant get books from storage")
		log.Println(err)
		return Books{}, err
	}

	for _, book := range books {
		if strings.Contains(strings.ToLower(book.Title), strings.ToLower(filter.Title)) {
			resBooks = append(resBooks, book)
		}
	}
	return resBooks, nil
}
