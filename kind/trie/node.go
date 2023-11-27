package trie

type node struct {
	value interface{}
	isEnd bool
	sub   map[rune]*node
}

func (n *node) set(key string, value any) (isNew bool) {
	if key == "" {
		return
	}

	current := n

	for _, r := range key {
		if val, exists := current.sub[r]; exists {
			current = val
			continue
		}

		isNew = true

		nd := &node{}
		if current.sub == nil {
			current.sub = make(map[rune]*node, 1)
		}

		current.sub[r] = nd
		current = nd
	}

	current.value = value
	current.isEnd = true
	return
}

func (n *node) get(key string) (value any, exists bool) {
	if key == "" {
		return
	}

	current := n

	for _, r := range key {
		nd, yes := current.sub[r]
		if !yes {
			return
		}

		current = nd
	}

	if current.isEnd {
		return current.value, true
	}
	return
}

func (n *node) exists(r rune) (exists bool) {
	_, exists = n.sub[r]
	return
}
