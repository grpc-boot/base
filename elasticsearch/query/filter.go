package query

type Filter struct {
	list []any
}

func (f *Filter) WithIds(item IdsItem) *Filter {
	f.list = append(f.list, NewIds(item))

	return f
}

func (f *Filter) WithTerm(field string, item TermItem) *Filter {
	f.list = append(f.list, NewTerm(field, item))
	return f
}

func (f *Filter) WithTerms(field string, item TermsItem) *Filter {
	f.list = append(f.list, NewTerms(field, item))
	return f
}

func (f *Filter) WithRange(field string, item RangeItem) *Filter {
	f.list = append(f.list, NewRange(field, item))
	return f
}

func (f *Filter) WithExists(item ExistsItem) *Filter {
	f.list = append(f.list, NewExists(item))
	return f
}

func (f *Filter) WithPrefix(field string, item PrefixItem) *Filter {
	f.list = append(f.list, NewPrefix(field, item))
	return f
}

func (f *Filter) WithWildcard(field string, item WildcardItem) *Filter {
	f.list = append(f.list, NewWildcard(field, item))
	return f
}

func (f *Filter) WithRegExp(field string, item RegExpItem) *Filter {
	f.list = append(f.list, NewRegExp(field, item))
	return f
}

func (f *Filter) WithFuzzy(field string, item FuzzyItem) *Filter {
	f.list = append(f.list, NewFuzzy(field, item))
	return f
}
