package lru

import (
	"container/list"
	"sync"
)

type entry[K comparable, V any] struct {
	key   K
	value V
}

type Cache[K comparable, V any] struct {
	items    map[K]*list.Element
	capacity int
	evict    *list.List
	mu       sync.Mutex
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity <= 0 {
		panic("capacity must be greater than 0")
	}

	return &Cache[K, V]{
		items:    make(map[K]*list.Element, capacity),
		capacity: capacity,
		evict:    list.New(),
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		kv := elem.Value.(*entry[K, V])
		kv.value = value
		c.evict.MoveToFront(elem)
		return
	}

	if c.evict.Len() >= c.capacity {
		c.evictLRU()
	}

	elem := c.evict.PushFront(&entry[K, V]{key: key, value: value})
	c.items[key] = elem
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.evict.MoveToFront(elem)
		return elem.Value.(*entry[K, V]).value, true
	}

	var zero V
	return zero, false
}

func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.evict.Init()
	c.items = make(map[K]*list.Element, c.capacity)
}

func (c *Cache[K, V]) evictLRU() {
	lru := c.evict.Back()
	if lru == nil {
		return
	}

	kv := lru.Value.(*entry[K, V])
	delete(c.items, kv.key)
	c.evict.Remove(lru)
}
