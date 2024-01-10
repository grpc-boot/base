package base

type Column interface {
	Unsigned() bool
	Comment() string
	Field() string
	Name() string
	CanNull() bool
	Size() int
	Scale() int
	IsPrimaryKey() bool
	AutoIncrement() bool
	GoType() string
	Key() string
	Default() string
	Enums() []string
}
