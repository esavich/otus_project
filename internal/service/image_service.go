package service

import (
	"fmt"
	"image"
	"log/slog"
	"net/http"
)

type ImageGetter interface {
	GetResizedImage(width, height int, imgURL string, header http.Header) (image.Image, error)
}

type resizer interface {
	ResizeImg(img image.Image, w int, h int) image.Image
}

type downloader interface {
	Download(imgURL string, header http.Header) (image.Image, error)
}

type SimpleImageService struct {
	dl downloader
	rz resizer
}

func NewSimpleImageService(dl downloader, rz resizer) *SimpleImageService {
	return &SimpleImageService{
		dl: dl,
		rz: rz,
	}
}

func (svc *SimpleImageService) GetResizedImage(
	width, height int,
	imgURL string,
	header http.Header,
) (image.Image, error) {
	img, err := svc.dl.Download(imgURL, header)
	if err != nil {
		err = fmt.Errorf("failed to download image: %w", err)
		slog.Error(err.Error())
		return nil, err
	}
	slog.Info("Image downloaded")
	slog.Info("Resizing image", slog.Int("width", width), slog.Int("height", height))
	resized := svc.rz.ResizeImg(img, width, height)

	return resized, nil
}
