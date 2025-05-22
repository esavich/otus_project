package diskcache

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{255, 255, 255, 255})

	return img
}

func TestSetAndGet(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewDiskCacheWrapper(2, dir)
	require.NoError(t, err)

	img := createTestImage()
	key := "test-key"
	err = cache.Set(key, img)
	require.NoError(t, err)

	gotImg, ok := cache.Get(key)
	require.True(t, ok)
	require.NotNil(t, gotImg)
}

func TestGetNotFound(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewDiskCacheWrapper(2, dir)
	require.NoError(t, err)

	_, ok := cache.Get("not-exist")
	require.False(t, ok)
}

func TestClearDiskCache(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewDiskCacheWrapper(2, dir)
	require.NoError(t, err)

	img := createTestImage()
	key := "clear-key"
	err = cache.Set(key, img)
	require.NoError(t, err)

	filePath := cache.getFilePath(key)
	_, err = os.Stat(filePath)
	require.NoError(t, err)

	err = cache.ClearDiskCache()
	require.NoError(t, err)

	_, err = os.Stat(filePath)
	require.Error(t, err)
}

func TestGetFilePath(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewDiskCacheWrapper(2, dir)
	require.NoError(t, err)

	key := "some-key"
	path := cache.getFilePath(key)
	require.Equal(t, dir, filepath.Dir(path))
	require.Equal(t, ".jpg", filepath.Ext(path))
}

func TestSetInvalidPath(t *testing.T) {
	cache := &Wrapper{
		basePath: "/invalid/path/for/test",
	}
	img := createTestImage()
	err := cache.Set("key", img)
	require.Error(t, err)
}

func TestSetAndOldCacheRemoving(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewDiskCacheWrapper(1, dir)
	require.NoError(t, err)

	img1 := createTestImage()
	img2 := createTestImage()
	err = cache.Set("key1", img1)
	require.NoError(t, err)
	err = cache.Set("key2", img2)
	require.NoError(t, err)

	_, ok := cache.Get("key1")
	require.False(t, ok)
	_, ok = cache.Get("key2")
	require.True(t, ok)
}
