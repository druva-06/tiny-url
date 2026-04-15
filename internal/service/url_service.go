package service

import (
	"strconv"

	"github.com/druva-06/tiny-url/internal/repository"
	"github.com/jxskiss/base62"
)

type URLService struct {
	repo *repository.URLRepository
}

func NewURLService(r *repository.URLRepository) *URLService {
	return &URLService{repo: r}
}

func (s *URLService) CreateShortURL(longUrl string) (string, error) {
	id, err := s.repo.Create(longUrl)
	if err != nil {
		return "", err
	}
	shortCode := base62.EncodeToString([]byte(strconv.FormatInt(id, 10)))
	if err = s.repo.UpdateShortCode(id, shortCode); err != nil {
		return "", err
	}
	return shortCode, nil
}

func (s *URLService) GetLongURL(shortCode string) (string, error) {
	longUrl, err := s.repo.GetLongURL(shortCode)
	if err != nil {
		return "", err
	}
	return longUrl, err
}
