package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type t_list struct {
	list
}

func checkInvariancts(t *testing.T, t0 t_list) {
	if t0.len == 0 {
		require.Nil(t, t0.front)
		require.Nil(t, t0.back)
	} else {
		require.Nil(t, t0.front.Prev)
		require.Nil(t, t0.back.Next)

		current := t0.front
		var previous *ListItem = nil

		for current != nil {
			// Проверяем, что current->prev ведет на предыдущий узел, который мы запомнили
			require.Equal(t, current.Prev, previous)

			// Двигаемся вперед
			previous = current
			current = current.Next
		}

		// Если мы дошли до конца, previous должен быть равен tail
		require.Equal(t, previous, t0.back)
	}
}

func getElements(l List) []int {
	elems := make([]int, 0, l.Len())
	for i := l.Front(); i != nil; i = i.Next {
		elems = append(elems, i.Value.(int))
	}
	return elems
}

func Test_list_ctor(t *testing.T) {
	var t0 t_list

	assert.Equal(t, 0, t0.len)
	assert.Nil(t, t0.front)
	assert.Nil(t, t0.back)
}

func Test_list_Init(t *testing.T) {
	var t0 t_list

	t0.Init(222)

	assert.Equal(t, 1, t0.len)
	assert.Equal(t, t0.front, t0.back)

	item := t0.front
	assert.Nil(t, item.Next)
	assert.Nil(t, item.Prev)
	assert.Equal(t, 222, item.Value)
}

func TestList_Len(t *testing.T) {
	var t0 t_list

	assert.Equal(t, 0, t0.Len())
	t0.len = 22
	assert.Equal(t, 22, t0.Len())
}

func TestList_Front(t *testing.T) {
	var t0 t_list

	assert.Nil(t, t0.Front())
	var t1 ListItem
	t0.front = &t1
	assert.Equal(t, &t1, t0.Front())
}

func TestList_Back(t *testing.T) {
	var t0 t_list

	assert.Nil(t, t0.Back())
	var t1 ListItem
	t0.back = &t1
	assert.Equal(t, &t1, t0.Back())
}

func TestList_PushFront(t *testing.T) {
	var t0 t_list

	checkInvariancts(t, t0)
	assert.Nil(t, t0.Front())

	res := t0.PushFront(33)

	checkInvariancts(t, t0)
	assert.Nil(t, t0.front.Prev)
	assert.Nil(t, t0.front.Next)
	assert.Equal(t, 33, t0.front.Value)
	assert.Equal(t, res, t0.front)
	assert.Equal(t, 1, t0.len)

	res = t0.PushFront(44)

	checkInvariancts(t, t0)
	assert.NotNil(t, t0.front.Next)
	assert.Equal(t, 33, t0.front.Next.Value)
	assert.Equal(t, 44, t0.front.Value)
	assert.Equal(t, res, t0.front)
	assert.Equal(t, 2, t0.len)
}

func TestList_PushBack(t *testing.T) {
	var t0 t_list

	checkInvariancts(t, t0)
	assert.Nil(t, t0.Back())

	res := t0.PushBack(33)

	checkInvariancts(t, t0)
	assert.Nil(t, t0.back.Prev)
	assert.Nil(t, t0.back.Next)
	assert.Equal(t, 33, t0.back.Value)
	assert.Equal(t, res, t0.back)
	assert.Equal(t, 1, t0.len)

	res = t0.PushBack(44)

	checkInvariancts(t, t0)
	assert.NotNil(t, t0.back.Prev)
	assert.Equal(t, 33, t0.back.Prev.Value)
	assert.Equal(t, 44, t0.back.Value)
	assert.Equal(t, 2, t0.len)
	assert.Equal(t, res, t0.back)
}

func TestList_RemoveFront(t *testing.T) {
	var t0 t_list

	t0.PushBack(10) //
	t0.PushBack(20) //
	t0.PushBack(30) // [10, 20, 30]

	checkInvariancts(t, t0)
	assert.Equal(t, 10, t0.front.Value)
	assert.Equal(t, 3, t0.len)

	t0.RemoveFront()

	checkInvariancts(t, t0)
	assert.Equal(t, 20, t0.front.Value)
	assert.Equal(t, 2, t0.len)

	t0.RemoveFront()
	checkInvariancts(t, t0)
	assert.Equal(t, 30, t0.front.Value)
	assert.Equal(t, 1, t0.len)

	t0.RemoveFront()
	checkInvariancts(t, t0)
	assert.Equal(t, 0, t0.len)
}

func TestList_RemoveBack(t *testing.T) {
	var t0 t_list

	t0.PushBack(10) //
	t0.PushBack(20) //
	t0.PushBack(30) //
	t0.PushBack(40) //

	assert.Equal(t, 40, t0.back.Value)
	assert.Equal(t, 4, t0.len)

	checkInvariancts(t, t0)

	t0.RemoveBack()

	checkInvariancts(t, t0)
	assert.Equal(t, 30, t0.back.Value)
	assert.Equal(t, 3, t0.len)

	t0.RemoveBack()

	checkInvariancts(t, t0)
	assert.Equal(t, 20, t0.back.Value)
	assert.Equal(t, 2, t0.len)

	t0.RemoveBack()

	checkInvariancts(t, t0)
	assert.Equal(t, 10, t0.back.Value)
	assert.Equal(t, 1, t0.len)

	t0.RemoveBack()

	checkInvariancts(t, t0)
	assert.Equal(t, 0, t0.len)
}

func TestList_MoveToFront(t *testing.T) {
	var t0 t_list

	t0.PushBack(10) //

	checkInvariancts(t, t0)

	t0.MoveToFront(t0.front)

	checkInvariancts(t, t0)
	assert.Equal(t, 1, t0.len)
	assert.Equal(t, 10, t0.front.Value)

	t0.PushBack(20) //

	t0.MoveToFront(t0.front)
	checkInvariancts(t, t0)
	assert.Equal(t, 2, t0.len)
	assert.Equal(t, 10, t0.front.Value)
	assert.Equal(t, 20, t0.front.Next.Value)

	t0.MoveToFront(t0.front.Next)
	checkInvariancts(t, t0)
	assert.Equal(t, 2, t0.len)
	assert.Equal(t, 20, t0.front.Value)
	assert.Equal(t, 10, t0.front.Next.Value)

	t0.PushFront(30) //

	t0.MoveToFront(t0.front.Next)
	checkInvariancts(t, t0)
	assert.Equal(t, 3, t0.len)
	assert.Equal(t, 20, t0.front.Value)
	assert.Equal(t, 30, t0.front.Next.Value)
	assert.Equal(t, 10, t0.front.Next.Next.Value)
}

func TestList_Remove(t *testing.T) {
	var t0 t_list

	t10 := t0.PushBack(10) //
	t20 := t0.PushBack(20) //
	t30 := t0.PushBack(30) //
	t40 := t0.PushBack(40) //

	checkInvariancts(t, t0)

	t0.Remove(t20)
	checkInvariancts(t, t0)

	t0.Remove(t10)
	checkInvariancts(t, t0)

	t0.Remove(t40)
	checkInvariancts(t, t0)

	t0.Remove(t30)
	checkInvariancts(t, t0)
}

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
}
