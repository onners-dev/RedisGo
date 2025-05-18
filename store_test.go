package main

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	store := NewStore()
	store.Set("foo", "bar")
	val, ok := store.Get("foo")
	if !ok || val != "bar" {
		t.Fatalf("expected to get bar, got %v (ok=%v)", val, ok)
	}
}

func TestDel(t *testing.T) {
	store := NewStore()
	store.Set("foo", "bar")
	deleted := store.Del("foo")
	if !deleted {
		t.Fatal("expected to delete key foo")
	}
	_, ok := store.Get("foo")
	if ok {
		t.Fatal("expected key foo to be deleted")
	}
}

func TestExpire(t *testing.T) {
	store := NewStore()
	store.Set("foo", "bar")
	ok := store.Expire("foo", 1)
	if !ok {
		t.Fatal("expected to set expiry")
	}
	time.Sleep(2 * time.Second)
	val, has := store.Get("foo")
	if has {
		t.Fatalf("expected foo to expire, got value: %v", val)
	}
}

func TestKeys(t *testing.T) {
	store := NewStore()
	store.Set("a", "1")
	store.Set("b", "2")
	keys := store.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	for _, k := range keys {
		if k != "a" && k != "b" {
			t.Fatalf("unexpected key: %v", k)
		}
	}
}

func TestDumpAll(t *testing.T) {
	store := NewStore()
	store.Set("x", "alpha")
	store.Set("y", "beta")
	all := store.DumpAll()
	if len(all) != 2 {
		t.Fatalf("expected 2 keys in dump, got %d", len(all))
	}
	if all["x"] != "alpha" || all["y"] != "beta" {
		t.Fatalf("values mismatch in dump: %v", all)
	}
}

func TestExpireDoesNotAffectOtherKeys(t *testing.T) {
	store := NewStore()
	store.Set("foo", "bar")
	store.Set("baz", "qux")
	store.Expire("foo", 1)
	time.Sleep(2 * time.Second)
	_, ok := store.Get("foo")
	if ok {
		t.Fatal("expected foo to be expired")
	}
	val, ok := store.Get("baz")
	if !ok || val != "qux" {
		t.Fatal("expected baz to still exist")
	}
}
