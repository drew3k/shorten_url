package service

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"shortUrl/shorten_url/internal/domain"
	"shortUrl/shorten_url/internal/repository"
	"time"
)

type URLService interface {
	Create(original string) (*domain.URL, error)
	Get(id int) (*domain.URL, error)
}

type UrlService struct {
	repo repository.URLRepository
}

func NewUrlRepository(repo repository.URLRepository) *UrlService {
	return &UrlService{repo: repo}
}

func generateHash() string {
	randomString := uuid.New().String()

	hashed := sha1.New()
	hashed.Write([]byte(randomString))
	hash := hex.EncodeToString(hashed.Sum(nil))

	return hash[:6]
}

func (s *UrlService) Create(original string) (*domain.URL, error) {
	url := &domain.URL{
		ID:        1,
		Original:  original,
		Shortened: generateHash(),
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *UrlService) Get(id int) (*domain.URL, error) {
	url, err := s.repo.Get(id)
	if err != nil {
		return url, errors.New("such url doesn't exist")
	}

	return url, nil
}
