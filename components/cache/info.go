package cache

type Info struct {
	LocalDir       string           `json:"localDir"`
	FlushInterval  string           `json:"flushInterval"`
	Length         int64            `json:"length"`
	LatestSyncTime string           `json:"latestSyncTime"`
	Keys           map[int][]string `json:"keys"`
	Items          []Item           `json:"items"`
}

type Item struct {
	Key       string `json:"key"`
	UpdatedAt int64  `json:"updatedAt"`
	UpdateCnt uint64 `json:"updateCount"`
	HitCnt    uint64 `json:"hitCount"`
	MissCnt   uint64 `json:"missCount"`
	CreatedAt int64  `msg:"createdAt"`
	Value     string `msg:"value"`
}
