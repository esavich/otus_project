package downloader

import (
	"context"
	"fmt"
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

func (d *Downloader) Download(url string, header http.Header) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("cant create request: %w", err)
	}

	req.Header = header

	slog.Info(fmt.Sprintf("Downloading: %s (%s) with headers: %v ", url, req.URL.String(), req.Header))
	resp, err := d.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cant do request: %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("cant read response body: %w", err)
	}

	return body, nil
}
