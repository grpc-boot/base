package connctx

type linkedNode struct {
	next  *linkedNode
	value interface{}
}

func (ln *linkedNode) append(n *linkedNode) {
	ln.next = n
}

func (ln *linkedNode) prepend(n *linkedNode) {
	n.next = ln
}
