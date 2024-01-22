package query

type Fuzzy struct {
	field          string
	Value          string `json:"value"`
	Fuzziness      string `json:"fuzziness,omitempty"`
	MaxExpansions  int    `json:"max_expansions,omitempty"`
	PrefixLength   int    `json:"prefix_length,omitempty"`
	Rewrite        string `json:"rewrite,omitempty"`
	Transpositions bool   `json:"transpositions"`
}
