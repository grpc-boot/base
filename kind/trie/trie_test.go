package trie

import "testing"

func TestTrie_Get(t *testing.T) {
	te := New()

	te.Set("我", "好的")
	t.Logf("length:%d", te.Length())

	te.Set("你好", "好的啊")
	te.Set("你好", "好的啊")
	t.Logf("length:%d", te.Length())
	te.Set("他也好", "真好的")
	te.Set("他也好", "真好的啊")
	t.Logf("length:%d", te.Length())

	value, exists := te.Get("我")
	t.Logf("value:%v exists: %t", value, exists)

	if !exists {
		t.Fatalf("want true, got %t", exists)
	}

	value, exists = te.Get("你好")
	t.Logf("value:%v exists: %t", value, exists)

	if !exists {
		t.Fatalf("want true, got %t", exists)
	}

	value, exists = te.Get("他也好")
	t.Logf("value:%v exists: %t", value, exists)

	if !exists {
		t.Fatalf("want true, got %t", exists)
	}
}

func words() []string {
	return []string{
		`我`,
		`中国人`,
		`很好`,
	}
}

func TestNewTrieSet(t *testing.T) {
	ts := NewTrieSet()

	wordList := words()

	for _, word := range wordList {
		ts.Add(word)
	}

	res := ts.HasKey(`当时有好多中国人`)
	if !res {
		t.Fatalf("want true, got %v", res)
	}

	res = ts.HasKey(`当时有好多很好`)
	if !res {
		t.Fatalf("want true, got %v", res)
	}

	res = ts.HasKey(`当时有好多很`)
	if res {
		t.Fatalf("want false, got %v", res)
	}

	rd := ts.ReplaceKey(`当时我看见有好多中国人`, '*')
	t.Logf("rd: %s", rd)
}
