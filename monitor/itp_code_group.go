package monitor

import (
	"sort"

	"google.golang.org/grpc/codes"
)

// ItpCodeGroup Increment of time period group，某个时间段状态码增量组，按照路径划分
type ItpCodeGroup struct {
	name     string
	codeList []codes.Code
	groups   map[string]*ItpCode
}

func NewItpCodeGroup(name string, codeList []codes.Code) *ItpCodeGroup {
	ig := &ItpCodeGroup{
		name:     name,
		codeList: codeList,
		groups:   map[string]*ItpCode{},
	}

	return ig
}

func (icg *ItpCodeGroup) Path(path, name string) *ItpCodeGroup {
	icg.groups[path] = NewItpCode(name, icg.codeList)
	return icg
}

func (icg *ItpCodeGroup) Add(path string, code codes.Code, delta uint64) (newValue uint64, exists bool) {
	if val, ok := icg.groups[path]; ok {
		return val.Add(code, delta)
	}
	return
}

func (icg *ItpCodeGroup) Reset(times int) {
	if len(icg.groups) < 1 {
		return
	}

	for _, path := range icg.groups {
		path.Reset(times)
	}
}

func (icg *ItpCodeGroup) Info() (info []CodeInfo) {
	var (
		index int
	)

	info = make([]CodeInfo, len(icg.groups))

	for path, q := range icg.groups {
		qInfo := q.Info()
		qInfo.Path = path
		info[index] = CodeInfo{
			GaugeName: icg.name,
			Info:      qInfo,
		}
		index++
	}

	sort.SliceStable(info, func(i, j int) bool {
		if info[i].Value == info[j].Value {
			return info[i].Total > info[j].Total
		}

		return info[i].Value > info[j].Value
	})

	return
}
