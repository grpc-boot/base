package components

import (
	"strconv"
	"testing"

	"github.com/grpc-boot/base/v2/internal"
	"github.com/grpc-boot/base/v2/kind"

	"go.uber.org/atomic"
)

var (
	hashGroup = &Group{}
	hostList  = []string{
		"192.168.1.135:3551:v0",
		"192.168.1.135:3551:v1",
		"192.168.1.135:3551:v2",
		"192.168.1.135:3551:v3",
		"192.168.1.135:3551:v4",
		"192.168.1.135:3552:v0",
		"192.168.1.135:3552:v1",
		"192.168.1.135:3552:v2",
		"192.168.1.135:3552:v3",
		"192.168.1.135:3552:v4",
		"192.168.1.135:3553:v0",
		"192.168.1.135:3553:v1",
		"192.168.1.135:3553:v2",
		"192.168.1.135:3553:v3",
		"192.168.1.135:3553:v4",
		"192.168.1.135:3554:v0",
		"192.168.1.135:3554:v1",
		"192.168.1.135:3554:v2",
		"192.168.1.135:3554:v3",
		"192.168.1.135:3554:v4",
		"192.168.1.135:3555:v0",
		"192.168.1.135:3555:v1",
		"192.168.1.135:3555:v2",
		"192.168.1.135:3555:v3",
		"192.168.1.135:3555:v4",
		"192.168.1.135:3556:v0",
		"192.168.1.135:3556:v1",
		"192.168.1.135:3556:v2",
		"192.168.1.135:3556:v3",
		"192.168.1.135:3556:v4",
		"192.168.1.135:3557:v0",
		"192.168.1.135:3557:v1",
		"192.168.1.135:3557:v2",
		"192.168.1.135:3557:v3",
		"192.168.1.135:3557:v4",
		"192.168.1.135:3558:v0",
		"192.168.1.135:3558:v1",
		"192.168.1.135:3558:v2",
		"192.168.1.135:3558:v3",
		"192.168.1.135:3558:v4",
		"192.168.1.135:3559:v0",
		"192.168.1.135:3559:v1",
		"192.168.1.135:3559:v2",
		"192.168.1.135:3559:v3",
		"192.168.1.135:3559:v4",
	}
)

type Data struct {
	kind.CanHash

	id string
}

func (d *Data) HashCode() (hashValue uint32) {
	return kind.Uint32Hash(internal.String2Bytes(d.id))
}

type Group struct {
	ring HashRing
}

func TestHashRing_Get(t *testing.T) {
	serverList := make([]kind.CanHash, 0, len(hostList))

	for _, server := range hostList {
		serverList = append(serverList, &Data{
			id: server,
		})
	}

	hashGroup.ring = NewHashRing(serverList...)

	for end := 1 << 25; end > 0; end-- {
		_, err := hashGroup.ring.Get(strconv.Itoa(end))
		if err != nil {
			t.Fatal(err)
		}
	}

	hashGroup.ring.Range(func(index int, server kind.CanHash, hitCount uint64) (handled bool) {
		t.Logf("index:%d, server.id:%s, hitCount:%d", index, server.(*Data).id, hitCount)
		return
	})
}

// go test -bench=. -benchmem -v
// BenchmarkHashRing_Get-8   	 7349292	       155.2 ns/op
func BenchmarkHashRing_Get(b *testing.B) {
	serverList := make([]kind.CanHash, 0, len(hostList))

	for _, server := range hostList {
		serverList = append(serverList, &Data{
			id: server,
		})
	}

	hashGroup.ring = NewHashRing(serverList...)
	var val atomic.Uint64

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := hashGroup.ring.Get(internal.String2Bytes(strconv.FormatUint(val.Add(1), 10)))
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}
