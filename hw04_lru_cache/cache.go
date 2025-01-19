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
	sync.Mutex
	queue List
	items map[Key]*ListItem
}

// возвращаемое значение - флаг, присутствовал ли элемент в кэше.
func (l *lruCache) Set(key Key, value interface{}) bool {
	l.Lock()
	defer l.Unlock()
	if oldItem, ok := l.items[key]; ok {
		// если элемент присутствует в словаре, то обновить его значение и переместить элемент в начало очереди
		oldItem.Value = value
		return true
	}

	// если элемента нет в словаре, то добавить в словарь и в начало очереди. Делаем в два этапа:

	// Сперва: если размер очереди больше ёмкости кэша, то необходимо удалить последний элемент из очереди
	// и его значение из словаря
	if l.capacity == l.queue.Len() {
		backItem := l.queue.Back()
		l.queue.Remove(backItem)

		for keyToDel, valueToDel := range l.items {
			if valueToDel == backItem {
				delete(l.items, keyToDel)
			}
		}
	}
	// Затем: собственно добавление в словарь и в начало очереди.
	l.items[key] = l.queue.PushFront(value)

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.Lock()
	defer l.Unlock()
	if foundItem, ok := l.items[key]; ok {
		// если элемент присутствует в словаре, то переместить элемент в начало очереди и вернуть его значение и true
		l.queue.MoveToFront(foundItem)
		return foundItem.Value, true
	}
	// если элемента нет в словаре, то вернуть nil и false (работа с кешом похожа на работу с map)

	return nil, false
}

func (l *lruCache) Clear() {
	l.Lock()
	defer l.Unlock()
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
