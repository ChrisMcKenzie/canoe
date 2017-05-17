package main

import (
	"log"
	"net/http"

	"github.com/ChrisMcKenzie/canoe"
	"github.com/ChrisMcKenzie/canoe/html"
)

func main() {
	fs := html.NewHTTPFragmentService()
	h := canoe.NewHandler(http.Dir("examples"), fs)
	err := http.ListenAndServeTLS(":18443", "server.crt", "server.key", h)
	if err != nil {
		log.Fatal(err)
	}
}
