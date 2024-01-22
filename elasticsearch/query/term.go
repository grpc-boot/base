package query

type Term struct {
	field           string
	Value           any     `json:"value"`
	Boost           float64 `json:"boost,omitempty"`
	CaseInsensitive bool    `json:"case_insensitive,omitempty"`
}

func TermQuery(field string, value any) Term {
	return Term{
		field: field,
		Value: value,
	}
}
