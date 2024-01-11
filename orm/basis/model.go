package basis

type Model interface {
	TableName() string
	PrimaryField() string
	PrimaryValue() any
	GetDefault(name string) string
	GetLabel(name string) string
	GetEnums(name string) []string
	GetSize(name string) int
	Query() Query
	Deleter() Delete
	Inserter() Insert
	Updater() Update
}
