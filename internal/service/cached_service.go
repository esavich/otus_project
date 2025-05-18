package service

import (
	"fmt"
	"image"
	"log/slog"
	"net/http"

	"github.com/esavich/otus_project/internal/cache"
	_ "github.com/esavich/otus_project/internal/cache"
)

type CachedImageService struct {
	cache cache.Cache
	is    ImageGetter
}

func NewCachedImageService(is ImageGetter) *CachedImageService {
	return &CachedImageService{
		is:    is,
		cache: cache.NewCache(5),
	}
}

func (svc *CachedImageService) GetResizedImage(
	width, height int,
	imgURL string,
	header http.Header,
) (image.Image, error) {
	key := cache.Key(fmt.Sprintf("%d-%d-%s", width, height, imgURL))

	slog.Info("Cache key:" + string(key))

	slog.Info(fmt.Sprintf("Trying to get image from cache: %s", key))
	if cachedImg, found := svc.cache.Get(key); found {
		slog.Info(fmt.Sprintf("Cache hit: %s", key))
		return cachedImg.(image.Image), nil
	}

	slog.Info("Cache miss, downloading image")

	resizedImage, err := svc.is.GetResizedImage(width, height, imgURL, header)
	if err != nil {
		return nil, err
	}

	// cache the resized image
	svc.cache.Set(key, resizedImage)
	slog.Info(fmt.Sprintf("Cache set: %s", key))

	return resizedImage, nil
}
