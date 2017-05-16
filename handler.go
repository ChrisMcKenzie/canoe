package canoe

import (
	"fmt"
	"net/http"
	"time"

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

	fragments       map[string]*html.Fragment
	fragmentService html.FragmentService
}

// NewHandler creates a new canoe.Handler with the given http.FileSystem
func NewHandler(fs http.FileSystem) http.Handler {
	return &Handler{fs, make(map[string]*html.Fragment), html.NewHTTPFragmentService()}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Unlock HTTP/2 server push
	p, ok := w.(http.Pusher)
	if !ok {
		// handler HTTP/1.1 case
		// http.Error(w, "http 1.1 not supported", http.StatusInternalServerError)
		fmt.Println("http 1.1 in use")
	}

	fmt.Println(r.URL.Path)
	switch r.URL.Path {
	case "/":
		f, err := h.fs.Open("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		parser, err := html.NewParser(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		for fragment := range parser.Fragments() {
			go h.fragmentService.Prime(fragment.ID, fragment.Href)
			fmt.Printf("pushing %s \n", fragment.ID)
			p.Push("/fragment?id="+fragment.ID, nil)
		}
		parser.Render(w)
	case "/fragment":
		id := r.URL.Query().Get("id")
		fmt.Printf("loading %s\n", id)
		h.fragmentService.Render(w, id)
		// fmt.Fprintf(w, "<template>Fragment ID: %s</template>", id)
	case "/body":
		fmt.Println("sleeping")
		time.Sleep(time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<template> hello </template>"))
	}

}
