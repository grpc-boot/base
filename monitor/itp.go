package monitor

import (
	"math"

	"go.uber.org/atomic"
)

// Itp Increment of time period，某个时间段增量
type Itp struct {
	Gauge

	total       atomic.Uint64
	latestTotal atomic.Uint64
}

func NewItp(name string, value uint64) *Itp {
	q := &Itp{}

	q.name = name

	if value > 0 {
		q.Set(value)
	}

	return q
}

func (i *Itp) Reset(times int) (value uint64) {
	total, latestTotal := i.total.Load(), i.latestTotal.Load()
	if total >= latestTotal {
		value = total - latestTotal
	} else {
		value = math.MaxUint64 - latestTotal + total
	}

	value = uint64(math.Round(float64(value) / float64(times)))

	i.latestTotal.Store(total)
	i.Set(value)
	return
}

func (i *Itp) Add(delta uint64) (newValue uint64) {
	return i.total.Add(delta)
}

func (i *Itp) Total() uint64 {
	return i.total.Load()
}
