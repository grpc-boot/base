package elasticsearch

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var (
	p *Pool
)

func init() {
	opt := DefaultOption()
	opt.BaseUrl = "http://127.0.0.1:9200"
	p = NewPool(opt)
}

func TestPool_Index(t *testing.T) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		index       = "access_log_2024"
	)

	defer cancel()

	res, err := p.IndexDel(ctx, index)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("data: %+v", res)

	res, err = p.Index(ctx, index)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	t.Logf("hasErr: %v res: %+v", res.HasError(), res)

	properties := Properties{
		"id": Property{
			"type": "long",
		},
		"name": Property{
			"type": "keyword",
		},
		"remark": Property{
			"type": "text",
		},
		"status": Property{
			"type": "short",
		},
		"created_at": Property{
			"type":   "date",
			"format": "yyyy-MM-dd HH:mm:ss||epoch_second",
		},
	}

	res, err = p.IndexMapping(ctx, index, properties)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("hasErr: %v res: %+v", res.HasError(), res)

	res, err = p.IndexSetting(ctx, index, WithArg("index", Property{
		"max_result_window": 100,
	}))
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	t.Logf("hasErr: %v res: %+v", res.HasError(), res)

	properties["updated_at"] = Property{
		"type":   "date",
		"format": "yyyy-MM-dd HH:mm:ss||epoch_second",
	}

	res, err = p.IndexMapping(ctx, index, properties)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("hasErr: %v res: %+v", res.HasError(), res)

	mapping, err := p.IndexMappingGet(ctx, index)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	t.Logf("error: %v res: %+v", err, mapping)

	resp, err := p.IndexSettingGet(ctx, index)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("status: %d body: %s", resp.GetStatus(), resp.GetBody())
}

func TestPool_DocIndex(t *testing.T) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		index       = "access_log_2024"
	)

	defer cancel()

	res, err := p.DocIndex(ctx, index, Body{
		"id":         1,
		"name":       time.Now().Format(time.Layout),
		"remark":     "cool perfect 好的",
		"status":     1,
		"created_at": time.Now().Unix(),
		"updated_at": time.Now().Format(time.DateTime),
	})

	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("error: %+v res: %+v", err, res)

	docRes, err := p.DocMGet(ctx, index, res.Id)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("error: %+v res: %+v", err, docRes)

	resp, err := p.DocUpdate(ctx, index, res.Id, Body{
		"script": fmt.Sprintf("ctx._source.id = %d", time.Now().Unix()),
	})

	t.Logf("status: %d body: %s", resp.GetStatus(), resp.GetBody())

	resp, err = p.DocUpdateWithOptimistic(ctx, index, res.Id, Body{
		"doc": WithArg("id", time.Now().Unix()),
	}, res.SeqNo, res.PrimaryTerm)

	t.Logf("status: %d body: %s", resp.GetStatus(), resp.GetBody())

	res, err = p.DocDel(ctx, index, "")
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("error: %+v res: %+v", err, docRes)
}
