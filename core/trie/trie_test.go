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
