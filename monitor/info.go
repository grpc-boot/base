package monitor

import "google.golang.org/grpc/codes"

type MonitorInfo struct {
	ResetCount uint64            `json:"resetCount"`
	ResetAt    string            `json:"resetAt"`
	CodesInfo  map[string][]Info `json:"codesInfo"`
	GaugesInfo []Info            `json:"gaugesInfo"`
}

type Info struct {
	Path  string `json:"path"`
	Name  string `json:"name"`
	Total uint64 `json:"total"`
	Value uint64 `json:"value"`
	Sub   []Item `json:"sub,omitempty"`
}

type Item struct {
	Name  string     `json:"name"`
	Code  codes.Code `json:"code"`
	Total uint64     `json:"total"`
	Value uint64     `json:"value"`
}
