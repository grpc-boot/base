package components

import (
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

func TestCode12(t *testing.T) {
	var (
		id  uint64
		num = 150
	)

	rand.Seed(uint64(time.Now().Unix()))

	for j := uint64(1); j < 123; j++ {
		code6, _ := DefaultIdCode.Id2Code32(uint32(j))
		decode6Id, _ := DefaultIdCode.Code2Id(code6)
		if j != decode6Id {
			t.Fatalf("id:%d code:%s decodeId:%d", j, code6, decode6Id)
		}

		code12, _ := DefaultIdCode.Id2Code64(j)
		decode12Id, _ := DefaultIdCode.Code2Id(code12)
		if j != decode12Id {
			t.Fatalf("id:%d code:%s decodeId:%d", j, code12, decode12Id)
		}

		t.Logf("id: %d 6:%s 12:%s", j, code6, code12)
	}

	for i := 0; i < num; i++ {
		id = rand.Uint64()
		code12, _ := DefaultIdCode.Id2Code64(id)
		decode12Id, _ := DefaultIdCode.Code2Id(code12)
		if id != decode12Id {
			t.Fatalf("id:%d code:%s decodeId:%d", id, code12, decode12Id)
		}

		t.Logf("id: %d 12:%s", id, code12)
	}
}

// BenchmarkId2Code32-11    	21275890	        56.12 ns/op
func BenchmarkId2Code32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := rand.Uint32()
		code, err := DefaultIdCode.Id2Code32(id)

		if err != nil {
			b.Fatalf("id:%d code:%s", id, code)
		}

		decodeId, _ := DefaultIdCode.Code2Id(code)
		if id != uint32(decodeId) {
			b.Fatalf("id:%d decodeId:%d code:%s", id, decodeId, code)
		}
	}
}

// BenchmarkId2Code64-11    	10593519	        95.92 ns/op
func BenchmarkId2Code64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := rand.Uint64()
		code, err := DefaultIdCode.Id2Code64(id)

		if err != nil {
			b.Fatalf("id:%d code:%s", id, code)
		}

		decodeId, _ := DefaultIdCode.Code2Id(code)
		if id != decodeId {
			b.Fatalf("id:%d decodeId:%d code:%s", id, decodeId, code)
		}
	}
}
