package result

type Result struct {
	Error  *Error `json:"error"`
	Status int    `json:"status"`
}

func (r *Result) HasError() bool {
	return r.Error != nil
}

type Error struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
	Index  string `json:"index"`
}
