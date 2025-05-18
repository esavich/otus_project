package diskcache

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/esavich/otus_project/internal/cache"
)

type Wrapper struct {
	memCache cache.Cache
	basePath string
	mutex    sync.Mutex
}

func NewDiskCacheWrapper(capacity int, diskPath string) (*Wrapper, error) {
	err := os.MkdirAll(diskPath, 0o755)
	if err != nil {
		return nil, fmt.Errorf("can't create or open cache dir: %w", err)
	}

	wrapper := &Wrapper{
		memCache: cache.NewCache(capacity),
		basePath: diskPath,
	}

	err = wrapper.ClearDiskCache()
	if err != nil {
		return nil, err
	}

	return wrapper, nil
}

func (dc *Wrapper) Set(key string, data image.Image) error {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	filePath := dc.getFilePath(key)

	outFile, err := os.Create(filePath)
	if err != nil {
		slog.Error(err.Error())
	}
	defer outFile.Close()
	err = jpeg.Encode(outFile, data, nil)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	// func to remove file from disk
	removeCallback := func(value interface{}) {
		fileToDeletePath := value.(string)
		slog.Error(fmt.Sprintf("Removing file: %s", fileToDeletePath))
		err := os.Remove(fileToDeletePath)
		if err != nil {
			slog.Error(fmt.Sprintf("Can't remove file %s: %s", fileToDeletePath, err))
		}
	}

	dc.memCache.Set(cache.Key(key), filePath, removeCallback)

	return nil
}

func (dc *Wrapper) Get(key string) (image.Image, bool) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	cachedPath, found := dc.memCache.Get(cache.Key(key))
	if !found {
		return nil, false
	}

	filePath, ok := cachedPath.(string)
	if !ok {
		slog.Error("Cache value is not a string")
		return nil, false
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error(fmt.Sprintf("Can't read file %s from disk: %s", filePath, err))
		return nil, false
	}
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		slog.Error(fmt.Sprintf("Can't decode jpeg: %s", err))
		return nil, false
	}
	return img, true
}

func (dc *Wrapper) getFilePath(key string) string {
	// hash name to avoid long names and special symbols compatibility problems
	h := sha256.New()
	h.Write([]byte(key))
	hash := hex.EncodeToString(h.Sum(nil))

	return filepath.Join(dc.basePath, fmt.Sprintf("%s.jpg", hash))
}

func (dc *Wrapper) ClearDiskCache() error {
	d, err := os.Open(dc.basePath)
	if err != nil {
		return fmt.Errorf("can't open cache dir: %w", err)
	}
	defer d.Close()

	files, err := d.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("can't read cache dir: %w", err)
	}

	for _, file := range files {
		filePath := filepath.Join(dc.basePath, file)
		err = os.RemoveAll(filePath)
		if err != nil {
			return fmt.Errorf("can't remove %s: %w", filePath, err)
		}
	}
	return nil
}
