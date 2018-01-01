package main

import (
	"flag"

	"github.com/jinzhu/gorm"

	"log"

	"github.com/RomanDerkach/homework/api"
	"github.com/RomanDerkach/homework/storage"
)

var (
	storagePath = flag.String("storagePath", "storage/Books.json", "set path the storage file")
	useSQL      = flag.Bool("useSQL", false, "use sql db instead of simple json file")
)

func main() {
	var (
		dbStorage *gorm.DB
		err       error
		store     api.Storage
	)

	defer func() {
		if dbStorage != nil {
			err = dbStorage.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	flag.Parse()

	switch *useSQL {
	case true:
		dbStorage, err = storage.InitDB("sqlite3", "storage/book.db")
		if err != nil {
			log.Fatal(err)
		}
		store = storage.NewSQLStorage(dbStorage)
	case false:
		store = storage.NewJSONStorage(*storagePath)
	}
	handler := api.NewHandler(store)
	log.Fatal(api.Server(handler))
}
