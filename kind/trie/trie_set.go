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

func (ts *TrieSet) HasKey(words string) bool {
	var (
		cursor  int
		start   int
		current = ts.sub
		data    = []rune(words)
	)

	for start < len(data) {
		r := data[start]

		if !current.exists(r) {
			current = ts.sub
			// 一直没有找到
			if cursor == 0 {
				start++
				continue
			}

			// 有找到的开始
			start = cursor + 1
			cursor = 0
			continue
		}

		if !current.sub[r].isEnd {
			// 记录找到的位置
			if cursor == 0 {
				cursor = start
			}

			current = current.sub[r]
			start++
			continue
		}

		return true
	}

	return false
}

func (ts *TrieSet) ReplaceKey(words string, starChar byte) (newWords string) {
	var (
		cursor  int
		start   int
		current = ts.sub
		data    = []rune(words)
	)

	for start < len(data) {
		r := data[start]

		if !current.exists(r) {
			current = ts.sub
			// 一直没有找到
			if cursor == 0 {
				start++
				continue
			}

			// 有找到的开始
			start = cursor + 1
			cursor = 0
			continue
		}

		if !current.sub[r].isEnd {
			// 记录找到的位置
			if cursor == 0 {
				cursor = start
			}

			current = current.sub[r]
			start++
			continue
		}

		for ; cursor <= start; cursor++ {
			data[cursor] = rune(starChar)
		}

		cursor = 0
		current = ts.sub

		start++
	}

	return string(data)
}

func (ts *TrieSet) Length() int64 {
	return ts.length
}
