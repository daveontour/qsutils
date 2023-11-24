// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package list implements a doubly linked list.
//
// To iterate over a list (where l is a *List):
//
//	for e := l.Front(); e != nil; e = e.Next() {
//		// do something with e.Value
//	}
package list

import "time"

// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.
type List struct {
	root Element // sentinel list element, only &root, root.prev, and root.next are used
	len  int     // current list length excluding (this) sentinel element
}

// Element is an element of a linked list.
type Element struct {
	next, prev *Element

	// The list to which this element belongs.
	list *List

	// The value stored with this element.
	Value any

	// Oredering options
	Priority int
	DateTime time.Time
}

// Next returns the next list element or nil.
func (e *Element) Next() *Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *Element) Prev() *Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Init initializes or clears list l.
func (l *List) Init() *List {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// New returns an initialized list.
func New() *List { return new(List).Init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *List) Len() int { return l.len }

// Front returns the first element of list l or nil if the list is empty.
func (l *List) Front() (*Element, bool) {
	if l.len == 0 {
		return nil, false
	}
	return l.root.next, true
}

// FrontPop returns the first element of list l or nil if the list is empty and then removes the element from the list
func (l *List) FrontPop() (*Element, bool) {
	if l.len == 0 {
		return nil, false
	}
	e := l.root.next
	l.remove(e)
	return e, true
}

// Back returns the last element of list l or nil if the list is empty.
func (l *List) Back() (*Element, bool) {
	if l.len == 0 {
		return nil, false
	}
	return l.root.prev, true
}

// Front returns the first element of list l or nil if the list is empty.
func (l *List) FrontOnly() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil if the list is empty.
func (l *List) BackOnly() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero List value.
func (l *List) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *List) insert(e, at *Element) *Element {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertBefore inserts e before at, increments l.len, and returns e.
func (l *List) insertBefore(e, at *Element) *Element {
	e.next = at
	e.prev = at.prev
	e.next.prev = e
	e.prev.next = e
	e.list = l
	l.len++
	return e
}

// SortByPriority sorts the list by priority
func (l *List) SortByPriority() {
	l.lazyInit()
	if l.len == 0 {
		return
	}
	for i := l.FrontOnly(); i != nil; i = i.Next() {
		for j := i.Next(); j != nil; j = j.Next() {
			if i.Priority < j.Priority {
				i.Value, j.Value = j.Value, i.Value
				i.Priority, j.Priority = j.Priority, i.Priority
				i.DateTime, j.DateTime = j.DateTime, i.DateTime
			}
		}
	}
}

// SortByDateTime sorts the list by DateTime
func (l *List) SortByDateTime() {
	l.lazyInit()
	if l.len == 0 {
		return
	}
	for i := l.FrontOnly(); i != nil; i = i.Next() {
		for j := i.Next(); j != nil; j = j.Next() {
			if i.DateTime.Before(j.DateTime) {
				i.Value, j.Value = j.Value, i.Value
				i.Priority, j.Priority = j.Priority, i.Priority
				i.DateTime, j.DateTime = j.DateTime, i.DateTime
			}
		}
	}
}

// InsertByPriority inserts e based on priority
func (l *List) InsertByPriority(v any, p int) *Element {
	l.lazyInit()
	e := &Element{Value: v, Priority: p}
	if l.len == 0 {
		return l.insert(e, &l.root)
	}
	for i, _ := l.Front(); i != nil; i = i.Next() {
		if i.Priority > e.Priority {
			return l.insertBefore(e, i)
		}
	}
	return l.insert(e, l.BackOnly())
}

// InsertByDateTime inserts e based on DateTime
func (l *List) InsertByDateTime(v any, t time.Time) *Element {
	l.lazyInit()
	e := &Element{Value: v, DateTime: t}
	if l.len == 0 {
		return l.insert(e, &l.root)
	}

	for i, _ := l.Front(); i != nil; i = i.Next() {
		if (!i.DateTime.IsZero() && i.DateTime.Before(e.DateTime)) || i.DateTime.IsZero() {
			return l.insertBefore(e, i)
		}
	}
	return l.insert(e, l.BackOnly())
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (l *List) insertValue(v any, at *Element) *Element {
	return l.insert(&Element{Value: v}, at)
}

// remove removes e from its list, decrements l.len
func (l *List) remove(e *Element) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
}

// move moves e to next to at.
func (l *List) move(e, at *Element) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (l *List) Remove(e *Element) any {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *List) PushFront(v any) *Element {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *List) PushBack(v any) *Element {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *List) InsertBefore(v any, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *List) InsertAfter(v any, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *List) MoveToFront(e *Element) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, &l.root)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *List) MoveToBack(e *Element) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, l.root.prev)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *List) MoveBefore(e, mark *Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *List) MoveAfter(e, mark *Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *List) PushBackList(other *List) {
	l.lazyInit()
	for i, e := other.Len(), other.FrontOnly(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of another list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *List) PushFrontList(other *List) {
	l.lazyInit()
	for i, e := other.Len(), other.BackOnly(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}
