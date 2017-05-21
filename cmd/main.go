package main

import (
	"log"
	"net/http"

	"github.com/chrismckenzie/canoe"
)

func main() {
	c := canoe.NewMemoryCache()
	r := canoe.NewHTTPFragmentResolver(http.DefaultClient, c)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/fragments/", canoe.NewFragmentHandler(canoe.WithHTTPResolver(r)))
	http.Handle("/", canoe.NewHandler("/fragments/", http.Dir("examples"), c))

	err := http.ListenAndServeTLS(":18443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal(err)
	}
}
