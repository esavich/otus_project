package cache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100, nil)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200, nil)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300, nil)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5)
		c.Set("aaa", 100, nil)
		c.Set("bbb", 200, nil)

		c.Clear()

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("capacity logic", func(t *testing.T) {
		c := NewCache(3)
		c.Set("aaa", 100, nil)
		c.Set("bbb", 200, nil)
		c.Set("ccc", 300, nil)
		c.Set("ddd", 400, nil)

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("LRU logic", func(t *testing.T) {
		c := NewCache(3)
		c.Set("aaa", 100, nil)
		c.Set("bbb", 200, nil)
		c.Set("ccc", 300, nil)

		c.Get("aaa")           // aaa used
		c.Set("bbb", 500, nil) // bbb updated
		// ccc not used, will be removed

		c.Set("ddd", 400, nil)

		// aaa exists
		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		// bbb exists
		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 500, val)
		// ddd exists
		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400, val)

		// ccc not exists
		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i, nil)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000)))) //nolint:gosec
		}
	}()

	wg.Wait()
}

func TestCacheCallback(t *testing.T) {
	c := NewCache(2)

	var callbackCalled bool
	callback := func(value interface{}) {
		callbackCalled = true
		require.Equal(t, 100, value) // check that the right removing value is passed
	}

	// Заполняем кеш до предела
	c.Set("aaa", 100, nil)
	c.Set("bbb", 200, nil)

	// add new element to remove 100 value
	c.Set("ccc", 300, callback)

	require.True(t, callbackCalled, "Callback must be called")
}
