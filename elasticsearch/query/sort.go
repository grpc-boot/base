package query

const (
	Asc  = `asc`
	Desc = `desc`
)

type Sort struct {
	field  string
	Order  string `json:"order"`
	Format string `json:"format,omitempty"`
	Mode   string `json:"mode,omitempty"`
}
