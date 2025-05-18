package service

import (
	"fmt"
	"image"
	"log/slog"
	"net/http"

	"github.com/esavich/otus_project/internal/downloader"
	"github.com/esavich/otus_project/internal/resizer"
)

type ImageGetter interface {
	GetResizedImage(width, height int, imgURL string, header http.Header) (image.Image, error)
}

type SimpleImageService struct {
	dl *downloader.Downloader
	rz *resizer.Resizer
}

func NewSimpleImageService() *SimpleImageService {
	dl := downloader.NewDownloader()
	rz := resizer.NewResizer()
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
	resized, _ := svc.rz.ResizeImg(img, width, height)

	return resized, nil
}
