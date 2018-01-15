package main

import (
	"flag"

	"log"

	"github.com/RomanDerkach/homework/api"
	"github.com/RomanDerkach/homework/storage"
	"github.com/pkg/errors"
)

var (
	fileStoragePath = flag.String("storagePath", "storage/Books.json", "set path the storage file")
	sqlStoragePath  = flag.String("useSQL", "", "set path the sql storage file")
	host            = flag.String("host", ":8081", "host for http server")
)

func main() {
	var (
		err   error
		store api.Storage
	)

	flag.Parse()

	if len(*sqlStoragePath) != 0 {
		store, err = storage.NewSQLStorage(*sqlStoragePath)
	} else {
		store, err = storage.NewJSONStorage(*fileStoragePath)
	}
	if err != nil {
		log.Fatal(errors.Wrap(err, "cannot create new storage"))
	}

	defer func() {
		if err := store.Close(); err != nil {
			err = errors.Wrap(err, "can't close storage")
			log.Fatal(err)
		}
	}()

	handler, err := api.NewHandler(store)
	if err != nil {
		log.Fatal(errors.Wrap(err, "cannot create new handler"))
	}
	log.Fatal(api.Server(handler, *host))
}
