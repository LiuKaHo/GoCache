package GoCache

import (
	"container/list"
	"fmt"
)

type Cache struct {
	maxBytes int64
	useBytes int64
	ll       *list.List
	cache    map[string]*list.Element
	OnDelete func(key string, value Value)
}

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

func (cache *Cache) Insert(key string, value Value) error {
	if cache.maxBytes < int64(value.Len()) {
		return fmt.Errorf("out of memory")
	}
	//key exists
	if ele, ok := cache.cache[key]; ok {
		cache.useBytes = cache.useBytes - int64(ele.Value.(*entry).value.Len()) + int64(value.Len())
		ele.Value = &entry{
			key:   key,
			value: value,
		}

		cache.ll.MoveToFront(ele)
	} else { //not exists
		ele := cache.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		cache.useBytes += int64(value.Len())
		cache.cache[key] = ele
	}

	for cache.maxBytes < cache.useBytes {
		cache.Remove()
	}

	return nil
}

func (cache *Cache) Remove() {
	last := cache.ll.Back()
	if last != nil {
		cache.ll.Remove(last)
		lastEntry := last.Value.(*entry)
		cache.useBytes -= int64(lastEntry.value.Len())
		delete(cache.cache, lastEntry.key)
		if cache.OnDelete != nil {
			cache.OnDelete(lastEntry.key, lastEntry.value)
		}
	}

}
