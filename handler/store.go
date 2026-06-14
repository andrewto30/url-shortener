package handler

import "sync"

type URLStore struct {
	mu   sync.Mutex
	urls map[string]string
}

func NewURLStore() *URLStore {
	return &URLStore{
		urls: make(map[string]string),
	}
}

func (u *URLStore) Save(key, url string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.urls[key] = url
}

func (u *URLStore) Get(key string) (string, bool) {
	u.mu.Lock()
	defer u.mu.Unlock()

	v, ok := u.urls[key]

	return v, ok
}
