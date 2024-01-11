package basis

import "testing"

func TestMapSetter_Build(t *testing.T) {
	ss := StringSetter(`id=234, name='sfasfd', blob='8Y'`)
	str, err := ss.Build()
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	t.Logf("ss: %s", str)

	ms := MapSetter{
		"id":   234,
		"name": "sfasfd",
		"blob": []byte{24, 56, 89},
	}

	str, err = ms.Build()
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	t.Logf("ms: %s", str)
}

func BenchmarkParseMapping(b *testing.B) {

}
