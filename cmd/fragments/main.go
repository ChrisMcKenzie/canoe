package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("examples/fragments")))

	log.Fatal(http.ListenAndServe(":8081", nil))
}
