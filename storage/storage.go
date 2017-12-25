package storage

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

var storagePath = flag.String("storagePath", "storage/Books.json", "set path the storage file")

type Book struct {
	ID     string   `json:"id, omitempty"`
	Title  string   `json:"title, omitempty"`
	Ganres []string `json:"ganres, omitempty"`
	Pages  int      `json:"pages, omitempty"`
	Price  float64  `json:"price, omitempty"`
}

type Books []Book

type BookFilter struct {
	Title string `json:"title, omitempty"`
}

func GetDBData() []byte {
	raw, err := ioutil.ReadFile(*storagePath)
	if err != nil {
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

func SaveNewBook(book Book) {
	books := append(GetBooksData(), book)
	SaveBookData(books)
}

//SaveBookData saves all changes with books
func SaveBookData(books Books) {
	raw, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		log.Println(err)
	}
	err = ioutil.WriteFile(*storagePath, raw, 0777)
	if err != nil {
		log.Println(err)
	}
}
