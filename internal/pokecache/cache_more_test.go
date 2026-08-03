package pokecache

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

// A good cache reports (nil, false) for a key it never stored — it must not
// invent a value for a miss.
func TestGetMissingKey(t *testing.T) {
	cache := NewCache(time.Minute)

	val, ok := cache.Get("https://example.com/never-added")
	if ok {
		t.Errorf("expected ok=false for a key that was never added")
	}
	if val != nil {
		t.Errorf("expected a nil value for a missing key, got %q", val)
	}
}

// Adding the same key twice should overwrite, so the most recent value wins.
func TestAddOverwrites(t *testing.T) {
	cache := NewCache(time.Minute)

	cache.Add("k", []byte("old"))
	cache.Add("k", []byte("new"))

	val, ok := cache.Get("k")
	if !ok {
		t.Fatalf("expected to find key after overwrite")
	}
	if string(val) != "new" {
		t.Errorf("expected overwritten value %q, got %q", "new", string(val))
	}
}

// The cache must not reap entries before their TTL elapses.
func TestNotReapedBeforeTTL(t *testing.T) {
	const ttl = 50 * time.Millisecond
	cache := NewCache(ttl)
	cache.Add("k", []byte("v"))

	time.Sleep(ttl / 5) // comfortably before expiry
	if _, ok := cache.Get("k"); !ok {
		t.Errorf("entry was reaped before its TTL elapsed")
	}
}

// A cache shared across goroutines (as it is in this program) must be safe
// under concurrent access. Run `go test -race ./...` to actually detect races.
func TestConcurrentAccess(t *testing.T) {
	cache := NewCache(time.Minute)

	const workers = 50
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(n int) {
			defer wg.Done()
			key := "key-" + strconv.Itoa(n)
			cache.Add(key, []byte(strconv.Itoa(n)))
			cache.Get(key)
		}(i)
	}
	wg.Wait()
}
