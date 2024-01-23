package query

import "github.com/grpc-boot/base/v2/kind"

const (
	KindTerm  = 1
	KindTerms = 2
)

type Item struct {
	kind            uint8
	field           string
	boost           float64
	query           string
	value           any
	values          []any
	caseInsensitive bool
	rewrite         string
	operator        string
	analyzer        string
}

func (i *Item) ToParam() kind.JsonParam {
	switch i.kind {
	case KindTerm:
		return i.toTerm()
	case KindTerms:
		return i.toTerms()
	default:
		return nil
	}
}

func NewItem(kind uint8, field string) *Item {
	return &Item{
		kind:  kind,
		field: field,
	}
}

func (i *Item) WithQuery(query string) *Item {
	i.query = query
	return i
}

func (i *Item) WithBoost(value float64) *Item {
	i.boost = value
	return i
}

func (i *Item) WithAnalyzer(analyzer string) *Item {
	i.analyzer = analyzer
	return i
}

func (i *Item) WithValue(value any) *Item {
	i.value = value
	return i
}

func (i *Item) WithValues(values []any) *Item {
	i.values = values
	return i
}

func (i *Item) WithCaseInsensitive(ci bool) *Item {
	i.caseInsensitive = ci
	return i
}

func (i *Item) WithRewrite(rewrite string) *Item {
	i.rewrite = rewrite
	return i
}

func (i *Item) WithOperator(opt string) *Item {
	i.operator = opt
	return i
}

func (i *Item) toPrefix() kind.JsonParam {
	if i.rewrite == "" && !i.caseInsensitive {
		return nil
	}
	return nil
}

func (i *Item) toTerms() kind.JsonParam {
	if i.boost == 0 && !i.caseInsensitive {
		return kind.JsonParam{
			"terms": kind.JsonParam{
				i.field: i.values,
			},
		}
	}

	item := make(kind.JsonParam, 3)
	item[i.field] = i.value
	if i.boost > 0 {
		item["boost"] = i.boost
	}

	if i.caseInsensitive {
		item["case_insensitive"] = i.caseInsensitive
	}

	return kind.JsonParam{
		"terms": kind.JsonParam{
			i.field: item,
		},
	}
}

func (i *Item) toTerm() kind.JsonParam {
	if i.boost == 0 && !i.caseInsensitive {
		return kind.JsonParam{
			"term": kind.JsonParam{
				i.field: i.value,
			},
		}
	}

	item := make(kind.JsonParam, 3)
	item["value"] = i.value
	if i.boost > 0 {
		item["boost"] = i.boost
	}

	if i.caseInsensitive {
		item["case_insensitive"] = i.caseInsensitive
	}

	return kind.JsonParam{
		"term": kind.JsonParam{
			i.field: item,
		},
	}
}
