package storage

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

//SQLStorage describes a storage
type SQLStorage struct {
	// dbtype string
	// dbname string
	storage *gorm.DB
}

//NewSQLStorage constructor for SQL storage
func NewSQLStorage(sqlDBpath string) (*SQLStorage, error) {
	dbStorage, err := InitDB("sqlite3", sqlDBpath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize DB")
	}
	return &SQLStorage{storage: dbStorage}, nil
}

//Close function close DB connection
func (s *SQLStorage) Close() error {
	return s.storage.Close()
}

//GetBooks return objects of Books
func (s *SQLStorage) GetBooks() (Books, error) {
	var books Books
	return books, s.storage.Find(&books).Error
}

//GetBook get book by it's id
func (s *SQLStorage) GetBook(bookID string) (Book, error) {
	var book Book
	err := s.storage.First(&book, bookID).Error
	if err == gorm.ErrRecordNotFound {
		return book, ErrNotFound
	}
	return book, err
}

//SaveBook save new books to DB
func (s *SQLStorage) SaveBook(book Book) error {
	err := s.storage.Create(&book).Error
	return errors.Wrap(err, "can't save book data")
}

//DeleteBook removes a book from storage
func (s *SQLStorage) DeleteBook(bookID string) error {
	query := s.storage.Delete(&Book{ID: bookID})
	if query.Error != nil {
		return errors.Wrap(query.Error, "can't delete book")
	}
	if query.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

//UpdateBook update a book in a storage
func (s *SQLStorage) UpdateBook(bookID string, updBook Book) (Book, error) {
	err := s.storage.Model(&Book{}).Where("id = ?", bookID).Update(&updBook).First(&updBook).Error
	if err == gorm.ErrRecordNotFound {
		return Book{}, ErrNotFound
	}
	return updBook, err
}

//FilterBooks get books that follow the filter
func (s *SQLStorage) FilterBooks(filter BookFilter) (Books, error) {
	var books Books
	return books, s.storage.Where("lower(title) LIKE ?", "%"+strings.ToLower(filter.Title)+"%").Find(&books).Error
}
