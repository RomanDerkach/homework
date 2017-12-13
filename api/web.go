package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    // "github.com/gorilla/mux"
    "github.com/RomanDerkach/homework/storage"
)


func Server() {
    http.HandleFunc("/books", BooksHandler)
    http.ListenAndServe(":8081", nil)
}

func BooksHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        fmt.Println("AAA")
    } else {
        jsonBody, err := json.Marshal(storage.GetBooks())
        if err != nil {
            http.Error(w, "Error converting results to json",
            http.StatusInternalServerError)
        } 
        w.Write(jsonBody)
    }    
}