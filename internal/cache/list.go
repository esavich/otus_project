package cache

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
	frontItem *ListItem
	backItem  *ListItem
	len       int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.frontItem
}

func (l *list) Back() *ListItem {
	return l.backItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.len++
	item := &ListItem{
		Value: v,
		Next:  l.Front(),
	}
	if l.frontItem == nil {
		l.backItem = item
		l.frontItem = item

		return item
	}

	l.frontItem.Prev = item
	l.frontItem = item

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.len++
	item := &ListItem{
		Value: v,
		Prev:  l.Back(),
	}
	if l.backItem == nil {
		l.backItem = item
		l.frontItem = item

		return item
	}

	l.backItem.Next = item
	l.backItem = item

	return item
}

func (l *list) Remove(i *ListItem) {
	l.len--
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.frontItem = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.backItem = i.Prev
	}
}

func (l *list) MoveToFront(i *ListItem) {
	// already front
	if i.Prev == nil {
		return
	}

	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
