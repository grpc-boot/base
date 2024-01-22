package query

type MultiMatch struct {
	Query    string   `json:"query"`
	Fields   []string `json:"fields"`
	Type     string   `json:"type,omitempty"`
	Operator string   `json:"operator,omitempty"`
	Analyzer string   `json:"analyzer,omitempty"`
}
