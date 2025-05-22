package service

import (
	"errors"
	"image"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockDownloader struct {
	mock.Mock
}

func (m *MockDownloader) Download(imgURL string, header http.Header) (image.Image, error) {
	args := m.Called(imgURL, header)

	img := args.Get(0)
	if img == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(image.Image), args.Error(1)
}

type MockResizer struct {
	mock.Mock
}

func (m *MockResizer) ResizeImg(img image.Image, w int, h int) image.Image {
	args := m.Called(img, w, h)
	return args.Get(0).(image.Image)
}

const testImgURL = "http://example.com/image.jpg"

func TestSimpleImageService_GetResizedImage_Success(t *testing.T) {
	mockDownloader := new(MockDownloader)
	mockResizer := new(MockResizer)

	service := NewSimpleImageService(mockDownloader, mockResizer)

	headers := http.Header{"Authorization": []string{"Bearer token"}}

	// Test images
	testImage := image.NewRGBA(image.Rect(0, 0, 100, 100))
	resizedImage := image.NewRGBA(image.Rect(0, 0, 50, 60))

	// setup mocks
	mockDownloader.On("Download", testImgURL, headers).Return(testImage, nil)
	mockResizer.On("ResizeImg", testImage, 50, 60).Return(resizedImage)

	result, err := service.GetResizedImage(50, 60, testImgURL, headers)

	require.NoError(t, err)
	require.Equal(t, resizedImage, result)

	// assert calls
	mockDownloader.AssertCalled(t, "Download", testImgURL, headers)
	mockResizer.AssertCalled(t, "ResizeImg", testImage, 50, 60)
}

func TestSimpleImageService_GetResizedImage_DownloadNil(t *testing.T) {
	mockDownloader := new(MockDownloader)
	mockResizer := new(MockResizer)

	service := NewSimpleImageService(mockDownloader, mockResizer)

	headers := http.Header{"Authorization": []string{"Bearer token"}}

	mockDownloader.On("Download", testImgURL, headers).Return(nil, errors.New("download error"))

	result, err := service.GetResizedImage(50, 60, testImgURL, headers)

	require.Error(t, err)
	require.Nil(t, result)

	mockDownloader.AssertCalled(t, "Download", testImgURL, headers)
	mockResizer.AssertNotCalled(t, "ResizeImg")
}
