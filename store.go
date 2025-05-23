package main

import (
	"errors"
	"strconv"
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

// TTL returns the remaining time to live of a key in seconds. OR:
// -2 means that the key does not exist.
// -1 means that the key exists but has no expiry set.
func (s *Store) TTL(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[key]
	if !exists {
		return -2
	}
	exp, hasExp := s.expires[key]
	if !hasExp {
		return -1
	}
	ttl := int(time.Until(exp).Seconds())
	if ttl < 0 {
		return -2
	}
	return ttl
}

// Incr increments a key's integer value by 1, setting it to 0 if it doesn't exist.
func (s *Store) Incr(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[key]
	var n int
	var err error
	if ok {
		n, err = strconv.Atoi(val)
		if err != nil {
			return 0, errors.New("value is not an integer")
		}
	} else {
		n = 0
	}
	n++
	s.data[key] = strconv.Itoa(n)
	delete(s.expires, key)
	return n, nil
}

// Decr decrements a key's integer value by 1, setting it to 0 if it doesn't exist.
func (s *Store) Decr(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[key]
	var n int
	var err error
	if ok {
		n, err = strconv.Atoi(val)
		if err != nil {
			return 0, errors.New("value is not an integer")
		}
	} else {
		n = 0
	}
	n--
	s.data[key] = strconv.Itoa(n)
	delete(s.expires, key)
	return n, nil
}

// MSet sets multiple key-value pairs. Expects even number of args: key1, val1, key2, etc
func (s *Store) MSet(keysValues ...string) error {
	if len(keysValues)%2 != 0 {
		return errors.New("MSET requires an even number of arguments")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < len(keysValues); i += 2 {
		s.data[keysValues[i]] = keysValues[i+1]
		delete(s.expires, keysValues[i]) // Clear expiry on update
	}
	return nil
}

// MGet returns values for the given keys in order. Missing keys return "".
func (s *Store) MGet(keys ...string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	values := make([]string, len(keys))
	now := time.Now()
	for i, key := range keys {
		if exp, ok := s.expires[key]; ok && now.After(exp) {
			values[i] = ""
			continue
		}
		val, ok := s.data[key]
		if ok {
			values[i] = val
		} else {
			values[i] = ""
		}
	}
	return values
}
