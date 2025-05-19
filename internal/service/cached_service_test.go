package service

import (
	"errors"
	"image"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(key string) (image.Image, bool) {
	args := m.Called(key)
	img := args.Get(0)
	if img == nil {
		return nil, args.Bool(1)
	}
	return img.(image.Image), args.Bool(1)
}

func (m *MockCache) Set(key string, img image.Image) error {
	args := m.Called(key, img)
	return args.Error(0)
}

type MockImageGetter struct {
	mock.Mock
}

func (m *MockImageGetter) GetResizedImage(width, height int, imgURL string, header http.Header) (image.Image, error) {
	args := m.Called(width, height, imgURL, header)
	img := args.Get(0)
	if img == nil {
		return nil, args.Error(1)
	}
	return img.(image.Image), args.Error(1)
}

func TestCachedImageService_GetResizedImage_CacheHit(t *testing.T) {
	cache := new(MockCache)
	imageGetter := new(MockImageGetter)
	svc := NewCachedImageService(imageGetter, cache)

	headers := http.Header{}
	key := "50-60-" + testImgURL
	cachedImg := image.NewRGBA(image.Rect(0, 0, 50, 60))

	cache.On("Get", key).Return(cachedImg, true)

	result, err := svc.GetResizedImage(50, 60, testImgURL, headers)
	require.NoError(t, err)
	require.Equal(t, cachedImg, result)

	cache.AssertCalled(t, "Get", key)
	imageGetter.AssertNotCalled(t, "GetResizedImage")
}

func TestCachedImageService_GetResizedImage_CacheMiss_Success(t *testing.T) {
	cache := new(MockCache)
	imageGetter := new(MockImageGetter)
	svc := NewCachedImageService(imageGetter, cache)

	headers := http.Header{}
	key := "50-60-" + testImgURL
	resizedImg := image.NewRGBA(image.Rect(0, 0, 50, 60))

	cache.On("Get", key).Return(nil, false)
	imageGetter.On("GetResizedImage", 50, 60, testImgURL, headers).Return(resizedImg, nil)
	cache.On("Set", key, resizedImg).Return(nil)

	result, err := svc.GetResizedImage(50, 60, testImgURL, headers)
	require.NoError(t, err)
	require.Equal(t, resizedImg, result)

	cache.AssertCalled(t, "Get", key)
	imageGetter.AssertCalled(t, "GetResizedImage", 50, 60, testImgURL, headers)
	cache.AssertCalled(t, "Set", key, resizedImg)
}

func TestCachedImageService_GetResizedImage_CacheMiss_ExternalError(t *testing.T) {
	cache := new(MockCache)
	imageGetter := new(MockImageGetter)
	svc := NewCachedImageService(imageGetter, cache)

	headers := http.Header{}
	imgURL := "http://example.com/image.jpg"
	key := "50-60-" + imgURL

	cache.On("Get", key).Return(nil, false)
	imageGetter.On("GetResizedImage", 50, 60, imgURL, headers).Return(nil, errors.New("external error"))

	result, err := svc.GetResizedImage(50, 60, imgURL, headers)
	require.Error(t, err)
	require.ErrorContains(t, err, "external error")
	require.Nil(t, result)

	cache.AssertCalled(t, "Get", key)
	imageGetter.AssertCalled(t, "GetResizedImage", 50, 60, imgURL, headers)
	cache.AssertNotCalled(t, "Set", mock.Anything, mock.Anything)
}

func TestCachedImageService_GetResizedImage_CacheMiss_CacheSetError(t *testing.T) {
	cache := new(MockCache)
	imageGetter := new(MockImageGetter)
	svc := NewCachedImageService(imageGetter, cache)

	headers := http.Header{}
	key := "50-60-" + testImgURL
	resizedImg := image.NewRGBA(image.Rect(0, 0, 50, 60))

	cache.On("Get", key).Return(nil, false)
	imageGetter.On("GetResizedImage", 50, 60, testImgURL, headers).Return(resizedImg, nil)
	cache.On("Set", key, resizedImg).Return(errors.New("cache set error"))

	result, err := svc.GetResizedImage(50, 60, testImgURL, headers)
	require.Error(t, err)
	require.ErrorContains(t, err, "cache set error")
	require.Nil(t, result)

	cache.AssertCalled(t, "Get", key)
	imageGetter.AssertCalled(t, "GetResizedImage", 50, 60, testImgURL, headers)
	cache.AssertCalled(t, "Set", key, resizedImg)
}
