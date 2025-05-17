package resize

import (
	"fmt"
	"image/jpeg"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/esavich/otus_project/internal/downloader"
	"github.com/esavich/otus_project/internal/resizer"
)

type Handler struct{}

func NewResizeHandler() *Handler {
	return &Handler{}
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

	d := downloader.NewDownloader()
	img, err := d.Download(imgURL, r.Header)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Cant get image: "+err.Error(), http.StatusBadGateway)
		return
	}
	slog.Info("Image downloaded")
	slog.Info("Resizing image", slog.Int("width", iw), slog.Int("height", ih))
	resized, _ := resizer.ResizeImg(img, iw, ih)

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
	// костыль тк дефолтный роутер переписывает двойные слеши в один
	// так же добавим дефолтный http префикс если его нет

	switch {
	case strings.HasPrefix(u, "http:/") && !strings.HasPrefix(u, "http://"):
		u = strings.Replace(u, "http:/", "http://", 1)
	case strings.HasPrefix(u, "https:/") && !strings.HasPrefix(u, "https://"):
		u = strings.Replace(u, "https:/", "https://", 1)
	case !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://"):
		u = "http://" + u
	}

	// провалидируем что получился корректный урл
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
