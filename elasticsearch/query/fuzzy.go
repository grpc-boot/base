package query

type Fuzzy struct {
	Fuzzy map[string]FuzzyItem `json:"fuzzy"`
}

type FuzzyItem struct {
	Value          string `json:"value"`
	Fuzziness      string `json:"fuzziness,omitempty"`
	MaxExpansions  int    `json:"max_expansions,omitempty"`
	PrefixLength   int    `json:"prefix_length,omitempty"`
	Rewrite        string `json:"rewrite,omitempty"`
	Transpositions bool   `json:"transpositions"`
}

func NewFuzzy(field string, item FuzzyItem) Fuzzy {
	return Fuzzy{
		Fuzzy: map[string]FuzzyItem{
			field: item,
		},
	}
}
