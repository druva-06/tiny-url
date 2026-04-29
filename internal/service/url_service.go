package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/druva-06/tiny-url/internal/dto/request"
	"github.com/druva-06/tiny-url/internal/dto/response"
	"github.com/druva-06/tiny-url/internal/repository"
	"github.com/druva-06/tiny-url/internal/repository/cache"
	"github.com/druva-06/tiny-url/internal/util"
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
	shortCode := util.GenerateShortCode()
	_, err := s.repo.Create(shortCode, longUrl)
	if err != nil {
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
	} else {
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
		err := s.rdb.Set(ctx, cacheKey, longUrl, 24*time.Hour)
		if err != nil {
			log.Printf("[GetOriginalURL] REDIS SET FAILED code=%s err=%v", shortCode, err.Error())
		} else {
			log.Printf("[GetOriginalURL] REDIS SET SUCCESS code=%s", shortCode)
		}
	}()
	log.Printf("[GetOriginalURL] END code=%s", shortCode)
	return longUrl, err
}

func (s *URLService) UpdateLongUrl(ctx context.Context, shortCode string, request request.CreateShortURLRequest) (response.ShortUrlResponse, error) {
	longUrl := request.LongUrl
	cacheKey := "url:short:" + shortCode
	log.Printf("[UpdateLongUrl] SERVICE START code=%s long_url=%s", shortCode, longUrl)
	_, err := s.repo.GetLongURL(shortCode)
	if err == sql.ErrNoRows {
		log.Printf("[UpdateLongUrl] SHORT_URL NOT EXIST code=%s", shortCode)
		return response.ShortUrlResponse{}, errors.New("short code not found")
	} else if err != nil {
		return response.ShortUrlResponse{}, err
	}
	if err = s.repo.UpdateLongUrl(shortCode, request.LongUrl); err != nil {
		return response.ShortUrlResponse{}, err
	}
	go func() {
		log.Printf("[UpdateLongUrl] CACHE START code=%s long_url=%s", shortCode, longUrl)
		exists, err := s.rdb.Exists(ctx, cacheKey)
		if err != nil {
			log.Printf("[UpdateLongUrl] REDIS ERROR cache_key=%s err=%v", cacheKey, err.Error())
			return
		}
		if exists != 1 {
			log.Printf("[UpdateLongUrl] REDIS MISS cache_key=%s long_url=%s", cacheKey, longUrl)
			return
		}
		err = s.rdb.Set(ctx, cacheKey, longUrl, 24*time.Hour)
		if err != nil {
			log.Printf("[UpdateLongUrl] REDIS SET FAILED code=%s err=%v", shortCode, err.Error())
		} else {
			log.Printf("[UpdateLongUrl] REDIS SET SUCCESS cache_key=%s", cacheKey)
		}
	}()
	return response.ShortUrlResponse{ShortCode: shortCode, LongUrl: longUrl}, nil
}
