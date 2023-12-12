package monitor

import "go.uber.org/atomic"

type Gauge struct {
	name  string
	value atomic.Uint64
}

func (g *Gauge) Set(value uint64) {
	g.value.Store(value)
}

func (g *Gauge) Value() uint64 {
	return g.value.Load()
}
