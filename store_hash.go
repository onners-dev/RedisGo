package main

import (
	"errors"
)

// HSet sets field in the hash stored at key to value.
// Returns 1 if field is a new field in the hash, 0 if field existed and was updated.
func (s *Store) HSet(key, field, value string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, _ := s.getOrInitHash(key)
	_, exists := v.Hash[field]
	v.Hash[field] = value
	if exists {
		return 0 // overwritten
	}
	return 1 // new field
}

// HGet gets the value of a field in the hash stored at key.
func (s *Store) HGet(key, field string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	if !ok || val.Type != HashType {
		return "", false
	}
	v, ok := val.Hash[field]
	return v, ok
}

// HDel removes fields from the hash stored at key.
// Returns the number of fields that were removed.
func (s *Store) HDel(key string, fields ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[key]
	if !ok || val.Type != HashType {
		return 0
	}
	deleted := 0
	for _, field := range fields {
		if _, present := val.Hash[field]; present {
			delete(val.Hash, field)
			deleted++
		}
	}
	return deleted
}

// HGetAll gets all field-value pairs in the hash stored at key.
func (s *Store) HGetAll(key string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	if !ok || val.Type != HashType {
		return nil, errors.New("no such key or not a hash")
	}
	// Return a copy to avoid race conditions.
	out := make(map[string]string, len(val.Hash))
	for f, v := range val.Hash {
		out[f] = v
	}
	return out, nil
}

// getOrInitHash gets or initializes a hash value at key.
func (s *Store) getOrInitHash(key string) (*Value, bool) {
	val, ok := s.data[key]
	if ok && val.Type == HashType {
		return val, true
	}
	newHash := &Value{Type: HashType, Hash: make(map[string]string)}
	s.data[key] = newHash
	return newHash, false
}
