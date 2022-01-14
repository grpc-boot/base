package base

import (
	"errors"
	"math"
	"sort"
	"sync"

	"go.uber.org/atomic"
)

var ErrNoServer = errors.New("no server")

// HashRing Hash环
type HashRing interface {
	// Store 存储servers
	Store(servers ...CanHash)
	// Get 获取server
	Get(key interface{}) (server CanHash, err error)
	// Index 根据index获取server
	Index(index int) (server CanHash, err error)
	// Add 添加server
	Add(server CanHash)
	// Remove 移除server
	Remove(server CanHash)
	// Length 获取servers长度
	Length() int
	// Range 遍历servers
	Range(handler func(index int, server CanHash, hitCount uint64) (handled bool))
}

type Node struct {
	server    CanHash
	hashValue uint32
	hitCount  atomic.Uint64
}

type NodeList []Node

func (n NodeList) Len() int {
	return len(n)
}

func (n NodeList) Less(i, j int) bool {
	return n[i].hashValue < n[j].hashValue
}

func (n NodeList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

type hashRing struct {
	HashRing

	nodes NodeList
	mutex sync.RWMutex
}

func NewHashRing(servers ...CanHash) HashRing {
	r := &hashRing{}
	r.Store(servers...)
	return r
}

func (hr *hashRing) Store(servers ...CanHash) {
	hr.mutex.Lock()
	nodes := make(NodeList, len(servers), len(servers))

	for index, _ := range servers {
		nodes[index] = Node{
			server:    servers[index],
			hashValue: servers[index].HashCode(),
		}
	}

	hr.nodes = nodes
	sort.Sort(hr.nodes)

	hr.mutex.Unlock()
}

func (hr *hashRing) Add(server CanHash) {
	hr.mutex.Lock()
	hr.nodes = append(hr.nodes, Node{
		server:    server,
		hashValue: server.HashCode(),
	})
	sort.Sort(hr.nodes)
	hr.mutex.Unlock()
}

func (hr *hashRing) Remove(server CanHash) {
	hr.mutex.Lock()

	value := server.HashCode()

	index := sort.Search(len(hr.nodes), func(i int) bool {
		return hr.nodes[i].hashValue >= value
	})

	if index < len(hr.nodes) && hr.nodes[index].hashValue == value {
		if len(hr.nodes) == 1 {
			hr.nodes = hr.nodes[:0]
		} else {
			hr.nodes = append(hr.nodes[0:index], hr.nodes[index+1:len(hr.nodes)]...)
			sort.Sort(hr.nodes)
		}
	}

	hr.mutex.Unlock()
}

func (hr *hashRing) Get(key interface{}) (server CanHash, err error) {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	length := len(hr.nodes)
	if length == 0 {
		return nil, ErrNoServer
	}

	if length < 2 {
		hr.nodes[0].hitCount.Inc()
		return hr.nodes[0].server, nil
	}

	value := HashValue(key)
	index := sort.Search(length, func(i int) bool {
		return hr.nodes[i].hashValue >= value
	})

	if index == length || index == 0 {
		if (value - hr.nodes[length-1].hashValue) < (math.MaxUint32 - value + hr.nodes[0].hashValue) {
			hr.nodes[length-1].hitCount.Inc()
			return hr.nodes[length-1].server, nil
		}

		hr.nodes[0].hitCount.Inc()
		return hr.nodes[0].server, nil
	}

	if (hr.nodes[index].hashValue - value) > (value - hr.nodes[index-1].hashValue) {
		hr.nodes[index-1].hitCount.Inc()
		return hr.nodes[index-1].server, nil
	}

	hr.nodes[index].hitCount.Inc()
	return hr.nodes[index].server, nil
}

func (hr *hashRing) Index(index int) (server CanHash, err error) {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	length := len(hr.nodes)
	if length == 0 {
		return nil, ErrNoServer
	}

	if length <= index {
		return nil, ErrNoServer
	}

	hr.nodes[index].hitCount.Inc()
	return hr.nodes[index].server, nil
}

func (hr *hashRing) Length() int {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()
	return len(hr.nodes)
}

func (hr *hashRing) Range(handler func(index int, server CanHash, hitCount uint64) (handled bool)) {
	hr.mutex.RLock()
	for index, _ := range hr.nodes {
		//标记已处理
		if handler(index, hr.nodes[index].server, hr.nodes[index].hitCount.Load()) {
			break
		}
	}
	hr.mutex.RUnlock()
}
