package query

type Match struct {
	field    string
	Query    string `json:"query"`
	Operator string `json:"operator,omitempty"`
	Analyzer string `json:"analyzer,omitempty"`
}
