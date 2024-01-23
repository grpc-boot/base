package query

type Prefix struct {
	Prefix map[string]PrefixItem `json:"prefix"`
}

type PrefixItem struct {
	Value           string `json:"value"`
	Rewrite         string `json:"rewrite,omitempty"`
	CaseInsensitive bool   `json:"case_insensitive,omitempty"`
}

func NewPrefix(field string, item PrefixItem) Prefix {
	return Prefix{
		Prefix: map[string]PrefixItem{
			field: item,
		},
	}
}
