// Package list impl a doubly linked list.
//
// To iterate over a list (where l is a *List):
//  for e := l.Front(); e != nil; e = e.Next() {
//      // do something with e.Value
//  }
package list

type Element struct {
	Val        interface{}
	next, prev *Element

	list *List
}

func (e *Element) Next() *Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

func (e *Element) Prev() *Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// List a doubly linked list.
// The zero value for List is an empty list ready to use.
type List struct {
	root Element // sentinel element, only &root, root.prev, root.next are used
	len  int     // current list length excluding sentinel element
}

// lazyInit lazily initializes a zero List value.
func (l *List) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

func (l *List) insert(e, at *Element) *Element {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

func (l *List) insertVal(v interface{}, at *Element) *Element {
	return l.insert(&Element{Val: v}, at)
}

func (l *List) remove(e *Element) *Element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.prev = nil
	e.next = nil
	e.list = nil
	l.len--
	return e
}

func (l List) move(e, at *Element) *Element {
	if e == at {
		return e
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e

	return e
}

// Init initializes or clears list l.
func (l *List) Init() *List {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

func New() *List {
	return new(List).Init()
}

func (l *List) Len() int {
	return l.len
}

func (l *List) Front() *Element {
	if l.len == 0 {
		return nil
	}

	return l.root.next
}

func (l *List) Back() *Element {
	if l.len == 0 {
		return nil
	}

	return l.root.prev
}

func (l *List) Remove(e *Element) interface{} {
	if e.list == l {
		l.remove(e)
	}

	return e.Val
}

func (l *List) PushFront(v interface{}) *Element {
	l.lazyInit()
	return l.insertVal(v, &l.root)
}

func (l *List) PushBack(v interface{}) *Element {
	l.lazyInit()
	return l.insertVal(v, l.root.prev)
}

func (l *List) InsertBefore(v interface{}, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	return l.insertVal(v, mark.prev)
}

func (l *List) InsertAfter(v interface{}, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	return l.insertVal(v, mark)
}

func (l *List) MoveToFront(e *Element) {
	if e.list != l || l.root.next == e {
		return
	}
	l.move(e, &l.root)
}

func (l *List) MoveToBack(e *Element) {
	if e.list != l || l.root.prev == e {
		return
	}
	l.move(e, l.root.prev)
}

func (l *List) MoveBefore(e, mark *Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

func (l *List) MoveAfter(e, mark *Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

func (l *List) PushBackList(other *List) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertVal(e.Val, l.root.prev)
	}
}

func (l *List) PushFrountList(other *List) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertVal(e.Val, &l.root)
	}
}
