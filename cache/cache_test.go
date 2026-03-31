package cache

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	c := New(Options{Capacity: 10, DefaultTTL: time.Minute})
	defer c.Close()

	c.Set("k", "v", 0)
	got, ok := c.Get("k")
	if !ok {
		t.Fatalf("expected cache hit")
	}
	if got.(string) != "v" {
		t.Fatalf("Get() = %v, want v", got)
	}
}

func TestTTLExpiration(t *testing.T) {
	c := New(Options{Capacity: 10, DefaultTTL: time.Millisecond * 30})
	defer c.Close()

	c.Set("k", "v", 0)
	time.Sleep(time.Millisecond * 50)
	if _, ok := c.Get("k"); ok {
		t.Fatalf("expected cache miss after ttl")
	}
}

func TestLRUEviction(t *testing.T) {
	c := New(Options{Capacity: 2, DefaultTTL: time.Minute})
	defer c.Close()

	c.Set("a", 1, 0)
	c.Set("b", 2, 0)
	if _, ok := c.Get("a"); !ok {
		t.Fatalf("expected a to exist")
	}
	c.Set("c", 3, 0)

	if _, ok := c.Get("b"); ok {
		t.Fatalf("expected b to be evicted as LRU")
	}
	if _, ok := c.Get("a"); !ok {
		t.Fatalf("expected a to remain")
	}
	if _, ok := c.Get("c"); !ok {
		t.Fatalf("expected c to remain")
	}
}

func TestGetOrLoadSingleflight(t *testing.T) {
	c := New(Options{Capacity: 20, DefaultTTL: time.Minute})
	defer c.Close()

	var loaderCount int32
	loader := func(ctx context.Context) (any, error) {
		atomic.AddInt32(&loaderCount, 1)
		time.Sleep(time.Millisecond * 20)
		return "loaded", nil
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err := c.GetOrLoad(context.Background(), "x", time.Minute, loader)
			if err != nil {
				t.Errorf("GetOrLoad() err = %v", err)
				return
			}
			if v.(string) != "loaded" {
				t.Errorf("GetOrLoad() = %v, want loaded", v)
			}
		}()
	}
	wg.Wait()

	if got := atomic.LoadInt32(&loaderCount); got != 1 {
		t.Fatalf("loaderCount = %d, want 1", got)
	}
}
