package storage

import (
    //"encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

type Book struct{
    Id string
    Title string
    Ganres []string
    Pages int
    Price float64
}

func GetBooks() []byte{
    //dir, _ := os.Getwd()
    //fmt.Println(dir)
    raw, err := ioutil.ReadFile("storage/Books.json")
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    return raw
}
