package main

import (
	"flag"

	"log"

	"github.com/RomanDerkach/homework/api"
)

func main() {
	flag.Parse()
	log.Fatal(api.Server())
}
