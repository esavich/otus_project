//go:build integration

package tests

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/esavich/otus_project/internal/diskcache"
	"github.com/esavich/otus_project/internal/downloader"
	"github.com/esavich/otus_project/internal/resizer"
	"github.com/esavich/otus_project/internal/service"
)

// from https://golang.testcontainers.org/examples/nginx/

type nginxContainer struct {
	testcontainers.Container
	URI string
}

func startContainer(ctx context.Context) (*nginxContainer, error) {
	imagesDirectory, err := filepath.Abs(filepath.Join(".", "examples"))

	req := testcontainers.ContainerRequest{
		Image:        "nginx:latest",
		ExposedPorts: []string{"80/tcp"},
		WaitingFor:   wait.ForHTTP("/").WithStartupTimeout(10 * time.Second),
		// https://golang.testcontainers.org/features/files_and_mounts/#copying-directories-to-a-container
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      imagesDirectory,
				ContainerFilePath: "/usr/share/nginx/html/",
				FileMode:          0o777,
			},
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	var nginxC *nginxContainer
	if container != nil {
		nginxC = &nginxContainer{Container: container}
	}
	if err != nil {
		return nginxC, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nginxC, err
	}

	mappedPort, err := container.MappedPort(ctx, "80")
	if err != nil {
		return nginxC, err
	}

	nginxC.URI = fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())
	return nginxC, nil
}

func TestCachedService(t *testing.T) {
	ctx := context.Background()
	nginxC, err := startContainer(ctx)
	require.NoError(t, err)

	fmt.Println("Nginx container started at:", nginxC.URI)

	imageService := service.NewSimpleImageService(
		downloader.NewDownloader(5*time.Second),
		resizer.NewResizer(),
	)
	dir := t.TempDir()
	dc, err := diskcache.NewDiskCacheWrapper(3, dir)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating disk cache: %s", err))
		return
	}
	cachedService := service.NewCachedImageService(imageService, dc)
	headers := http.Header{
		"Authorization": []string{"Token mock-token"},
		"User-Agent":    []string{"Mozilla/5.0"},
	}

	t.Run("invalid url", func(t *testing.T) {
		imgURL := "invalid"
		result, err := cachedService.GetResizedImage(50, 60, imgURL, headers)

		require.Error(t, err)
		require.ErrorContains(t, err, "cant do request:")
		require.Nil(t, result)
	})

	t.Run("success", func(t *testing.T) {
		imgURL := nginxC.URI + "/examples/gopher.jpg"
		result, err := cachedService.GetResizedImage(50, 60, imgURL, headers)

		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("404", func(t *testing.T) {
		imgURL := nginxC.URI + "/examples/gopher-404.jpg"
		result, err := cachedService.GetResizedImage(50, 60, imgURL, headers)

		require.Error(t, err)
		require.ErrorContains(t, err, "invalid status: 404 Not Found")
		require.Nil(t, result)
	})

	t.Run("not image", func(t *testing.T) {
		imgURL := nginxC.URI + "/examples/1.txt"
		result, err := cachedService.GetResizedImage(50, 60, imgURL, headers)

		require.Error(t, err)
		require.ErrorContains(t, err, "cant decode jpeg")
		require.Nil(t, result)
	})

	t.Run("broken image", func(t *testing.T) {
		imgURL := nginxC.URI + "/examples/bad.jpg"
		result, err := cachedService.GetResizedImage(50, 60, imgURL, headers)

		require.Error(t, err)
		require.ErrorContains(t, err, "cant decode jpeg")
		require.Nil(t, result)
	})

	t.Run("check cache by log", func(t *testing.T) {
		imgURL := nginxC.URI + "/examples/gopher.jpg"

		// fake logger for test
		var logBuf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&logBuf, nil))
		oldLogger := slog.Default()
		slog.SetDefault(logger)
		// revert logger after test
		defer slog.SetDefault(oldLogger)

		// cache miss
		_, err := cachedService.GetResizedImage(220, 300, imgURL, headers)
		require.NoError(t, err)
		require.Contains(t, logBuf.String(), "downloading")

		for i := 0; i < 10; i++ {

			// clear log
			logBuf.Reset()

			// must be from cache
			_, err = cachedService.GetResizedImage(220, 300, imgURL, headers)
			require.NoError(t, err)

			// check cache hit
			require.Contains(t, logBuf.String(), "Cache hit")
		}
	})

	testcontainers.CleanupContainer(t, nginxC)
	require.NoError(t, err)
}
