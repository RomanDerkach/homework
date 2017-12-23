package main

import (
	"flag"

	"github.com/RomanDerkach/homework/api"
)

func main() {
	flag.Parse()
	api.Server()
}
