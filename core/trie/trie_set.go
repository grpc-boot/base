package trie

var (
	setValue = struct{}{}
)

type TrieSet struct {
	sub    *node
	length int64
}

func NewTrieSet() *TrieSet {
	return &TrieSet{
		sub: &node{},
	}
}

func (ts *TrieSet) Add(key string) (isNew bool) {
	isNew = ts.sub.set(key, setValue)
	if isNew {
		ts.length++
	}

	return
}

func (ts *TrieSet) Exists(key string) (exists bool) {
	_, exists = ts.sub.get(key)
	return
}

func (ts *TrieSet) Length() int64 {
	return ts.length
}
