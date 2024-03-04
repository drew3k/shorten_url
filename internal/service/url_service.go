package service

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"shortUrl/shorten_url/internal/domain"
	"shortUrl/shorten_url/internal/repository"
	"time"
)

type URLService interface {
	Create(original string) (*domain.URL, error)
	Get(shortened string) (*domain.URL, error)
}

type UrlService struct {
	repo repository.URLRepository
}

func NewUrlRepository(repo repository.URLRepository) *UrlService {
	return &UrlService{repo: repo}
}

func generateHash(original string) string {
	hashed := sha1.New()
	hashed.Write([]byte(original))
	hash := hex.EncodeToString(hashed.Sum(nil))

	return hash[:6]
}

func (s *UrlService) Create(original string) (*domain.URL, error) {
	shortened := generateHash(original)
	url := &domain.URL{
		Original:  original,
		Shortened: fmt.Sprintf("https://test/%s", shortened),
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *UrlService) Get(shortened string) (*domain.URL, error) {
	url, err := s.repo.Get(shortened)
	if err != nil {
		return url, errors.New("such url doesn't exist")
	}

	return url, nil
}
