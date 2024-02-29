package redistack

import (
	"testing"
	"time"

	"github.com/redis/rueidis"
)

var (
	client    rueidis.Client
	streamKey = `key:stream`
	group     = `group2`
	consumer  = `consumer2`
)

func init() {
	var err error
	client, err = rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"10.16.49.131:6379"},
	})

	if err != nil {
		panic(err)
	}
}

func TestStreamAdd(t *testing.T) {
	id, err := StreamAdd(time.Second, client, streamKey, "10000", FieldValue{
		"k1":   "value1",
		"key2": "34",
	})
	if err != nil {
		t.Fatalf("produce msg failed with error: %+v", err)
	}
	t.Logf("produce msg success and got id: %s", id)
}

func TestStreamGroupAdd(t *testing.T) {
	isNew, err := StreamGroupCreate(time.Second, client, streamKey, group, `0-0`)
	if err != nil {
		t.Fatalf("create group failed with error: %+v", err)
	}

	if isNew {
		t.Logf("create new group")
	} else {
		t.Logf("group is exists")
	}
}

func TestStreamGroupCreateConsumer(t *testing.T) {
	isNew, err := StreamGroupCreateConsumer(time.Second, client, streamKey, group, consumer)
	if err != nil {
		t.Fatalf("create group consumer failed with error: %+v", err)
	}

	if isNew {
		t.Logf("create new group consumer")
	} else {
		t.Logf("group consumer is exists")
	}
}
