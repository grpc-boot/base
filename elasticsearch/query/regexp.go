package query

type RegExp struct {
	RegExp map[string]RegExpItem `json:"regexp"`
}

type RegExpItem struct {
	Value                 string `json:"value"`
	Flags                 string `json:"flags"`
	MaxDeterminizedStates int    `json:"max_determinized_states,omitempty"`
	Rewrite               string `json:"rewrite,omitempty"`
	CaseInsensitive       bool   `json:"case_insensitive,omitempty"`
}

func NewRegExp(field string, item RegExpItem) RegExp {
	return RegExp{
		RegExp: map[string]RegExpItem{
			field: item,
		},
	}
}
