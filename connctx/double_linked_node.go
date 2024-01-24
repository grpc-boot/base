package connctx

type doubleLinkednode struct {
	prev  *doubleLinkednode
	next  *doubleLinkednode
	value any
}

func (dln *doubleLinkednode) remove() (value any) {
	value = dln.value

	if dln.next != nil {
		dln.next.prev = nil
	}

	if dln.prev != nil {
		dln.prev.next = nil
	}

	dln.prev = nil
	dln.next = nil
	dln.value = nil
	return
}

func (dln *doubleLinkednode) append(a *doubleLinkednode) {
	dln.next = a
	a.prev = dln
}

func (dln *doubleLinkednode) prepend(p *doubleLinkednode) {
	dln.prev = p
	p.next = dln
}
