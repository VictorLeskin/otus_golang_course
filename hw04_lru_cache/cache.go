package hw04lrucache

type Key string
type cacheItem struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (r *lruCache) Set(key Key, value interface{}) bool {
	// check the key
	if item, exist := r.items[key]; exist {
		// there is such key -> move its forward and update its value
		r.queue.MoveToFront(item)
		item.Value = cacheItem{key: key, value: value}
		return true
	}

	// check if the cache is full
	if r.queue.Len() == r.capacity {
		// remove oldest
		b := r.queue.Back()
		r.queue.Remove(b)
		delete(r.items, b.Value.(cacheItem).key)
	}

	// add new element to head of list and to the map
	item := r.queue.PushFront(cacheItem{key: key, value: value})
	r.items[key] = item

	return false
}

func (r *lruCache) Get(key Key) (interface{}, bool) {
	if item, exist := r.items[key]; exist {
		// there is such key -> move its forward and update its value
		r.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, true
	}
	return nil, false
}

func (r *lruCache) Clear() {
	for b := r.queue.Back(); b != nil; b = b.Prev {
		r.queue.Remove(b)
		delete(r.items, b.Value.(cacheItem).key)
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
