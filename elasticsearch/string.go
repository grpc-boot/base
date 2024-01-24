package elasticsearch

type StringQuery struct {
	QueryString String `json:"query_string"`
}

type String struct {
	Query string `json:"query"`
}
