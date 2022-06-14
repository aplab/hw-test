package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	sync.RWMutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	_, ok := l.items[key]
	l.Lock()
	defer l.Unlock()
	l.items[key] = l.queue.PushFront(cacheItem{
		key:   key,
		value: value,
	})
	if l.queue.Len() > l.capacity {
		b := l.queue.Back()
		l.queue.Remove(b)
		delete(l.items, b.Value.(cacheItem).key)
	}
	return ok
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.RLock()
	defer l.RUnlock()
	if i, ok := l.items[key]; ok {
		l.queue.MoveToFront(i)
		return i.Value.(cacheItem).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.Lock()
	defer l.Unlock()
	l.items = make(map[Key]*ListItem, l.capacity)
}
