package query

type Ids struct {
	Ids IdsItem `json:"ids"`
}

type IdsItem struct {
	Values []string `json:"values"`
}

func NewIds(item IdsItem) Ids {
	return Ids{
		Ids: item,
	}
}
