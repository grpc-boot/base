package query

import (
	"github.com/grpc-boot/base/v2/kind"
	"github.com/grpc-boot/base/v2/utils"
)

type Filter struct {
	list []kind.JsonParam
}

func (f *Filter) Marshal() []byte {
	data, _ := utils.JsonMarshal(f.list)
	return data
}

func (f *Filter) WithIds(ids Ids) *Filter {
	f.list = append(f.list, kind.JsonParam{
		"ids": ids,
	})

	return f
}

func (f *Filter) WithTerm(t Term) *Filter {
	f.list = append(f.list, kind.JsonParam{
		"term": kind.JsonParam{
			t.field: t,
		},
	})
	return f
}

func (f *Filter) WithTerms(ts Terms) *Filter {
	p := kind.JsonParam{
		ts.field: ts.values,
	}

	if ts.Boost != 0 {
		p["boost"] = ts.Boost
	}

	f.list = append(f.list, kind.JsonParam{
		"terms": p,
	})
	return f
}

func (f *Filter) WithRange(r Range) *Filter {
	f.list = append(f.list, kind.JsonParam{
		"range": kind.JsonParam{
			r.field: r,
		},
	})

	return f
}

func (f *Filter) WithExists(e Exists) *Filter {
	f.list = append(f.list, kind.JsonParam{
		"exists": e,
	})

	return f
}

func (f *Filter) WithPrefix(p Prefix) *Filter {
	f.list = append(f.list, kind.JsonParam{
		"prefix": kind.JsonParam{
			p.field: p,
		},
	})

	return f
}

func (f *Filter) WithWildcard(w Wildcard) *Filter {
	f.list = append(f.list, kind.JsonParam{
		"wildcard": kind.JsonParam{
			w.field: w,
		},
	})

	return f
}

func (f *Filter) WithRegExp(re RegExp) *Filter {
	f.list = append(f.list, kind.JsonParam{
		"regexp": kind.JsonParam{
			re.field: re,
		},
	})

	return f
}

func (f *Filter) WithFuzzy(fz Fuzzy) *Filter {
	f.list = append(f.list, kind.JsonParam{
		"fuzzy": kind.JsonParam{
			fz.field: fz,
		},
	})

	return f
}
