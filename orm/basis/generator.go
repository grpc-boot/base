package basis

type Generator interface {
	ShowTables(pattern string) (tables []string, err error)
	LoadTableSchema(table string) (t *Table, err error)
}
