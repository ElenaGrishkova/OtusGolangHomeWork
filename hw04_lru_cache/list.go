package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	// List // Remove me after realization.
	// Place your code here.
	length    int
	frontItem *ListItem
	backItem  *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.frontItem
}

func (l *list) Back() *ListItem {
	return l.backItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	if pushEmpty := pushIntoEmptyList(v, l); pushEmpty {
		return l.frontItem
	}
	oldFrontItem := l.Front()
	newFrontItem := &ListItem{v, oldFrontItem, nil}
	oldFrontItem.Prev = newFrontItem
	l.length++
	l.frontItem = newFrontItem
	return newFrontItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	if pushEmpty := pushIntoEmptyList(v, l); pushEmpty {
		return l.backItem
	}
	oldBackItem := l.Back()
	newBackItem := &ListItem{v, nil, oldBackItem}
	oldBackItem.Next = newBackItem
	l.length++
	l.backItem = newBackItem
	return newBackItem
}

func (l *list) Remove(i *ListItem) {
	switch {
	case l.length == 1:
		l.frontItem = nil
		l.backItem = nil
	case l.frontItem == i:
		l.frontItem = l.frontItem.Next
		l.frontItem.Prev = nil
	case l.backItem == i:
		l.backItem = l.backItem.Prev
		l.backItem.Next = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	i.Prev = nil
	i.Next = nil
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	switch {
	case l.frontItem == i:
		return
	case l.backItem == i:
		l.backItem = l.backItem.Prev
		l.backItem.Next = nil
	}
	// Предыдущий "первый элемент".
	oldFrontItem := l.Front()
	// Предыдущий "первый элемент" становится вторым. Обновим его ссылку на Prev.
	oldFrontItem.Prev = i

	// Общая характеристика списка "Первый элемент" обновляется.
	l.frontItem = i

	// Старые элементы, которые были раньше связаны с i теперь должны смотреть друг на друга.
	i.Prev.Next = i.Next
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	// У нового первого элемента обновлены ссылки.
	i.Next = oldFrontItem
	i.Prev = nil
}

func pushIntoEmptyList(v interface{}, l *list) bool {
	if l.length == 0 {
		singleItem := &ListItem{v, nil, nil}
		l.frontItem = singleItem
		l.backItem = singleItem
		l.length++
		return true
	}
	return false
}

func NewList() List {
	return new(list)
}
