package storage

import (
	"fmt"

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
func NewSQLStorage(sqlDB *gorm.DB) *SQLStorage {
	return &SQLStorage{
		storage: sqlDB,
	}
}

//GetBooks return objects of Books
func (s *SQLStorage) GetBooks() (Books, error) {
	books := Books{}
	return books, s.storage.Find(&books).Error
}

//GetBook get book by it's id
func (s *SQLStorage) GetBook(bookID string) (Book, error) {
	book := Book{}
	err := s.storage.Where("id = ?", bookID).First(&book).Error
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
	err := s.storage.Where("id = ?", bookID).Delete(&Book{}).Error
	//when deleting it wont return that nothing found
	fmt.Println(err)
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
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
	books := Books{}
	return books, s.storage.Where("title LIKE ?", fmt.Sprintf("%%%s%%", filter.Title)).Find(&books).Error
}
