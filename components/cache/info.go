package cache

type Info struct {
	LocalDir       string `json:"localDir"`
	FlushInterval  string `json:"flushInterval"`
	Length         int64  `json:"length"`
	LatestSyncTime string `json:"latestSyncTime"`
	Items          []Item `json:"items"`
}
