package utils

import "testing"

func TestMinuteAxis(t *testing.T) {
	axis := MinuteAxis()
	t.Logf("%v", axis)
}

func TestFiveMinuteAxis(t *testing.T) {
	axis := FiveMinuteAxis()
	t.Logf("%v", axis)
}
