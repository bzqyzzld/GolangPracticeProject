package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes int64
	nBytes   int64
	ll       *list.List
	cache    map[string]*list.Element
	OnEvent  func(key string, value Value)
}

// New, 创建新的缓存
func New(maxBytes int64, onEvent func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
		OnEvent:  onEvent,
	}
}

// Get, 正确获取到数据之后会把元素移动到左边
func (cache *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := cache.cache[key]; ok {
		cache.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest
func (cache *Cache) RemoveOldest() {
	ele := cache.ll.Back()
	if ele != nil {
		cache.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(cache.cache, kv.key)
		cache.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if cache.OnEvent != nil {
			cache.OnEvent(kv.key, kv.value)
		}
	}
}

func (cache *Cache) Set(key string, value Value) {
	if ele, ok := cache.cache[key]; ok {
		cache.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		cache.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := cache.ll.PushFront(&entry{key, value})
		cache.cache[key] = ele
		cache.nBytes += int64(len(key)) + int64(value.Len())
	}
	for cache.nBytes > cache.maxBytes {
		cache.RemoveOldest()
	}
}

func (cache *Cache) Len() int {
	return cache.ll.Len()
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}
