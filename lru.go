package lru

import "container/list"

type Element struct {
	_key  string
	list  *list.Element
	value interface{}
}

type Cache struct {
	capacity int
	_queue  *list.List
	_storage map[string]Element
}

func New(capacity int ) *Cache {
	return &Cache{
		capacity: capacity,
		_queue: list.New(),
		_storage: make(map[string]Element),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	v, ok := c._storage[key]
	if ok {
		c.reChange(key, v.list, value)
		return
	}

	if c._queue.Len() == c.capacity {
		c.purge()
	}

	element := Element{
		_key: key,
		value: value,
	}

	element.list = c._queue.PushBack(element)
	c._storage[key] = element
}

func (c *Cache) reChange(key string, l *list.Element, value interface{}) {
	c._queue.MoveToFront(l)

	c._storage[key] = Element{
		_key: key,
		list: l,
		value: value,
	}
}

func (c *Cache) purge() {
	element := c._queue.Back()
	c._queue.Remove(element)

	val, _ := element.Value.(Element)
	delete(c._storage, val._key)
}

func (c *Cache) Get(key string) interface{} {
	val, ok := c._storage[key]
	if !ok {
		return nil
	}

	c._queue.MoveToFront(val.list)
	return val.value
}