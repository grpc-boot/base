package elasticsearch

import (
	"context"
	"golang.org/x/exp/rand"
	"testing"
	"time"

	"github.com/grpc-boot/base/v2/elasticsearch/query"
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

func TestNewQueryString(t *testing.T) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		index       = "access_log_2024"
	)

	defer cancel()

	qs := NewQueryString("created_at:>0")
	//qs.Source = false
	qs.Sort = query.Sort{
		query.SortItem{
			"id": query.OrderItem{
				Order: query.SortDesc,
			},
		},
	}
	qs.SearchAfter = []any{1705754063}
	res, err := p.QueryWithString(ctx, index, qs)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	t.Logf("res: %+v", res)
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

func TestPool_Query(t *testing.T) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		index       = "access_log_2024"
	)
	defer cancel()

	query := AcquireQuery().
		From(index).
		Limit(5)
	defer query.Close()

	res, err := p.Query(ctx, query, "json")
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("res: %+v", res)
	records, err := res.ToRecord()
	t.Logf("records: %+v error: %v", records, err)
}

func TestPool_Bulk(t *testing.T) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		index       = "access_log_2024"
	)

	defer cancel()

	bi := NewBulkItem()
	bi.WithCreate(index, "2", Property{
		"name":       time.Now().Format(time.DateTime),
		"status":     rand.Intn(10),
		"created_at": time.Now().Unix() + 10,
		"updated_at": time.Now().Unix() + 10,
	}).
		WithUpdate(index, "2", Updater{
			Doc: Property{
				"updated_at": time.Now().Unix() + 15,
			},
		}).
		WithDelete(index, "1").
		WithDelete(index, "2")

	res, err := p.Bulk(ctx, bi)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	t.Logf("res: %+v", res)
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
