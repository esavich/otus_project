package cache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}, callback func(value interface{})) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    sync.Mutex
}

func (l *lruCache) Set(key Key, value interface{}, callback func(value interface{})) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	item, isInCache := l.items[key]
	if isInCache {
		ci := item.Value.(cacheItem)
		ci.value = value
		item.Value = ci
		l.queue.MoveToFront(item)

		return true
	}

	ci := cacheItem{
		key:   key,
		value: value,
	}
	if l.queue.Len() == l.capacity {
		lastItem := l.queue.Back()
		if callback != nil {
			callback(lastItem.Value.(cacheItem).value)
		}
		l.queue.Remove(lastItem)
		delete(l.items, lastItem.Value.(cacheItem).key)
	}
	newItem := l.queue.PushFront(ci)
	l.items[key] = newItem

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	item, isInCache := l.items[key]
	if isInCache {
		l.queue.MoveToFront(item)
		ci := item.Value.(cacheItem)

		return ci.value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mutex:    sync.Mutex{},
	}
}
