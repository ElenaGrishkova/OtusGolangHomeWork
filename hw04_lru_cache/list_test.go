package hw04lrucache

import (
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("customTest", func(t *testing.T) {
		l := NewList()

		// Добавление в пустой в начало
		l.PushFront(10) // [10]
		require.Equal(t, 1, l.Len())
		frontItem := l.Front()
		require.Equal(t, 10, frontItem.Value)
		require.Nil(t, frontItem.Next)
		require.Nil(t, frontItem.Prev)

		// Делаем снова пустым
		l.Remove(frontItem) // []
		require.Equal(t, 0, l.Len())
		require.Nil(t, frontItem.Next)
		require.Nil(t, frontItem.Prev)

		// Добавление в пустой в конец
		l.PushBack(20) // [20]
		require.Equal(t, 1, l.Len())
		backItem := l.Back()
		require.Equal(t, 20, backItem.Value)
		require.Nil(t, backItem.Next)
		require.Nil(t, backItem.Prev)

		// Делаем снова пустым
		l.Remove(backItem) // []
		require.Equal(t, 0, l.Len())
		require.Nil(t, backItem.Next)
		require.Nil(t, backItem.Prev)

		// Проверяем корректность связей
		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]

		frontItem = l.Front()
		middleItem := l.Front().Next
		backItem = l.Back()
		checkValuesOfSize3(t, frontItem, middleItem, backItem, 10, 20, 30)
		checkLinks(t, middleItem, frontItem, backItem)

		// Проверяем корректность связей MoveToFront
		l.MoveToFront(backItem)
		frontItem = l.Front()
		middleItem = l.Front().Next
		backItem = l.Back()
		checkValuesOfSize3(t, frontItem, middleItem, backItem, 30, 10, 20)
		checkLinks(t, middleItem, frontItem, backItem)
	})
}

func checkValuesOfSize3(t *testing.T, frontItem *ListItem, middleItem *ListItem, backItem *ListItem,
	//nolint:gofumpt
	expected1 interface{}, expected2 interface{}, expected3 interface{}) {
	t.Helper()
	require.Equal(t, expected1, frontItem.Value)
	require.Equal(t, expected2, middleItem.Value)
	require.Equal(t, expected3, backItem.Value)
}

func checkLinks(t *testing.T, middleItem *ListItem, frontItem *ListItem, backItem *ListItem) {
	t.Helper()
	require.Equal(t, middleItem, frontItem.Next)
	require.Nil(t, frontItem.Prev)
	require.Equal(t, backItem, middleItem.Next)
	require.Equal(t, frontItem, middleItem.Prev)
	require.Nil(t, backItem.Next)
	require.Equal(t, middleItem, backItem.Prev)
}
