package storage

import (
	"errors"

	"github.com/lib/pq"
)

type Book struct {
	ID     string         `gorm:"primary_key" json:"id, omitempty"`
	Title  string         `gorm:"not null"    json:"title, omitempty"`
	Genres pq.StringArray `gorm:"type:varchar(64)" json:"genres, omitempty"`
	Pages  int            `gorm:"not null"    json:"pages, omitempty"`
	Price  float64        `gorm:"not null"    json:"price, omitempty"`
}

//ErrNotFound for the error when we haven't found something
var ErrNotFound = errors.New("can't find the book with given ID")

var (
	//ErrTitleEmpty error on book validation
	ErrTitleEmpty  = errors.New("There is no title request")
	ErrGenresEmpty = errors.New("There is no genres in request")
	ErrPagesEmpty  = errors.New("There is no pages in request")
	ErrPriceEmpty  = errors.New("There is no price in request")
)

//Books is a list of book structures
type Books []Book

//BookFilter struct for book filtering
type BookFilter struct {
	Title string `json:"title, omitempty"`
}

//Validate ..
func (b Book) Validate() (err error) {
	if b.Title == "" {
		return ErrTitleEmpty
	}
	if len(b.Genres) == 0 {
		return ErrGenresEmpty
	}
	if b.Pages == 0 {
		return ErrPagesEmpty
	}
	if b.Price == 0 {
		return ErrPriceEmpty
	}
	return
}
