package service

import (
	"errors"
	"github.com/google/uuid"
	"shortUrl/shorten_url/internal/domain"
	"shortUrl/shorten_url/internal/repository"
	"time"
)

type URLService interface {
	Create(original string) (*domain.URL, error)
	Get(id string) (*domain.URL, error)
}

type UrlService struct {
	repo repository.URLRepository
}

func NewUrlRepository(repo repository.URLRepository) *UrlService {
	return &UrlService{repo: repo}
}

func (s *UrlService) Create(original string) (*domain.URL, error) {
	url := &domain.URL{
		ID:        uuid.New().String(),
		Original:  original,
		Shortened: "",
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *UrlService) Get(id string) (*domain.URL, error) {
	url, err := s.repo.Get(id)
	if err != nil {
		return url, errors.New("such url doesn't exist")
	}

	return url, nil
}
