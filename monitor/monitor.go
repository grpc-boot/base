package monitor

import (
	"time"

	"github.com/grpc-boot/base/v2/components"
	"github.com/grpc-boot/base/v2/logger"
	"github.com/grpc-boot/base/v2/utils"

	"go.uber.org/atomic"
	"google.golang.org/grpc/codes"
)

const (
	GaugeRequestCount  = `requestCount`
	GaugeResponseCount = `responseCount`
	GaugeRequestLen    = `requestLen`
	GaugeResponseLen   = `responseLen`
	GaugePanicCount    = `panicCount`
)

const (
	defaultResetSeconds = 60
	defaultResetTimes   = 1
)

var (
	RootPath = `/`
)

type Monitor struct {
	opt        Options
	resetAt    atomic.Time
	resetCount atomic.Uint64
	ticker     *time.Ticker
	codeGauges map[string]*ItpCodeGroup
	gauges     *ItpGroup
	storage    Storage
}

func NewMonitor(opt Options) *Monitor {
	opt = formatOpt(opt)

	m := &Monitor{
		opt:        opt,
		gauges:     NewItpGroup("gauge"),
		codeGauges: make(map[string]*ItpCodeGroup, len(opt.CodeGauges)),
		ticker:     time.NewTicker(time.Duration(opt.ResetSeconds) * time.Second),
	}

	for _, gauge := range opt.CodeGauges {
		m.codeGauges[gauge] = NewItpCodeGroup(gauge, opt.CodeList)
	}

	if len(opt.Gauges) > 0 {
		for _, gauge := range opt.Gauges {
			m.gauges.Path(gauge, gauge)
		}
	}

	m.Path(RootPath, "root")

	go m.tick()

	return m
}

func (m *Monitor) tick() {
	for range m.ticker.C {
		m.gauges.Reset(m.opt.ResetTimes)

		for _, group := range m.codeGauges {
			group.Reset(m.opt.ResetTimes)
		}

		m.resetCount.Inc()
		m.resetAt.Store(time.Now())

		if m.storage != nil {
			go utils.Recover("monitor storage", func(args ...any) {
				if err := m.storage(m.Info()); err != nil {
					logger.ZapError("storage monitor info failed",
						logger.Error(err),
					)
				}
			})
		}
	}
}

func (m *Monitor) WithStorage(s Storage) *Monitor {
	m.storage = s
	return m
}

func (m *Monitor) AddGauge(gauge string, delta uint64) (newValue uint64, exists bool) {
	return m.gauges.Add(gauge, delta)
}

func (m *Monitor) Add(gauge, path string, code codes.Code, delta uint64) (newValue uint64, exists bool) {
	if group, ok := m.codeGauges[gauge]; ok {
		newValue, exists = group.Add(path, code, delta)
		group.Add(RootPath, code, delta)
	}
	return
}

func (m *Monitor) AddWithStatus(gauge, path string, sts *components.Status, delta uint64) (newValue uint64, exists bool) {
	return m.Add(gauge, path, sts.Code, delta)
}

func (m *Monitor) Path(path, name string) *Monitor {
	for _, group := range m.codeGauges {
		group.Path(path, name)
	}

	return m
}

func (m *Monitor) Info() MonitorInfo {
	info := make(map[string][]Info, len(m.codeGauges))

	for gauge, group := range m.codeGauges {
		info[gauge] = group.Info()
	}

	return MonitorInfo{
		ResetAt:    m.resetAt.Load().Format(time.DateTime),
		ResetCount: m.resetCount.Load(),
		CodesInfo:  info,
		GaugesInfo: m.gauges.Info(),
	}
}
