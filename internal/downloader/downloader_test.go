package downloader

import (
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var headers = http.Header{
	"Authorization": []string{"Token mock-token"},
	"User-Agent":    []string{"Mozilla/5.0"},
}

func TestDownloader_OK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/image.jpg", r.URL.Path)

		assert.Equal(t, headers.Get("Authorization"), r.Header.Get("Authorization"))
		assert.Equal(t, headers.Get("User-Agent"), r.Header.Get("User-Agent"))

		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)

		// simple test image
		img := image.NewRGBA(image.Rect(0, 0, 1, 1))
		img.Set(0, 0, color.RGBA{255, 255, 255, 255})

		jpeg.Encode(w, img, nil)
	}))

	defer server.Close()

	d := NewDownloader()

	imgURL := server.URL + "/image.jpg"
	result, err := d.Download(imgURL, headers)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestDownloader_ErrorCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/image.jpg", r.URL.Path)

		assert.Equal(t, headers.Get("Authorization"), r.Header.Get("Authorization"))
		assert.Equal(t, headers.Get("User-Agent"), r.Header.Get("User-Agent"))

		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("1"))
	}))
	defer server.Close()

	d := NewDownloader()

	imgURL := server.URL + "/image.jpg"
	result, err := d.Download(imgURL, headers)

	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid status: 400")
}

func TestDownloader_InvalidUrl(t *testing.T) {
	d := NewDownloader()

	imgURL := "invalid/image.jpg"
	result, err := d.Download(imgURL, headers)

	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "cant do request")
}

func TestDownloader_NotJpeg(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/image.jpg", r.URL.Path)

		assert.Equal(t, headers.Get("Authorization"), r.Header.Get("Authorization"))
		assert.Equal(t, headers.Get("User-Agent"), r.Header.Get("User-Agent"))

		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("1"))
	}))
	defer server.Close()

	d := NewDownloader()

	imgURL := server.URL + "/image.jpg"
	result, err := d.Download(imgURL, headers)

	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "cant decode jpeg image")
}

func TestDownloader_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/image.jpg", r.URL.Path)

		assert.Equal(t, headers.Get("Authorization"), r.Header.Get("Authorization"))
		assert.Equal(t, headers.Get("User-Agent"), r.Header.Get("User-Agent"))

		// Simulate a long response time
		time.Sleep(3 * time.Second)

		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("1"))
	}))
	defer server.Close()

	d := NewDownloader()

	imgURL := server.URL + "/image.jpg"
	result, err := d.Download(imgURL, headers)

	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "context deadline exceeded")
}
