package service

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"shortUrl/shorten_url/internal/domain"
	"shortUrl/shorten_url/internal/repository"
	"time"
)

type URLService interface {
	Create(original string) (*domain.URL, error)
}

type UrlService struct {
	repo repository.URLRepository
}

func NewUrlService(repo repository.URLRepository) *UrlService {
	return &UrlService{repo}
}

func generateHash(original string) string {
	hashed := sha1.New()
	hashed.Write([]byte(original))
	hash := hex.EncodeToString(hashed.Sum(nil))

	return hash[:5]
}

func (s *UrlService) Create(original string) (*domain.URL, error) {
	shortened := generateHash(original)
	customDomain := "https://shrt.ly/"
	url := &domain.URL{
		Original:  original,
		Shortened: fmt.Sprintf("%s%s", customDomain, shortened),
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(url)
	if err != nil {
		return nil, err
	}

	return url, nil
}
