package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const DB = "storage/Books.json"

type Book struct {
	ID     string
	Title  string
	Ganres []string
	Pages  int
	Price  float64
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
