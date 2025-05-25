package main

import (
    "errors"
)

type ValueType int

const (
    StringType ValueType = iota
    ListType
)

type Value struct {
    Type ValueType
    Str  string
    List []string
}

func (s *Store) LPush(key string, values ...string) int {
    s.mu.Lock()
    defer s.mu.Unlock()

    v, _ := s.getOrInitList(key)
    // FIX: Insert from left to right, NOT reversed
    for i := 0; i < len(values); i++ {
        v.List = append([]string{values[i]}, v.List...)
    }
    return len(v.List)
}

func (s *Store) RPop(key string) (string, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    val, ok := s.data[key]
    if !ok || val.Type != ListType || len(val.List) == 0 {
        return "", errors.New("no such key or not a list or list empty")
    }
    item := val.List[len(val.List)-1]
    val.List = val.List[:len(val.List)-1]
    return item, nil
}

func (s *Store) LLen(key string) int {
    s.mu.RLock()
    defer s.mu.RUnlock()

    val, ok := s.data[key]
    if !ok || val.Type != ListType {
        return 0
    }
    return len(val.List)
}

func (s *Store) getOrInitList(key string) (*Value, bool) {
    val, ok := s.data[key]
    if ok && val.Type == ListType {
        return val, true
    }
    newList := &Value{Type: ListType, List: []string{}}
    s.data[key] = newList
    return newList, false
}
