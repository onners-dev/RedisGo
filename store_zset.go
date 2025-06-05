package main

import (
	"errors"
	"sort"
)

// ZSetEntry represents a member of a sorted set.
type ZSetEntry struct {
	Member string
	Score  float64
}

// Synchronize ZSetMap pointers after sorting/mutations.
func zsetMapSync(zs []ZSetEntry, zm map[string]*ZSetEntry) {
	for i := range zs {
		zm[zs[i].Member] = &zs[i]
	}
}

func (s *Store) ZAdd(key string, score float64, member string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, _ := s.getOrInitZSet(key)
	entry, exists := v.ZSetMap[member]
	if exists {
		// Update score if changed, then re-sort and resync
		if entry.Score != score {
			entry.Score = score
			sort.Slice(v.ZSet, func(i, j int) bool {
				if v.ZSet[i].Score == v.ZSet[j].Score {
					return v.ZSet[i].Member < v.ZSet[j].Member // lex order
				}
				return v.ZSet[i].Score < v.ZSet[j].Score
			})
			zsetMapSync(v.ZSet, v.ZSetMap)
		}
		return 0 // not a new member
	}
	// Add new member
	newEntry := ZSetEntry{Member: member, Score: score}
	v.ZSet = append(v.ZSet, newEntry)
	// Resort after insertion and resync pointers
	sort.Slice(v.ZSet, func(i, j int) bool {
		if v.ZSet[i].Score == v.ZSet[j].Score {
			return v.ZSet[i].Member < v.ZSet[j].Member
		}
		return v.ZSet[i].Score < v.ZSet[j].Score
	})
	zsetMapSync(v.ZSet, v.ZSetMap)
	return 1
}

func (s *Store) ZRem(key string, member string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.data[key]
	if !ok || val.Type != ZSetType {
		return 0
	}
	_, exists := val.ZSetMap[member]
	if !exists {
		return 0
	}
	// Remove from map and slice
	delete(val.ZSetMap, member)
	for i, entry := range val.ZSet {
		if entry.Member == member {
			val.ZSet = append(val.ZSet[:i], val.ZSet[i+1:]...)
			break
		}
	}
	zsetMapSync(val.ZSet, val.ZSetMap)
	return 1
}

func (s *Store) ZRange(key string, start, stop int) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]
	if !ok || val.Type != ZSetType {
		return nil, errors.New("no such key or not a zset")
	}
	n := len(val.ZSet)
	if n == 0 {
		return []string{}, nil
	}
	// handle negative indexes
	if start < 0 {
		start = n + start
	}
	if stop < 0 {
		stop = n + stop
	}
	if start < 0 {
		start = 0
	}
	if stop >= n {
		stop = n - 1
	}
	if start > stop || start >= n {
		return []string{}, nil
	}
	res := make([]string, 0, stop-start+1)
	for i := start; i <= stop; i++ {
		res = append(res, val.ZSet[i].Member)
	}
	return res, nil
}

func (s *Store) getOrInitZSet(key string) (*Value, bool) {
	val, ok := s.data[key]
	if ok && val.Type == ZSetType {
		return val, true
	}
	newZSet := &Value{
		Type:    ZSetType,
		ZSet:    []ZSetEntry{},
		ZSetMap: make(map[string]*ZSetEntry),
	}
	s.data[key] = newZSet
	return newZSet, false
}
