package main

import (
	"log"
	"net/http"

	"github.com/ChrisMcKenzie/canoe"
)

func main() {
	h := canoe.NewHandler(http.Dir("examples"))
	err := http.ListenAndServeTLS(":18443", "server.crt", "server.key", h)
	if err != nil {
		log.Fatal(err)
	}
}
