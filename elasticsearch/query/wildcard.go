package query

type Wildcard struct {
	Wildcard map[string]WildcardItem `json:"wildcard"`
}

type WildcardItem struct {
	Value           string  `json:"value"`
	Boost           float64 `json:"boost,omitempty"`
	Rewrite         string  `json:"rewrite,omitempty"`
	CaseInsensitive bool    `json:"case_insensitive,omitempty"`
}

func NewWildcard(field string, item WildcardItem) Wildcard {
	return Wildcard{
		Wildcard: map[string]WildcardItem{
			field: item,
		},
	}
}
