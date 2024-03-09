package repository

import (
	"shortUrl/shorten_url/internal/domain"
	"sync"
)

type URLRepository interface {
	Create(url *domain.URL) error
}

type InMemoryURLRepository struct {
	urls map[string]*domain.URL
	mu   sync.Mutex
}

func NewInMemoryURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{
		urls: make(map[string]*domain.URL),
	}
}

func (r *InMemoryURLRepository) Create(url *domain.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.urls[url.Shortened] = url
	return nil
}
