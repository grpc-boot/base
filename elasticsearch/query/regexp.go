package query

type RegExp struct {
	field                 string
	Value                 string `json:"value"`
	Flags                 string `json:"flags"`
	MaxDeterminizedStates int    `json:"max_determinized_states,omitempty"`
	Rewrite               string `json:"rewrite,omitempty"`
	CaseInsensitive       bool   `json:"case_insensitive,omitempty"`
}
