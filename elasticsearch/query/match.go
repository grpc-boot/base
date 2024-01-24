package query

import "github.com/grpc-boot/base/v2/kind"

type Match struct {
	field    string
	Query    string `json:"query"`
	Operator string `json:"operator,omitempty"`
	Analyzer string `json:"analyzer,omitempty"`
}

func (m Match) Param() kind.JsonParam {
	return kind.JsonParam{
		"match": kind.JsonParam{
			m.field: m,
		},
	}
}

func NewMatch(field, query string) Match {
	return Match{
		field: field,
		Query: query,
	}
}
