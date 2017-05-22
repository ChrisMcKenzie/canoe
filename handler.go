package canoe

import (
	"net/http"
	"path"

	"github.com/chrismckenzie/canoe/html"
)

// FragmentCache ...
type FragmentCache interface {
	Add(html.Fragment)
	Get(id string) (html.Fragment, bool)
}

// Handler is an http.Handler that will parse html for fragments and resolve them
// concurrently and push them to the client as they resolve. The client will
// use a service-worker to handle and compile the resulting page.
//
// In the case of HTTP/1.1 connections we will render the page on the server
// side and send it out.
type Handler struct {
	fs           http.FileSystem
	fc           FragmentCache
	fragmentPath string
}

// NewHandler creates a new canoe.Handler with the given http.FileSystem
func NewHandler(fp string, fs http.FileSystem, fc FragmentCache) http.Handler {
	return &Handler{fs, fc, fp}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Unlock HTTP/2 server push
	p, h2 := w.(http.Pusher)
	if h2 {
		p.Push("/static/canoe.js", nil)
	}

	f, err := h.fs.Open("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	hp, err := html.NewParser(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	for fragment := range hp.Fragments() {
		h.fc.Add(fragment)
		if h2 {
			p.Push(path.Join(h.fragmentPath, fragment.ID), nil)
		}
	}

	hp.Render(w)
}
