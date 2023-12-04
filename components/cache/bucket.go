package cache

import (
	"sync"

	"go.uber.org/atomic"
)

//go:generate msgp

type Bucket struct {
	mutex      sync.RWMutex
	fMutex     sync.Mutex
	hasChanged atomic.Bool
	Data       map[string]*Entry `msg:"data"`
}

type Entry struct {
	_lock     int64
	updateCnt uint64
	hitCnt    atomic.Uint64
	missCnt   atomic.Uint64
	UpdatedAt int64       `msg:"u"`
	CreatedAt int64       `msg:"c"`
	Value     interface{} `msg:"v"`
}
