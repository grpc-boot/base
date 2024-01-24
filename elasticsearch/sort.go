package elasticsearch

const (
	SortAsc  = `asc`
	SortDesc = `desc`
)

type Sort []SortItem

type SortItem map[string]OrderItem

type OrderItem struct {
	Order  string `json:"order"`
	Format string `json:"format,omitempty"`
	Mode   string `json:"mode,omitempty"`
}
