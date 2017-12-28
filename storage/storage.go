package storage

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
)

var storagePath = flag.String("storagePath", "storage/Books.json", "set path the storage file")

var (
	ErrTitleEmpty  = errors.New("There is no title request")
	ErrGenresEmpty = errors.New("There is no genres in request")
	ErrPagesEmpty  = errors.New("There is no pages in request")
	ErrPriceEmpty  = errors.New("There is no price in request")
)

type Book struct {
	ID     string   `json:"id, omitempty"`
	Title  string   `json:"title, omitempty"`
	Genres []string `json:"genres, omitempty"`
	Pages  int      `json:"pages, omitempty"`
	Price  float64  `json:"price, omitempty"`
}

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

type Books []Book

type BookFilter struct {
	Title string `json:"title, omitempty"`
}

func GetDBData() []byte {
	raw, err := ioutil.ReadFile(*storagePath)
	if err != nil {
		// SKIPPED ERROR CHECK!!!!!!!!!!!!!!!!!!!!!!!11111111oneoneone
		log.Println(err)
	}
	return raw
}

func GetBooksData() []Book {
	//dir, _ := os.Getwd()
	//fmt.Println(dir)
	var books Books
	raw := GetDBData()
	err := json.Unmarshal(raw, &books)
	if err != nil {
		log.Println(err)
	}
	return books
}

func SaveNewBook(book Book) error {
	books := append(GetBooksData(), book)
	return errors.Wrap(SaveBookData(books), "can't save book data")
}

//SaveBookData saves all changes with books
func SaveBookData(books Books) error {
	raw, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		return errors.Wrap(err, "can't marshal json")
	}
	return errors.Wrap(ioutil.WriteFile(*storagePath, raw, 0777), "can't write to file")
}
