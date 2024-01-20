package basis

type Generator interface {
	Executor
	ShowTables(pattern string) (tables []string, err error)
	LoadTableSchema(table string) (t *Table, err error)
}
