package repository

import "shortUrl/shorten_url/internal/domain"

type URLRepository interface {
	Create(url *domain.URL) error
	Get(id string) (*domain.URL, error)
}

type inMemoryURLRepository struct {
	urls map[string]*domain.URL
}

func NewInMemoryURLRepository() URLRepository {
	return &inMemoryURLRepository{
		urls: make(map[string]*domain.URL),
	}
}

func (r *inMemoryURLRepository) Create(url *domain.URL) error {
	r.urls[url.ID] = url
	return nil
}

func (r *inMemoryURLRepository) Get(id string) (*domain.URL, error) {
	url, ok := r.urls[id]
	if !ok {
		return nil, nil
	}
	return url, nil
}
