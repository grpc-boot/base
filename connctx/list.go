package connctx

type list struct {
	head   *linkedNode
	tail   *linkedNode
	length int
}

func (l *list) reset() {
	l.head = nil
	l.tail = nil
	l.length = 0
}

func (l *list) prepend(items ...any) {
	if len(items) < 1 {
		return
	}

	head := &linkedNode{value: items[0]}
	tail := head

	for index := 1; index < len(items); index++ {
		current := &linkedNode{value: items[index]}
		head.prepend(current)
		head = current
	}

	if l.head == nil {
		l.tail = tail
	} else {
		tail.append(l.head)
	}

	l.head = head
	l.length += len(items)
}

func (l *list) lpop() (value any) {
	if l.length < 1 {
		return
	}

	value = l.head.value
	l.head = l.head.next
	l.length--

	if l.head == nil {
		l.reset()
	}
	return
}

func (l *list) append(items ...any) {
	if len(items) < 1 {
		return
	}

	head := &linkedNode{value: items[0]}
	tail := head

	for index := 1; index < len(items); index++ {
		current := &linkedNode{value: items[index]}
		tail.append(current)
		tail = current
	}

	if l.head == nil {
		l.head = head
	} else {
		l.tail.append(head)
	}

	l.tail = tail
	l.length += len(items)
}

func (l *list) index(index int) (value any) {
	if index >= l.length || l.length < 1 {
		return
	}

	rIndex, err := realIndex(index, l.length)
	if err != nil {
		return
	}

	if rIndex == l.length-1 {
		value = l.tail.value
		return
	}

	current := l.head
	for rIndex != 0 {
		current = current.next
		rIndex--
	}
	value = current.value
	return
}

func (l *list) lrange(start, end int) (valueList []any, err error) {
	if l.length < 1 {
		return
	}

	startIndex, er := realIndex(start, l.length)
	if er != nil {
		return
	}

	endIndex, er := realIndex(end, l.length)
	if er != nil {
		endIndex = l.length - 1
	} else if endIndex > l.length-1 {
		endIndex = l.length - 1
	}

	if startIndex > endIndex {
		return
	}

	itemCount := endIndex - startIndex + 1
	valueList = make([]any, itemCount, itemCount)

	current := l.head
	for startIndex != 0 {
		current = current.next
		startIndex--
	}

	for i := 0; i < itemCount; i++ {
		valueList[i] = current.value
		current = current.next
	}

	return
}

func (l *list) trim(start, end int) {
	if l.length < 1 {
		return
	}

	if start > l.length {
		l.reset()
		return
	}

	startIndex, er := realIndex(start, l.length)
	if er != nil {
		startIndex = 0
	}

	endIndex, er := realIndex(end, l.length)
	if er != nil {
		endIndex = l.length - 1
	} else if endIndex > l.length-1 {
		endIndex = l.length - 1
	}

	if startIndex > endIndex {
		l.reset()
		return
	}

	if startIndex == 0 && endIndex == l.length-1 {
		return
	}

	itemCount := endIndex - startIndex + 1

	current := l.head
	for startIndex != 0 {
		next := current.next
		current.next = nil
		current = next
		startIndex--
	}

	l.head = current
	for i := 0; i < itemCount-1; i++ {
		current = current.next
	}
	current.next = nil
	l.tail = current

	l.length = itemCount
	return
}

func (l *list) set(index int, value any) (err error) {
	if index >= l.length || l.length < 1 {
		return ErrIndexOutOfRange
	}

	rIndex, err := realIndex(index, l.length)
	if err != nil {
		return err
	}

	if rIndex == l.length-1 {
		l.tail.value = value
		return
	}

	current := l.head
	for rIndex != 0 {
		current = current.next
		rIndex--
	}
	current.value = value
	return
}
