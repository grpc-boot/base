package convert

// go test -bench='_Convert$' -benchmem -run=none

import "testing"

type TestRow struct {
	Id        uint32
	Name      string
	Sex       bool
	IsOn      bool
	Balance   float64
	Age       uint8
	CreatedAt int64
}

func TestRow_Convert(t *testing.T) {
	caseList := []Row{
		{
			"id":        "1",
			"name":      "Macos",
			"sex":       "false",
			"isOn":      "0",
			"balance":   "34.674",
			"age":       "15",
			"createdAt": "1534234324",
		},
		{
			"id":        "2034",
			"name":      "Macoasdfasfs",
			"sex":       "1",
			"isOn":      "true",
			"balance":   "34.674",
			"age":       "15",
			"createdAt": "1534234324",
		},
	}

	for _, row := range caseList {
		data1 := TestRow{}

		err := row.Convert(&data1)
		if err != nil {
			t.Fatalf("want nil, got %s", err.Error())
		}

		t.Logf("data1:%+v", data1)
	}
}

func BenchmarkRow_Convert(b *testing.B) {
	row := Row{
		"id":        "1",
		"name":      "Macos",
		"sex":       "false",
		"isOn":      "0",
		"balance":   "34.674",
		"age":       "15",
		"createdAt": "1534234324",
	}

	data1 := TestRow{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		row.Convert(&data1)
	}
}
