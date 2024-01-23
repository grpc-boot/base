package query

type Exists struct {
	Exists ExistsItem `json:"exists"`
}

type ExistsItem struct {
	Field string `json:"field"`
}

func NewExists(item ExistsItem) Exists {
	return Exists{
		Exists: item,
	}
}
