package monitor

import (
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
)

var (
	Seperator byte = '-'
)

type MonitorInfo struct {
	Name       string                `json:"name"`
	ResetCount uint64                `json:"resetCount"`
	ResetAt    string                `json:"resetAt"`
	CodesInfo  map[string][]CodeInfo `json:"codesInfo"`
	GaugesInfo []Info                `json:"gaugesInfo"`
}

func (mi *MonitorInfo) Keys(prefix string) (gaugeKeys []string, codeKeys []string) {
	gaugeKeys = make([]string, len(mi.GaugesInfo))
	codeKeys = make([]string, 0, 256)
	index := 0

	for _, gauge := range mi.GaugesInfo {
		gaugeKeys[index] = gauge.Key(mi.Name, prefix)
		index++
	}

	for _, groups := range mi.CodesInfo {
		if len(groups) < 1 {
			continue
		}

		for _, gauge := range groups {
			codeKeys = append(codeKeys, gauge.Key(mi.Name, prefix))
		}
	}

	return
}

type Info struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Total uint64 `json:"total"`
	Value uint64 `json:"value"`
	Sub   []Item `json:"sub,omitempty"`
}

func (i *Info) Key(appName, prefix string) string {
	var buffer strings.Builder
	buffer.Grow(len(appName) + len(prefix) + len(i.Path) + 2)
	buffer.WriteString(appName)
	buffer.WriteByte(Seperator)
	buffer.WriteString(prefix)
	buffer.WriteByte(Seperator)
	buffer.WriteString(i.Path)
	return buffer.String()
}

type CodeInfo struct {
	Info

	GaugeName string `json:"gaugeName"`
}

func (ci *CodeInfo) Key(appName, prefix string) string {
	var buffer strings.Builder
	buffer.Grow(len(appName) + len(prefix) + len(ci.GaugeName) + 3)
	buffer.WriteString(appName)
	buffer.WriteByte(Seperator)
	buffer.WriteString(prefix)
	buffer.WriteByte(Seperator)
	buffer.WriteString(ci.GaugeName)
	buffer.WriteByte(Seperator)
	buffer.WriteString(ci.Path)
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

	buffer.Grow(len(prefix) + len(path) + 1)
	buffer.WriteString(prefix)
	buffer.WriteByte(Seperator)
	buffer.WriteString(path)

	return buffer.String()
}
