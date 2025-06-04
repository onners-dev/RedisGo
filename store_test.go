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

func TestSetOps(t *testing.T) {
	store := NewStore()

	// SAdd single member
	added := store.SAdd("myset", "a")
	if added != 1 {
		t.Fatalf("expected 1 added, got %d", added)
	}
	// SAdd duplicate member does not add again
	added = store.SAdd("myset", "a")
	if added != 0 {
		t.Fatalf("expected 0 added (duplicate), got %d", added)
	}
	// SAdd multiple members, including a duplicate
	added = store.SAdd("myset", "b", "c", "a")
	if added != 2 {
		t.Fatalf("expected 2 added, got %d", added)
	}

	// SMembers should return all members (order not guaranteed)
	members, err := store.SMembers("myset")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	want := map[string]bool{"a": true, "b": true, "c": true}
	if len(members) != len(want) {
		t.Fatalf("expected 3 members, got %d: %v", len(members), members)
	}
	for _, m := range members {
		if !want[m] {
			t.Errorf("unexpected member in set: %q", m)
		}
	}

	// SRem removes members
	removed := store.SRem("myset", "a", "d") // 'a' exists, 'd' does not
	if removed != 1 {
		t.Fatalf("expected 1 removed, got %d", removed)
	}
	members, _ = store.SMembers("myset")
	want = map[string]bool{"b": true, "c": true}
	if len(members) != len(want) {
		t.Fatalf("expected 2 members after SRem, got %d: %v", len(members), members)
	}

	// SAdd on key with wrong type creates a new set
	store.Set("notaset", "hello")
	added = store.SAdd("notaset", "x")
	if added != 1 {
		t.Fatalf("expected 1 added for new set, got %d", added)
	}
	val, ok := store.Get("notaset")
	if ok && val == "hello" {
		t.Fatalf("SAdd should overwrite previous string value")
	}

	// SMembers error on non-existing set
	_, err = store.SMembers("doesnotexist")
	if err == nil {
		t.Fatalf("expected error on non-existent set key")
	}
}

func TestHashOps(t *testing.T) {
	store := NewStore()

	// HSET new field
	n := store.HSet("h", "foo", "bar")
	if n != 1 {
		t.Fatalf("expected 1 for new field, got %d", n)
	}
	// HSET existing field (overwrite)
	n = store.HSet("h", "foo", "baz")
	if n != 0 {
		t.Fatalf("expected 0 for overwrite, got %d", n)
	}
	// HGET
	v, ok := store.HGet("h", "foo")
	if !ok || v != "baz" {
		t.Fatalf("expected baz, got %q", v)
	}
	// HDEL
	del := store.HDel("h", "foo")
	if del != 1 {
		t.Fatalf("expected 1 deleted, got %d", del)
	}
	// HGETALL
	store.HSet("h", "a", "1")
	store.HSet("h", "b", "2")
	all, err := store.HGetAll("h")
	if err != nil || all["a"] != "1" || all["b"] != "2" {
		t.Fatalf("HGETALL failed: got %v err=%v", all, err)
	}
}
