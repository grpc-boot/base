package http_client

import (
	"testing"
	"time"
)

var (
	c = NewClient(DefaultOptions())
)

func TestClient_GetTimeout(t *testing.T) {
	rp, err := c.GetTimeout(time.Second, "https://www.baidu.com", nil)

	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	t.Logf("%d: %s", rp.GetStatus(), rp.GetBody())
}
