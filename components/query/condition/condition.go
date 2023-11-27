package condition

type Condition interface {
	Build() (sql string, args []any)
}
