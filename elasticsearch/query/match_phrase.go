package query

type MatchPhrase struct {
	MatchPhrase map[string]MatchPhraseItem `json:"match_phrase"`
}

type MatchPhraseItem struct {
	Query    string `json:"query"`
	Operator string `json:"operator,omitempty"`
	Analyzer string `json:"analyzer,omitempty"`
}

func NewMatchPhrase(field string, item MatchPhraseItem) MatchPhrase {
	return MatchPhrase{
		MatchPhrase: map[string]MatchPhraseItem{
			field: item,
		},
	}
}
