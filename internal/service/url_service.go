package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/druva-06/tiny-url/internal/repository"
	"github.com/druva-06/tiny-url/internal/repository/cache"
	"github.com/jxskiss/base62"
	"github.com/redis/go-redis/v9"
)

type URLService struct {
	repo *repository.URLRepository
	rdb  *cache.URLCache
}

func NewURLService(r *repository.URLRepository, rdb *cache.URLCache) *URLService {
	return &URLService{repo: r, rdb: rdb}
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

func (s *URLService) GetLongURL(ctx context.Context, shortCode string) (string, error) {

	cacheKey := "url:short:" + shortCode
	log.Printf("[GetOriginalURL] START code=%s cacheKey=%s", shortCode, cacheKey)
	value, err := s.rdb.Get(ctx, cacheKey)
	if err == nil {
		log.Printf("[GetOriginalURL] CACHE HIT code=%s value=%s", shortCode, value)
		return value, err
	}
	if err == redis.Nil {
		log.Printf("[GetOriginalURL] CACHE MISS code=%s", shortCode)
	} else if err != nil {
		log.Printf("[GetOriginalURL] REDIS ERROR code=%s err=%v", shortCode, err.Error())
	}
	log.Printf("[GetOriginalURL] FETCHING FROM DB code=%s", shortCode)
	longUrl, err := s.repo.GetLongURL(shortCode)
	if err != nil {
		log.Printf("[GetOriginalURL] DB ERROR code=%s err=%v", shortCode, err)
		return "", err
	}
	if longUrl == "" {
		log.Printf("[GetOriginalURL] NOT FOUND code=%s", shortCode)
		return "", nil
	}
	log.Printf("[GetOriginalURL] DB HIT code=%s url=%s", shortCode, longUrl)
	go func() {
		err := s.rdb.Set(context.Background(), cacheKey, longUrl, 24*time.Hour)
		if err != nil {
			log.Printf("[GetOriginalURL] REDIS SET FAILED code=%s err=%v", shortCode, err.Error())
		} else {
			log.Printf("[GetOriginalURL] REDIS SET SUCCESS code=%s", shortCode)
		}
	}()
	log.Printf("[GetOriginalURL] END code=%s", shortCode)
	return longUrl, err
}

func (s *URLService) DeteleShortURL(ctx context.Context, shortcode string) error {

	err := s.repo.DeteleShortURL(shortcode)
	if err != nil {
		log.Printf("[DeleteShortURL] DB ERROR code=%s err=%v", shortcode, err)
		return err
	}
	cacheKey := "url:short:" + shortcode
	deleted, error := s.rdb.Del(ctx, cacheKey)
	if error != nil {
		log.Printf("[DeleteShortURL]Issue when deleting shortcode %s", shortcode)
		return error
	} else {
		log.Printf("[DeleteShortURL]successfully Deleted %v rows and key %s", deleted, shortcode)
	}
	return nil
}
