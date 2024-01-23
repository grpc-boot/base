package query

type Match struct {
	Match map[string]MatchItem `json:"match"`
}

type MatchItem struct {
	Query    string `json:"query"`
	Operator string `json:"operator,omitempty"`
	Analyzer string `json:"analyzer,omitempty"`
}

func NewMatch(field string, item MatchItem) Match {
	return Match{
		Match: map[string]MatchItem{
			field: item,
		},
	}
}
