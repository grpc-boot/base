package query

type Range struct {
	Range map[string]RangeItem `json:"range"`
}

type RangeItem struct {
	Gt       any     `json:"gt,omitempty"`
	Gte      any     `json:"gte,omitempty"`
	Lt       any     `json:"lt,omitempty"`
	Lte      any     `json:"lte,omitempty"`
	Boost    float64 `json:"boost,omitempty"`
	Format   string  `json:"format,omitempty"`
	Relation string  `json:"relation,omitempty"`
	TimeZone string  `json:"time_zone,omitempty"`
}

func NewRange(field string, item RangeItem) Range {
	return Range{
		Range: map[string]RangeItem{
			field: item,
		},
	}
}
