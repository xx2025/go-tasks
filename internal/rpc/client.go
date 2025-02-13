package rpc

import (
	"sync"
)

var clientPoolMap = make(map[string]*GRPCPool)
var mu sync.Mutex

func GetPool(host string) (*GRPCPool, error) {
	mu.Lock()
	defer mu.Unlock()
	pool, ok := clientPoolMap[host]
	if !ok {
		clientPoolMap[host] = newGRPCPool(host)
		pool = clientPoolMap[host]
	}
	return pool, nil
}
