package query

type Wildcard struct {
	field           string
	Value           string  `json:"value"`
	Boost           float64 `json:"boost,omitempty"`
	Rewrite         string  `json:"rewrite,omitempty"`
	CaseInsensitive bool    `json:"case_insensitive,omitempty"`
}
