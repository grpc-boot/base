package trie

type Trie struct {
	sub    *node
	length int64
}

func New() *Trie {
	return &Trie{
		sub: &node{},
	}
}

func (t *Trie) Set(key string, value any) (isNew bool) {
	isNew = t.sub.set(key, value)
	if isNew {
		t.length++
	}

	return
}

func (t *Trie) Get(key string) (value any, exists bool) {
	return t.sub.get(key)
}

func (t *Trie) Length() int64 {
	return t.length
}
