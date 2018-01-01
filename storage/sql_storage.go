package storage

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type SQLStorage struct {
	// dbtype string
	// dbname string
	storage *gorm.DB
}

func NewSQLStorage(sqlDB *gorm.DB) *SQLStorage {
	return &SQLStorage{
		storage: sqlDB,
	}
}

//GetDBData get byte array of data
func (s *SQLStorage) GetDBData() ([]byte, error) {
	books := Books{}
	err := s.storage.Find(&books).Error
	if err != nil {
		err = errors.Wrap(err, "cant get data from db")
		log.Println(err)
		return nil, err
	}
	fmt.Printf("%v", books)
	raw, err := json.Marshal(books)
	if err != nil {
		err = errors.Wrap(err, "cant convert db data to bytes")
		log.Println(err)
		return nil, err
	}
	return raw, nil
}

//GetBooksData return objects of Books
func (s *SQLStorage) GetBooksData() (Books, error) {
	books := Books{}
	return books, s.storage.Find(&books).Error
}

//SaveNewBook save new books to DB
func (s *SQLStorage) SaveNewBook(book Book) error {
	err := s.storage.Create(&book).Error
	return errors.Wrap(err, "can't save book data")
}

//SaveBookData saves all changes with books
func (s *SQLStorage) SaveBookData(books Books) error {
	// Need to rewrite all the methods up here
	return errors.New("Not implemented")
}
