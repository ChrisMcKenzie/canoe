package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/header", func(w http.ResponseWriter, r *http.Request) {
		log.Println("sending header")
		w.Header().Set("Etag", "\"somekey\"")
		fmt.Fprintf(w, "<template> <p> Hello, from Header </p> </template>")
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
