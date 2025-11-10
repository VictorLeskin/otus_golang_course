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
	Prev  *ListItem
	Next  *ListItem
}

type list struct {
	List  // Remove me after realization.
	front *ListItem
	back  *ListItem
	len   int
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l list) Len() int {
	return l.len
}

func NewListItem(v interface{}, prev *ListItem, next *ListItem) *ListItem {
	return &ListItem{Value: v, Prev: prev, Next: next}
}

func (l *list) Init(v interface{}) *ListItem {
	m := &ListItem{Value: v}
	l.front = m
	l.back = m
	l.len = 1
	return m
}

func (l *list) PushBack(v interface{}) *ListItem {
	if nil == l.back {
		return l.Init(v)
	}

	m := NewListItem(v, l.back, nil)

	l.back.Next = m
	l.back = m
	l.len++
	return m
}

func (l *list) PushFront(v interface{}) *ListItem {
	if nil == l.front {
		return l.Init(v)
	}

	m := NewListItem(v, nil, l.front)

	l.front.Prev = m
	l.front = m
	l.len++
	return m
}

func (l *list) RemoveFront() {
	if l.len == 1 {
		l.front = nil
		l.back = nil
	} else {
		l.front = l.front.Next
		l.front.Prev = nil
	}
	l.len--
}

func (l *list) RemoveBack() {
	if l.len == 1 {
		l.front = nil
		l.back = nil
	} else {
		l.back = l.back.Prev
		l.back.Next = nil
	}
	l.len--
}

func (l *list) Remove(j *ListItem) {
	switch {
	case l.front == j:
		l.RemoveFront()
	case l.back == j:
		l.RemoveBack()
	default: // item has prev and next
		j.Prev.Next = j.Next
		j.Next.Prev = j.Prev
		l.len--
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if l.front == i {
		return
	}

	l.Remove(i)

	i.Next = l.front
	i.Prev = nil

	l.front.Prev = i

	l.front = i

	l.len++
}

func NewList() List {
	return &list{
		front: nil,
		back:  nil,
		len:   0,
	}
}
