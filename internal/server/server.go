package server

import (
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/esavich/otus_project/internal/config"
	"github.com/esavich/otus_project/internal/handlers/resize"
	"github.com/esavich/otus_project/internal/service"
)

type Server struct {
	Config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Config: cfg,
	}
}

func (s *Server) Start() error {
	addr := net.JoinHostPort(s.Config.HTTP.Host, strconv.Itoa(s.Config.HTTP.Port))

	mux := http.NewServeMux()

	imageService := service.NewSimpleImageService()
	cachedService := service.NewCachedImageService(imageService)

	rh := resize.NewResizeHandler(cachedService)
	mux.HandleFunc("GET /fill/{width}/{height}/{url...}", rh.Resize)

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	slog.Info("Starting server at : " + addr)
	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
