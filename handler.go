package canoe

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ChrisMcKenzie/canoe/html"
)

// Handler is an http.Handler that will parse html for fragments and resolve them
// concurrently and push them to the client as they resolve. The client will
// use a service-worker to handle and compile the resulting page.
//
// In the case of HTTP/1.1 connections we will render the page on the server
// side and send it out.
type Handler struct {
	fs http.FileSystem

	fragmentService html.FragmentService
}

// NewHandler creates a new canoe.Handler with the given http.FileSystem
func NewHandler(fs http.FileSystem, frs html.FragmentService) http.Handler {
	return &Handler{fs, frs}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Unlock HTTP/2 server push
	p, ok := w.(http.Pusher)
	if !ok {
		// handler HTTP/1.1 case
		// http.Error(w, "http 1.1 not supported", http.StatusInternalServerError)
		fmt.Println("http 1.1 in use")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f, err := h.fs.Open("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		parser, err := html.NewParser(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		h.fragmentService.Prime(p, parser.Fragments())
		parser.Render(w)
	})

	mux.HandleFunc("/fragment/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/fragment/"):]

		log.Printf("rendering fragment: %s\n", id)
		log.Printf("loading %s\n", id)
		h.fragmentService.Render(w, id)
	})

	mux.HandleFunc("/body", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("sleeping")
		// time.Sleep(time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<template> hello </template>"))
	})

	mux.ServeHTTP(w, r)
}
