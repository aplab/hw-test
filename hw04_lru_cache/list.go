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
	first *ListItem
	last  *ListItem
	len   int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := &ListItem{
		Value: v,
	}
	if l.len == 0 {
		l.last = i
	} else {
		i.Next = l.first
		l.first.Prev = i
	}
	l.first = i
	l.len++
	return i
}

func (l *list) PushBack(v interface{}) *ListItem {
	i := &ListItem{
		Value: v,
	}
	if l.len == 0 {
		l.first = i
	} else {
		i.Prev = l.last
		l.last.Next = i
	}
	l.last = i
	l.len++
	return i
}

func (l *list) Remove(i *ListItem) {
	if l.len == 0 {
		return
	}
	if i == nil {
		return
	}
	if i.Prev == nil {
		l.first = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	if i.Next == nil {
		l.last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}
	if i.Prev == nil {
		return
	}
	i.Prev.Next = i.Next
	if i.Next == nil {
		l.last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	i.Next = l.first
	i.Prev = nil
	l.first.Prev = i
	l.first = i
}
