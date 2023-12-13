package monitor

import (
	"github.com/grpc-boot/base/v2/components"

	"google.golang.org/grpc/codes"
)

// ItpCode Increment of time period code，某个时间段状态码增量
type ItpCode struct {
	name string
	data map[codes.Code]*Itp
}

func NewItpCode(name string, codeList []codes.Code) *ItpCode {
	qg := &ItpCode{
		name: name,
		data: make(map[codes.Code]*Itp, len(codeList)),
	}

	for _, code := range codeList {
		qg.data[code] = NewItp(code.String(), 0)
	}

	return qg
}

func (ic *ItpCode) Add(code codes.Code, delta uint64) (newValue uint64, exists bool) {
	_, exists = ic.data[code]
	if exists {
		newValue = ic.data[code].Add(delta)
	}

	return
}

func (ic *ItpCode) AddWithStatus(sts *components.Status, delta uint64) (newValue uint64, exists bool) {
	return ic.Add(sts.Code, delta)
}

func (ic *ItpCode) Reset(times int) {
	for _, p := range ic.data {
		p.Reset(times)
	}
}

func (ic *ItpCode) Info() (info Info) {
	var (
		total uint64
		value uint64
	)

	info.Sub = make([]Item, 0, len(ic.data))

	for code, p := range ic.data {
		item := Item{
			Name:  p.name,
			Code:  code,
			Total: p.Total(),
			Value: p.Value(),
		}

		total += item.Total
		value += item.Value

		info.Sub = append(info.Sub, item)
	}

	info.Name = ic.name
	info.Total = total
	info.Value = value

	return
}
