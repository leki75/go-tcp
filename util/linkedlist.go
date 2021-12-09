package util

type item struct {
	data []byte
	next *item
}

type LinkedList struct {
	head *item
	tail *item
	size int
}

func (l *LinkedList) Add(b []byte) {
	l.size++
	item := &item{data: b}
	if l.head == nil {
		l.head = item
		l.tail = item
	} else {
		l.tail.next = item
		l.tail = item
	}
}

func (l *LinkedList) Remove() ([]byte, bool) {
	if l.head == nil {
		return []byte{}, false
	}

	l.size--
	data := l.head.data
	l.head = l.head.next
	return data, true
}

func (l *LinkedList) Size() int {
	return l.size
}
