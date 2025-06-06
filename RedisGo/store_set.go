package main

import (
	"errors"
)

// SAdd adds one or more members to a set.
// Returns the number of new elements actually added (not previously present).
func (s *Store) SAdd(key string, members ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, _ := s.getOrInitSet(key)
	added := 0
	for _, m := range members {
		if _, exists := v.Set[m]; !exists {
			v.Set[m] = struct{}{}
			added++
		}
	}
	return added
}

// SRem removes one or more members from the set.
// Returns the number of members actually removed.
func (s *Store) SRem(key string, members ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[key]
	if !ok || val.Type != SetType {
		return 0
	}
	removed := 0
	for _, m := range members {
		if _, exists := val.Set[m]; exists {
			delete(val.Set, m)
			removed++
		}
	}
	return removed
}

// SMembers returns a slice of all members in the set.
// If key doesn't exist or isn't a set, returns an error.
func (s *Store) SMembers(key string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	if !ok || val.Type != SetType {
		return nil, errors.New("no such key or not a set")
	}
	members := make([]string, 0, len(val.Set))
	for m := range val.Set {
		members = append(members, m)
	}
	return members, nil
}

// getOrInitSet retrieves a set for a key, or creates one if missing/wrong type.
func (s *Store) getOrInitSet(key string) (*Value, bool) {
	val, ok := s.data[key]
	if ok && val.Type == SetType {
		return val, true
	}
	newSet := &Value{Type: SetType, Set: make(map[string]struct{})}
	s.data[key] = newSet
	return newSet, false
}
