package components

import (
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

func TestCode2Id(t *testing.T) {
	var (
		num = 50
	)

	rand.Seed(uint64(time.Now().Unix()))

	t.Logf("max6:%d max8:%d", Max6(), Max8())

	for i := 1; i < num; i++ {
		id := int64(rand.Uint32())

		code6, _ := Code6(id)
		decode6Id, _ := Code2Id(code6)

		if id != decode6Id {
			t.Fatalf("id:%d code:%s decodeId:%d", id, code6, decode6Id)
		}

		code8, _ := Code8(id)
		decode8Id, _ := Code2Id(code8)

		if id != decode8Id {
			t.Fatalf("id:%d code:%s decodeId:%d", id, code8, decode8Id)
		}

		t.Logf("id: %d 6:%s 8:%s", id, code6, code8)
	}
}

// BenchmarkId2Code6-8   	 7852514	       146.3 ns/op
func BenchmarkId2Code6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := int64(rand.Uint32())
		code, err := Code6(id)

		if err != nil {
			b.Fatalf("id:%d code:%s", id, code)
		}
	}
}

// BenchmarkCode2Id-8   	 2015258	       608.7 ns/op
func BenchmarkCode2Id(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := int64(rand.Uint32())
		code, err := Code6(id)

		if err != nil {
			b.Fatalf("id:%d code:%s", id, code)
		}

		decodeId, _ := Code2Id(code)
		if id != decodeId {
			b.Fatalf("id:%d decodeId:%d code:%s", id, decodeId, code)
		}
	}
}

// BenchmarkIdCode_Code2Id-8   	 2647532	       468.5 ns/op
func BenchmarkIdCode_Code2Id(b *testing.B) {
	alphanumeric := []byte("D6uLv9xy7bB5AfCHNR8VK4aFzPcZeUghIM0Q1EXi3SjTkYmnGpqrJst")
	ic, _ := NewIdCode(alphanumeric, defaultSalt6-1000, defaultSalt8-20000)

	b.Logf("max6:%d max8:%d", ic.Max6(), ic.Max8())

	for i := 0; i < b.N; i++ {
		id := int64(rand.Uint32())
		code, err := ic.Code6(id)

		if err != nil {
			b.Fatalf("id:%d code:%s", id, code)
		}

		decodeId, _ := ic.Code2Id(code)
		if id != decodeId {
			b.Fatalf("id:%d decodeId:%d code:%s", id, decodeId, code)
		}
	}
}
