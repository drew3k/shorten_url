package repository

import "shortUrl/shorten_url/internal/domain"

type URLRepository interface {
	Create(url *domain.URL) error
	Get(id string) (*domain.URL, error)
}

type InMemoryURLRepository struct {
	urls map[string]*domain.URL
}

func NewInMemoryURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{
		urls: make(map[string]*domain.URL),
	}
}

func (r *InMemoryURLRepository) Create(url *domain.URL) error {
	r.urls[url.ID] = url
	return nil
}

func (r *InMemoryURLRepository) Get(id string) (*domain.URL, error) {
	url, ok := r.urls[id]
	if !ok {
		return nil, nil
	}
	return url, nil
}
