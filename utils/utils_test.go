package utils

import (
	"math"
	"testing"
)

func TestJoin(t *testing.T) {
	res := Join("|", []uint64{math.MaxUint64, math.MaxInt64}...)

	t.Logf("res: %s", res)
}

func TestAbs(t *testing.T) {
	float32Val := Abs(-45.4)
	if float32Val != 45.4 {
		t.Fatalf("want 45.4, got %v", float32Val)
	}

	int64Val := Abs[int64](-45)
	if int64Val != 45 {
		t.Fatalf("want 45, got %v", int64Val)
	}

	var val int32 = -980
	int32Val := Abs(val)
	if int32Val != 980 {
		t.Fatalf("want 32, got %v", int32Val)
	}
}

func TestJsonUnmarshalFile(t *testing.T) {
	type Data struct {
		Int    int     `json:"int"`
		Float  float64 `json:"float"`
		Bool   bool    `json:"bool"`
		String string  `json:"string"`
	}

	var data Data
	err := JsonUnmarshalFile("sample.json", &data)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	t.Logf("got data: %+v", data)
}
