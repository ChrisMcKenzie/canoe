package html

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"sync"
)

var (
	DefaultClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
)

type FragmentService interface {
	Render(w http.ResponseWriter, id string)
	Prime(id, url string)
}

type HTTPFragmentService struct {
	cache  map[string]io.Reader
	Client *http.Client

	mu sync.Mutex
}

func NewHTTPFragmentService() *HTTPFragmentService {
	return &HTTPFragmentService{make(map[string]io.Reader), DefaultClient, sync.Mutex{}}
}

func (fs *HTTPFragmentService) Prime(id, url string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if _, ok := fs.cache[id]; !ok {
		res, err := fs.Client.Get(url)
		if err != nil {
			log.Println(err)
		}

		if res.StatusCode < http.StatusBadRequest {
			fs.cache[id] = res.Body
		}
	}
}

func (fs *HTTPFragmentService) Render(w http.ResponseWriter, id string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if val, ok := fs.cache[id]; !ok {
		http.Error(w, "unable to find fragment by given id", http.StatusNotFound)
	} else {
		io.Copy(w, val)
	}
}
