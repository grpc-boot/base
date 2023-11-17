package condition

type Values []interface{}

type Condition interface {
	Build() (sql string, args []interface{})
}
