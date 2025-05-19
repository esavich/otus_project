package service

import (
	"fmt"
	"image"
	"log/slog"
	"net/http"
)

type disckCache interface {
	Set(key string, data image.Image) error
	Get(key string) (image.Image, bool)
}
type CachedImageService struct {
	cache disckCache
	is    ImageGetter
}

func NewCachedImageService(is ImageGetter, dc disckCache) *CachedImageService {
	return &CachedImageService{
		is:    is,
		cache: dc,
	}
}

func (svc *CachedImageService) GetResizedImage(
	width, height int,
	imgURL string,
	header http.Header,
) (image.Image, error) {
	key := fmt.Sprintf("%d-%d-%s", width, height, imgURL)

	slog.Debug("Cache key:" + key)

	slog.Info(fmt.Sprintf("Trying to get image from cache: %s", key))

	if cachedImg, found := svc.cache.Get(key); found {
		slog.Info(fmt.Sprintf("Cache hit: %s", key))
		return cachedImg, nil
	}

	slog.Info("Cache miss, downloading image")

	resizedImage, err := svc.is.GetResizedImage(width, height, imgURL, header)
	if err != nil {
		return nil, err
	}

	// cache the resized image
	err = svc.cache.Set(key, resizedImage)
	if err != nil {
		return nil, err
	}
	slog.Info(fmt.Sprintf("Cache set: %s", key))

	return resizedImage, nil
}
