package storage

import (
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func GetBooksData() []byte {
	//dir, _ := os.Getwd()
	//fmt.Println(dir)
	raw, err := ioutil.ReadFile("storage/Books.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return raw
}
