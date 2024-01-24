package query

import "github.com/grpc-boot/base/v2/kind"

type Term struct {
	field           string
	Value           any     `json:"value"`
	Boost           float64 `json:"boost,omitempty"`
	CaseInsensitive bool    `json:"case_insensitive,omitempty"`
}

func (t Term) Param() kind.JsonParam {
	return kind.JsonParam{
		"term": kind.JsonParam{
			t.field: t,
		},
	}
}

func NewTerm(field string, value any) Term {
	return Term{
		field: field,
		Value: value,
	}
}
