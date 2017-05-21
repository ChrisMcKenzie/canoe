package html

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
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
	Prime(http.Pusher, <-chan Fragment)
	Render(w http.ResponseWriter, id string)
}

type HTTPFragmentService struct {
	fragments map[string]Fragment
	client    *http.Client

	pusher http.Pusher

	mu sync.Mutex
}

func NewHTTPFragmentService() *HTTPFragmentService {
	fs := &HTTPFragmentService{
		fragments: make(map[string]Fragment),
		client:    DefaultClient,
		mu:        sync.Mutex{},
	}

	return fs
}

func (fs *HTTPFragmentService) Prime(p http.Pusher, fragments <-chan Fragment) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	for fragment := range fragments {
		fs.fragments[fragment.ID] = fragment
		p.Push("/fragment/"+fragment.ID, nil)
	}
}

// func (fs *HTTPFragmentService) updateCache(id, url string) error {
// 	res, err := fs.get(url)
// 	if err != nil {
// 		return err
// 	}
//
// 	fs.cache[id] = res
// 	return nil
// }

func (fs *HTTPFragmentService) check(id, url string) (bool, error) {
	res, err := fs.client.Head(url)
	if err != nil {
		return false, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return false, fmt.Errorf("unable to retrieve fragment: status: %d", res.StatusCode)
	}

	if res.Header.Get("Etag") == id {
		return false, nil
	}

	return true, nil
}

func (fs *HTTPFragmentService) get(url string) ([]byte, error) {
	res, err := fs.client.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("unable to retrieve fragment: status: %d", res.StatusCode)
	}

	byts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return byts, nil
}

func (fs *HTTPFragmentService) Render(w http.ResponseWriter, id string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fragment := fs.fragments[id]
	res, err := fs.get(fragment.Href)
	if err != nil {
		http.Error(w, "unable to resolve fragment", http.StatusNotFound)
	} else {
		w.Write(res)
	}
}
