package main

import (
	"sync"
	"time"
)

type Store struct {
	mu      sync.RWMutex
	data    map[string]string
	expires map[string]time.Time
}

func NewStore() *Store {
	s := &Store{
		data:    make(map[string]string),
		expires: make(map[string]time.Time),
	}
	go s.expiryLoop()
	return s
}

// Set stores a value for a given key.
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	delete(s.expires, key) // Remove expiry if value is updated
}

// Get retrieves the value for a given key and a boolean if it exists.
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if exp, ok := s.expires[key]; ok && time.Now().After(exp) {
		return "", false
	}
	val, ok := s.data[key]
	return val, ok
}

// Del removes a key from the store. Returns true if key was present.
func (s *Store) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, existed := s.data[key]
	delete(s.data, key)
	delete(s.expires, key)
	return existed
}

// Expire sets an expiration for a key in seconds.
func (s *Store) Expire(key string, seconds int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[key]
	if !ok {
		return false
	}
	s.expires[key] = time.Now().Add(time.Duration(seconds) * time.Second)
	return true
}

// expiryLoop runs in the background to remove expired keys.
func (s *Store) expiryLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for k, exp := range s.expires {
			if now.After(exp) {
				delete(s.data, k)
				delete(s.expires, k)
			}
		}
		s.mu.Unlock()
	}
}

// Keys returns a slice of all keys that are not expired.
func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var keys []string
	now := time.Now()
	for k := range s.data {
		if exp, ok := s.expires[k]; ok && now.After(exp) {
			continue
		}
		keys = append(keys, k)
	}
	return keys
}

// DumpAll returns a map of all key-value pairs that are not expired.
func (s *Store) DumpAll() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	now := time.Now()
	all := make(map[string]string)
	for k, v := range s.data {
		if exp, ok := s.expires[k]; ok && now.After(exp) {
			continue
		}
		all[k] = v
	}
	return all
}
