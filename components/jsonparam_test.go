package components

import "testing"

func jsonParamProvider() [][]byte {
	return [][]byte{
		[]byte(`{}`),
		[]byte(`{"id": 1, "name": "masco"}`),
		[]byte(`{"id": 2, "name": "masco", "score": 55.6}`),
		[]byte(`{"id": 2, "name": "masco", "score": 55.6, "sex":null}`),
		[]byte(`{"id": 2, "name": "masco", "score": 55.6, "sex":null, "tags":[1, 2, 3]}`),
		[]byte(`{"id": 2, "name": "masco", "score": 55.6, "sex":null, "tags":[1, 2, 3], "hobby":["basketball","football"]}`),
	}
}

func TestUnmarshalJsonParam(t *testing.T) {
	cL := jsonParamProvider()

	for _, data := range cL {
		p, err := UnmarshalJsonParam(data)
		if err != nil {
			t.Fatalf("want nil, got %+v", err)
		}
		t.Logf("%+v", p)
	}
}

func TestJsonParam_Get(t *testing.T) {
	cL := jsonParamProvider()

	for _, data := range cL {
		p, _ := UnmarshalJsonParam(data)
		t.Logf("id:%d name:%s score:%10.2f sex:%v tags:%v hobby:%v", p.Int64("id"), p.String("name"), p.Float64("score"), p.Int("sex"), p.Uint32Slice("tags"), p.StringSlice("hobby"))
	}
}
