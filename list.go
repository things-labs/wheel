package timewheel

type list struct {
	root Timer
	len  int // current list length excluding (this) sentinel element
}

// newList returns an initialized list.
func newList() *list { return new(list).Init() }

// Init initializes or clears list l.
func (l *list) Init() *list {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *list) Len() int { return l.len }

// lazyInit lazily initializes a zero list value.
func (l *list) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *list) insert(e, at *Timer) *Timer {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = l
	l.len++
	return e
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *list) remove(e *Timer) *Timer {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	return e
}

// Front returns the first element of list l or nil if the list is empty.
func (l *list) Front() *Timer {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// PushElementBack inserts a new element e at the back of list l and returns e.
func (l *list) PushElementBack(e *Timer) *Timer {
	l.lazyInit()
	return l.insert(e, l.root.prev)
}

// PopFront pop the first element of list l or nil if the list is empty.
func (l *list) PopFront() *Timer {
	l.lazyInit()
	if e := l.Front(); e != nil {
		return l.remove(e)
	}
	return nil
}

// SpliceBackList inserts an other list at the back of list l.
// and then remove all the other list element
// The lists l and other may be the same. They must not be nil.
func (l *list) SpliceBackList(other *list) {
	l.lazyInit()
	for other.len > 0 {
		l.PushElementBack(other.PopFront())
	}
}

// removeSelf remove self from list ,if it not on any list do nothing
func (t *Timer) removeSelf() {
	if t.list != nil {
		t.list.remove(t)
	}
}
