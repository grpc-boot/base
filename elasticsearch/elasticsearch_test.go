package elasticsearch

import (
	"context"
	"testing"
	"time"

	"github.com/grpc-boot/base/v2/elasticsearch/result"
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

	t.Logf("res: %+v", res)

	mapping := result.Mapping{
		Properties: result.MappingProperties{
			"id": result.MappingProperty{
				Type: "long",
			},
			"name": result.MappingProperty{
				Type: "keyword",
			},
			"remark": result.MappingProperty{
				Type: "text",
			},
			"status": result.MappingProperty{
				Type: "short",
			},
			"created_at": result.MappingProperty{
				Type:   "date",
				Format: "yyyy-MM-dd HH:mm:ss||epoch_second",
			},
		},
	}

	res, err = p.IndexMapping(ctx, index, mapping.Properties)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("res: %+v", res)

	res, err = p.IndexSetting(ctx, index, WithArg("index", Property{
		"max_result_window": 100,
	}))
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	t.Logf("res: %+v", res)

	mapping.Properties["updated_at"] = result.MappingProperty{
		Type:   "date",
		Format: "yyyy-MM-dd HH:mm:ss||epoch_second",
	}

	res, err = p.IndexMapping(ctx, index, mapping.Properties)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("res: %+v", res)

	gotMapping, err := p.IndexMappingGet(ctx, index)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	t.Logf("error: %v res: %+v", err, gotMapping)

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

	res, err := p.DocDel(ctx, index, "empty")
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("res: %+v", res)

	res, err = p.DocIndex(ctx, index, Body{
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
	t.Logf("res: %+v", res)

	res, err = p.DocIndexWithId(ctx, index, res.Id, Body{
		"name":       time.Now().Format(time.Layout),
		"remark":     "不好的",
		"status":     1,
		"created_at": time.Now().Unix(),
		"updated_at": time.Now().Format(time.DateTime),
	})

	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("res: %+v", res)

	docRes, err := p.DocMGet(ctx, index, res.Id)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("res: %+v", docRes)

	time.Sleep(time.Second)

	res, err = p.DocUpdateWithOptimistic(ctx, index, res.Id, Setter{
		"updated_at": time.Now().Unix(),
		"id":         100,
	}, res.SeqNo, res.PrimaryTerm)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	t.Logf("res: %+v", res)

	time.Sleep(time.Second)
	res, err = p.DocUpdate(ctx, index, res.Id, Setter{
		"created_at": time.Now().Unix(),
		"name":       "'s'adf我也是",
	})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	t.Logf("res: %+v", res)

	res, err = p.DocFieldIncr(ctx, index, res.Id, "id", 10)
	t.Logf("res: %+v", res)
}
