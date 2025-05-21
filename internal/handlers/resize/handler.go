package resize

import (
	"fmt"
	"image"
	"image/jpeg"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ImageGetter interface {
	GetResizedImage(width, height int, imgURL string, header http.Header) (image.Image, error)
}

type Handler struct {
	ig ImageGetter
}

func NewResizeHandler(ig ImageGetter) *Handler {
	return &Handler{
		ig: ig,
	}
}

func (h *Handler) Resize(w http.ResponseWriter, r *http.Request) {
	width := r.PathValue("width")
	height := r.PathValue("height")
	imgURL := r.PathValue("url")

	iw, err := convertDimension(width)
	if err != nil {
		http.Error(w, "Invalid width parameter: "+width, http.StatusBadRequest)
		return
	}
	ih, err := convertDimension(height)
	if err != nil {
		http.Error(w, "Invalid height parameter: "+height, http.StatusBadRequest)
		return
	}

	imgURL, err = processURL(imgURL)
	if err != nil {
		http.Error(w, "Invalid URL parameter: "+err.Error(), http.StatusBadRequest)
		return
	}

	if !checkJpg(imgURL) {
		http.Error(w, "Invalid URL parameter: not jpeg", http.StatusBadRequest)
		return
	}
	slog.Info("Params", slog.Int("width", iw), slog.Int("height", ih), slog.String("url", imgURL))

	resized, err := h.ig.GetResizedImage(iw, ih, imgURL, r.Header)
	if err != nil {
		http.Error(w, "Cant get image: "+err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(w, resized, nil)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Cant encode image: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func checkJpg(imgURL string) bool {
	lower := strings.ToLower(imgURL)
	return strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg")
}

func processURL(u string) (string, error) {
	// workaround because the default router rewrites double slashes to a single one
	// also add the default http prefix if it is missing

	switch {
	case strings.HasPrefix(u, "http:/") && !strings.HasPrefix(u, "http://"):
		u = strings.Replace(u, "http:/", "http://", 1)
	case strings.HasPrefix(u, "https:/") && !strings.HasPrefix(u, "https://"):
		u = strings.Replace(u, "https:/", "https://", 1)
	case !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://"):
		u = "http://" + u
	}

	// validate url
	_, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	return u, nil
}

func convertDimension(dimension string) (int, error) {
	value, err := strconv.Atoi(dimension)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("invalid value: %s", dimension)
	}
	return value, nil
}
