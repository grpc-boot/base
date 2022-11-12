//go:generate msgp
package cache

type Bucket struct {
	Items map[string]*Item `json:"items" msg:"items"`
}

type Item struct {
	Value    []byte `json:"value" msg:"value"`
	ExpireAt int64  `json:"expireAt" msg:"expireAt"`
}
