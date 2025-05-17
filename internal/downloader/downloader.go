package downloader

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Downloader struct {
	c *http.Client
}

func NewDownloader() *Downloader {
	return &Downloader{
		c: &http.Client{},
	}
}

func (d *Downloader) Download(url string, header http.Header) (image.Image, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("cant create request: %w", err)
	}

	req.Header = header

	slog.Info(fmt.Sprintf("Downloading: %s  with headers: %+v ", url, req.Header))
	resp, err := d.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cant do request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("cant read response body: %w", err)
	}

	jpegImage, err := jpeg.Decode(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("cant decode jpeg image: %w", err)
	}

	return jpegImage, nil
}
