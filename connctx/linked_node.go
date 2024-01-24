package connctx

type linkedNode struct {
	next  *linkedNode
	value any
}

func (ln *linkedNode) append(n *linkedNode) {
	ln.next = n
}

func (ln *linkedNode) prepend(n *linkedNode) {
	n.next = ln
}
