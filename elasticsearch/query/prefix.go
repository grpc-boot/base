package query

type Prefix struct {
	field           string
	Value           string `json:"value"`
	Rewrite         string `json:"rewrite,omitempty"`
	CaseInsensitive bool   `json:"case_insensitive,omitempty"`
}
