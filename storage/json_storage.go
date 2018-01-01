package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"

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

//GetDBData returns sequence of bytes which represents books
func (j *JSONStorage) GetDBData() ([]byte, error) {
	raw, err := ioutil.ReadFile(j.storage)
	if err != nil {
		err = errors.Wrap(err, "cant get data from json file")
		log.Println(err)
		return nil, err
	}
	return raw, nil
}

//GetBooksData return all books as a list of book struct
func (j *JSONStorage) GetBooksData() (Books, error) {
	//dir, _ := os.Getwd()
	//fmt.Println(dir)
	var books Books
	raw, err := j.GetDBData()
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

//SaveNewBook save only one new book to storage
func (j *JSONStorage) SaveNewBook(book Book) error {
	books, err := j.GetBooksData()
	if err != nil {
		// !!!!!!!!!!!!!!! PLEASE ERROR HANDLE
		return nil
	}
	books = append(books, book)
	return errors.Wrap(j.SaveBookData(books), "can't save book data")
}

//SaveBookData saves all changes with books
func (j *JSONStorage) SaveBookData(books Books) error {
	raw, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		return errors.Wrap(err, "can't marshal json")
	}
	return errors.Wrap(ioutil.WriteFile(j.storage, raw, 0777), "can't write to file")
}
