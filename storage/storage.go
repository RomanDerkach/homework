package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const DB = "storage/Books.json"

type Book struct {
	ID     string   `json:"id, omitempty"`
	Title  string   `json:"title, omitempty"`
	Ganres []string `json:"ganres, omitempty"`
	Pages  int      `json:"pages, omitempty"`
	Price  float64  `json:"price, omitempty"`
}

type BookFilter struct {
	Title string `json:"title, omitempty"`
}

func GetDBData() []byte {
	raw, err := ioutil.ReadFile(DB)
	if err != nil {
		log.Println(err)
	}
	return raw
}

func GetBooksData() []Book {
	//dir, _ := os.Getwd()
	//fmt.Println(dir)
	books := []Book{}
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
func SaveBookData(books []Book) {
	raw, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		log.Println(err)
	}
	err = ioutil.WriteFile(DB, raw, 0777)
	if err != nil {
		log.Println(err)
	}
}
