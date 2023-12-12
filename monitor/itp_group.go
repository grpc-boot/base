package monitor

import (
	"sort"
)

type ItpGroup struct {
	name   string
	groups map[string]*Itp
}

func NewItpGroup(name string) *ItpGroup {
	ig := &ItpGroup{
		name:   name,
		groups: map[string]*Itp{},
	}

	return ig
}

func (ig *ItpGroup) Path(path, name string) {
	ig.groups[path] = NewItp(name, 0)
}

func (ig *ItpGroup) Add(path string, delta uint64) (newValue uint64, exists bool) {
	if val, ok := ig.groups[path]; ok {
		exists = ok
		newValue = val.Add(delta)
	}
	return
}

func (ig *ItpGroup) Reset(times int) {
	if len(ig.groups) < 1 {
		return
	}

	for _, itp := range ig.groups {
		itp.Reset(times)
	}
}

func (ig *ItpGroup) Info() (info []Info) {
	var (
		index int
	)

	info = make([]Info, len(ig.groups))

	for path, g := range ig.groups {
		info[index] = Info{
			Path:  path,
			Name:  ig.name,
			Total: g.Total(),
			Value: g.Value(),
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
