package monitor

import (
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
)

var (
	InfoPrefix byte = 'i'
	ItemPrefix byte = 'c'
	Seperator  byte = '-'
)

type MonitorInfo struct {
	ResetCount uint64            `json:"resetCount"`
	ResetAt    string            `json:"resetAt"`
	CodesInfo  map[string][]Info `json:"codesInfo"`
	GaugesInfo []Info            `json:"gaugesInfo"`
}

func (mi *MonitorInfo) Keys(prefix string) (gaugeKeys []string, codeKeys []string) {
	gaugeKeys = make([]string, len(mi.GaugesInfo))
	codeKeys = make([]string, 0, 256)
	index := 0

	for _, gauge := range mi.GaugesInfo {
		gaugeKeys[index] = gauge.Key(prefix, "")
		index++
	}

	for gaugeName, groups := range mi.CodesInfo {
		if len(groups) < 1 {
			continue
		}

		for _, gauge := range groups {
			codeKeys = append(codeKeys, gauge.Key(prefix, gaugeName))
		}
	}

	return
}

type Info struct {
	Path  string `json:"path"`
	Name  string `json:"name"`
	Total uint64 `json:"total"`
	Value uint64 `json:"value"`
	Sub   []Item `json:"sub,omitempty"`
}

func (i *Info) Key(prefix, gaugeName string) string {
	var buffer strings.Builder
	buffer.Grow(len(prefix) + len(gaugeName) + len(i.Path) + 2)
	buffer.WriteString(prefix)
	buffer.WriteByte(Seperator)
	buffer.WriteString(gaugeName)
	buffer.WriteByte(Seperator)
	buffer.WriteString(i.Path)
	return buffer.String()
}

func (i *Info) Field(prefix string) string {
	var buffer strings.Builder
	buffer.Grow(len(prefix) + len(i.Path) + 2)
	buffer.WriteByte(InfoPrefix)
	buffer.WriteString(prefix)
	buffer.WriteByte(Seperator)
	buffer.WriteString(i.Path)
	return buffer.String()
}

type Item struct {
	Name  string     `json:"name"`
	Code  codes.Code `json:"code"`
	Total uint64     `json:"total"`
	Value uint64     `json:"value"`
}

func (it *Item) Field(prefix string) string {
	var (
		buffer strings.Builder
		path   = strconv.FormatUint(uint64(it.Code), 10)
	)

	buffer.Grow(len(prefix) + len(path) + 2)
	buffer.WriteByte(ItemPrefix)
	buffer.WriteString(prefix)
	buffer.WriteByte(Seperator)
	buffer.WriteString(path)

	return buffer.String()
}
