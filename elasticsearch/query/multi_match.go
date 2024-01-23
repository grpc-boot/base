package query

type MultiMatch struct {
	MultiMatch MultiMatchItem `json:"multi_match"`
}

type MultiMatchItem struct {
	Query    string   `json:"query"`
	Fields   []string `json:"fields"`
	Type     string   `json:"type,omitempty"`
	Operator string   `json:"operator,omitempty"`
	Analyzer string   `json:"analyzer,omitempty"`
}

func NewMultiMatch(item MultiMatchItem) MultiMatch {
	return MultiMatch{
		MultiMatch: item,
	}
}
