package limiter

import (
	"testing"
	"time"
)

func TestLeakyBucketTakeDeterministic(t *testing.T) {
	clock := newFakeClock()
	bucket := NewLeakyBucket(1, 3)
	bucket.now = clock.Now
	bucket.lastTime = clock.Now().UnixNano()

	for i := 0; i < 3; i++ {
		if !bucket.Take() {
			t.Fatalf("expected take to succeed before peak, index=%d", i)
		}
	}
	if bucket.Take() {
		t.Fatalf("expected overflow take to fail when bucket reaches peak")
	}

	clock.Advance(1 * time.Second)
	if !bucket.Take() {
		t.Fatalf("expected one slot available after leaking for 1s")
	}
	if bucket.Take() {
		t.Fatalf("expected bucket to be full again after immediate extra take")
	}

	clock.Advance(3 * time.Second)
	for i := 0; i < 3; i++ {
		if !bucket.Take() {
			t.Fatalf("expected leak to clear bucket, index=%d", i)
		}
	}
	if bucket.Take() {
		t.Fatalf("expected overflow take to fail after re-filling to peak")
	}
}
