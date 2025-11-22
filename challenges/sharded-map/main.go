package main

import (
	"hash/fnv"
	"sync"
	"sync/atomic"
)

type ShardedMap struct {
	shards     []map[string]int
	locks      []*sync.RWMutex
	totalCount *atomic.Int32
}

func NewShardedMap(numShards int) *ShardedMap {
	shards := make([]map[string]int, 0, numShards)
	locks := make([]*sync.RWMutex, 0, numShards)

	for range numShards {
		shards = append(shards, make(map[string]int))
		locks = append(locks, &sync.RWMutex{})
	}

	return &ShardedMap{
		shards:     shards,
		locks:      locks,
		totalCount: &atomic.Int32{},
	}
}

func (m *ShardedMap) Get(key string) (int, bool) {
	n := shard(key, len(m.shards))

	s := m.shards[n]
	l := m.locks[n]

	l.RLock()
	defer l.RUnlock()

	val, found := s[key]

	return val, found
}

func (m *ShardedMap) Put(key string, value int) {
	n := shard(key, len(m.shards))

	s := m.shards[n]
	l := m.locks[n]

	l.Lock()
	defer l.Unlock()

	_, found := s[key]
	if !found {
		m.totalCount.Add(1)
	}

	s[key] = value
}

func (m *ShardedMap) Delete(key string) {
	n := shard(key, len(m.shards))

	s := m.shards[n]
	l := m.locks[n]

	l.Lock()
	defer l.Unlock()

	_, found := s[key]
	if found {
		delete(s, key)
		m.totalCount.Add(-1)
	}
}

func (m *ShardedMap) Len() int {
	return int(m.totalCount.Load())
}

func shard(key string, numShards int) int {
	h := fnv.New32()
	h.Write([]byte(key))
	v := h.Sum32()

	return int(v) % numShards
}
