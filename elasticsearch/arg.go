package elasticsearch

type Arg struct {
	key   string
	value any
}

func WithArg(key string, value any) Arg {
	return Arg{
		key:   key,
		value: value,
	}
}
