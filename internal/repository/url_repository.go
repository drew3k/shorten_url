package repository

import (
	"shortUrl/shorten_url/internal/domain"
	"sync"
)

type URLRepository interface {
	Create(url *domain.URL) error
	Get(id int) (*domain.URL, error)
}

type InMemoryURLRepository struct {
	urls    map[int]*domain.URL
	counter int
	mu      sync.Mutex
}

func NewInMemoryURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{
		urls:    make(map[int]*domain.URL),
		counter: 1,
	}
}

func (r *InMemoryURLRepository) Create(url *domain.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	url.ID = r.counter

	r.urls[url.ID] = url
	r.counter++

	return nil
}

func (r *InMemoryURLRepository) Get(id int) (*domain.URL, error) {
	url, ok := r.urls[id]
	if !ok {
		return nil, nil
	}
	return url, nil
}
