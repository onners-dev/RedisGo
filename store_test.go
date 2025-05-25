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
func TestTTL(t *testing.T) {
	store := NewStore()
	if ttl := store.TTL("foo"); ttl != -2 {
		t.Fatalf("expected -2 for non-existent key, got %d", ttl)
	}
	store.Set("foo", "bar")
	if ttl := store.TTL("foo"); ttl != -1 {
		t.Fatalf("expected -1 for no-expiry key, got %d", ttl)
	}
	store.Expire("foo", 1)
	ttl := store.TTL("foo")
	if ttl <= 0 {
		t.Fatalf("expected >0 for key with expiry, got %d", ttl)
	}
	time.Sleep(2 * time.Second)
	if ttl := store.TTL("foo"); ttl != -2 {
		t.Fatalf("expected -2 for expired key, got %d", ttl)
	}
}

func TestIncrDecr(t *testing.T) {
	store := NewStore()
	n, err := store.Incr("foo")
	if err != nil || n != 1 {
		t.Fatalf("INCR on new key should return 1, got %d, err=%v", n, err)
	}
	n, err = store.Incr("foo")
	if err != nil || n != 2 {
		t.Fatalf("INCR again should return 2, got %d, err=%v", n, err)
	}
	n, err = store.Decr("foo")
	if err != nil || n != 1 {
		t.Fatalf("DECR should return 1, got %d, err=%v", n, err)
	}
	store.Set("bar", "notanint")
	_, err = store.Incr("bar")
	if err == nil {
		t.Fatal("INCR should error on non-integer value")
	}
}

func TestMSetMGet(t *testing.T) {
	store := NewStore()
	err := store.MSet("a", "1", "b", "2", "c", "3")
	if err != nil {
		t.Fatal("unexpected error from MSet:", err)
	}
	vals := store.MGet("a", "c", "d")
	if vals[0] != "1" || vals[1] != "3" || vals[2] != "" {
		t.Fatalf("expected [1 3 ], got %#v", vals)
	}
	err = store.MSet("onlykey")
	if err == nil {
		t.Fatal("expected error on odd number of MSet args")
	}
}

func TestListOps(t *testing.T) {
    store := NewStore()

    // Test LPUSH + LLEN
    if n := store.LPush("mylist", "a"); n != 1 {
        t.Fatalf("expected 1, got %d", n)
    }
    if n := store.LPush("mylist", "b", "c"); n != 3 {
        t.Fatalf("expected 3 after multi-LPUSH, got %d", n)
    }
    if l := store.LLen("mylist"); l != 3 {
        t.Fatalf("expected 3, got %d", l)
    }

    // Test RPOP order
    v, err := store.RPop("mylist")
    if err != nil || v != "a" {
        t.Fatalf("expected 'a', got '%s' (err=%v)", v, err)
    }
    v, err = store.RPop("mylist")
    if err != nil || v != "b" {
        t.Fatalf("expected 'b', got '%s' (err=%v)", v, err)
    }
    v, err = store.RPop("mylist")
    if err != nil || v != "c" {
        t.Fatalf("expected 'c', got '%s' (err=%v)", v, err)
    }

    // Empty pop
    _, err = store.RPop("mylist")
    if err == nil {
        t.Fatal("expected error on empty list")
    }
}