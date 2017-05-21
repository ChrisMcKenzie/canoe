package canoe

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/chrismckenzie/canoe/html"
)

// FragmentResolver is used to retrieve fragment.
type FragmentResolver interface {
	ResolveByID(id string) ([]byte, error)
}

// FragmentHandler ...
type FragmentHandler struct {
	fr FragmentResolver
}

// FragmentHandlerOption ...
type FragmentHandlerOption func(*FragmentHandler)

// NewFragmentHandler ...
func NewFragmentHandler(opts ...FragmentHandlerOption) *FragmentHandler {
	var fh FragmentHandler
	for _, opt := range opts {
		opt(&fh)
	}

	return &fh
}

func (fh *FragmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/fragments/"):]
	if len(id) <= 0 {
		http.Error(w, "id in path is required", http.StatusBadRequest)
		return
	}

	fragment, err := fh.fr.ResolveByID(id)
	if err != nil {
		http.Error(w, "unable to retrieve fragment", http.StatusNotFound)
	}

	w.Write(fragment)
}

// WithHTTPResolver ...
func WithHTTPResolver(fr *HTTPFragmentResolver) func(fh *FragmentHandler) {
	return func(fh *FragmentHandler) {
		fh.fr = fr
	}
}

// HTTPFragmentResolver ...
type HTTPFragmentResolver struct {
	client *http.Client
	cache  FragmentCache
}

// NewHTTPFragmentResolver ...
func NewHTTPFragmentResolver(c *http.Client, fc FragmentCache) *HTTPFragmentResolver {
	return &HTTPFragmentResolver{c, fc}
}

// ResolveByID ...
func (fr HTTPFragmentResolver) ResolveByID(id string) ([]byte, error) {
	if fragment, ok := fr.cache.Get(id); ok {
		res, err := fr.client.Get(fragment.Href)
		if err != nil {
			return nil, err
		}

		byt, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return byt, nil
	}

	return nil, fmt.Errorf("no fragment exist by given id: %s", id)
}

type MemoryCache struct {
	store map[string]html.Fragment
	mu    sync.Mutex
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{make(map[string]html.Fragment), sync.Mutex{}}
}

func (c *MemoryCache) Add(f html.Fragment) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[f.ID] = f
}

func (c *MemoryCache) Get(id string) (html.Fragment, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f, ok := c.store[id]
	return f, ok
}
