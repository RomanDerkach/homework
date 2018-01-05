package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/pkg/errors"
)

//JSONStorage this a struct for describing json storage
type JSONStorage struct {
	storage string
}

//NewJSONStorage create a new object of JSON storage
func NewJSONStorage(pathToStorage string) *JSONStorage {
	return &JSONStorage{
		storage: pathToStorage,
	}
}

func indexByID(id string, books Books) (int, error) {
	for index, book := range books {
		if book.ID == id {
			return index, nil
		}
	}
	return 0, ErrNotFound
}

func (j *JSONStorage) writeToFile(books Books) error {
	raw, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		return errors.Wrap(err, "can't marshal json")
	}
	return errors.Wrap(ioutil.WriteFile(j.storage, raw, 0644), "can't write to file")
}

//readFromFile returns sequence of bytes which represents books
func (j *JSONStorage) readFromFile() ([]byte, error) {
	raw, err := ioutil.ReadFile(j.storage)
	if err != nil {
		err = errors.Wrap(err, "cant get data from json file")
		log.Println(err)
		return nil, err
	}
	return raw, nil
}

//GetBooks return all books as a list of book struct
func (j *JSONStorage) GetBooks() (Books, error) {
	//dir, _ := os.Getwd()
	//fmt.Println(dir)
	var books Books
	raw, err := j.readFromFile()
	if err != nil {
		err = errors.Wrap(err, "Cant get books")
		log.Println(err)
		return nil, err
	}
	err = json.Unmarshal(raw, &books)
	if err != nil {
		err = errors.Wrap(err, "Cant unmarshal books")
		log.Println(err)
		return nil, err
	}
	return books, err
}

//SaveBook save only one new book to storage
func (j *JSONStorage) SaveBook(book Book) error {
	books, err := j.GetBooks()
	if err != nil {
		err = errors.Wrap(err, "can't get book data")
		log.Println(err)
		return nil
	}
	books = append(books, book)
	return errors.Wrap(j.writeToFile(books), "can't save book data")
}

//GetBook ...
func (j *JSONStorage) GetBook(bookID string) (Book, error) {
	book := Book{}
	books, err := j.GetBooks()
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
	books, err := j.GetBooks()
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
	books, err := j.GetBooks()
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
	books, err := j.GetBooks()
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
