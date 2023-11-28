package components

import (
	"os"
	"testing"
	"time"
)

func TestSnowFlake_Id(t *testing.T) {
	begin, _ := time.ParseInLocation("2006-01-02", `2023-01-01`, time.Local)
	sf, err := NewSFByIp(ModeWait, begin.Unix())
	if err != nil {
		t.Fatal(err)
	}

	id, _ := sf.Id()
	t.Log(id)
	ts, machineId, index := sf.Info(id)

	if index != 1 {
		t.Fatalf("want 1, got %d", index)
	}

	t.Logf("milliTimestamp:%d machine:%d, index:%d", ts, machineId, index)

	id, _ = sf.Id()
	t.Log(id)
	ts, machineId, index = sf.Info(id)

	if index != 2 {
		t.Fatalf("want 1, got %d", index)
	}

	t.Logf("milliTimestamp:%d machine:%d, index:%d", ts, machineId, index)
}

func TestNewSFByMachineFunc(t *testing.T) {
	os.Setenv("MNUM", "128")

	begin, _ := time.ParseInLocation("2006-01-02", `2021-01-01`, time.Local)
	sf, err := NewSFByMachineFunc(ModeWait, GetMachineIdByEnv("MNUM"), begin.Unix())
	if err != nil {
		t.Fatal(err)
	}

	id, _ := sf.Id()
	t.Log(id)
	ts, machineId, index := sf.Info(id)

	if index != 1 {
		t.Fatalf("want 1, got %d", index)
	}

	if machineId != 128 {
		t.Fatalf("want 128, got %d", machineId)
	}

	t.Logf("milliTimestamp:%d machine:%d, index:%d", ts, machineId, index)
}
