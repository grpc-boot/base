package query

type Terms struct {
	field  string
	values []any
	Boost  float64
}
