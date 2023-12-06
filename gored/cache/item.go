package cache

//go:generate msgp

// Item 缓存Item
type Item struct {
	CreatedAt int64       `msg:"c"`
	UpdatedAt int64       `msg:"u"`
	Value     interface{} `msg:"v"`
}
