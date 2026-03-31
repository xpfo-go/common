package limiter

import (
	"testing"
	"time"
)

type fakeClock struct {
	current time.Time
}

func newFakeClock() *fakeClock {
	return &fakeClock{current: time.Unix(0, 0)}
}

func (f *fakeClock) Now() time.Time {
	return f.current
}

func (f *fakeClock) Advance(d time.Duration) {
	f.current = f.current.Add(d)
}

func TestTokenBucketTakeDeterministic(t *testing.T) {
	clock := newFakeClock()
	bucket := NewTokenBucket(2, 3)
	bucket.now = clock.Now
	bucket.lastTime = clock.Now().UnixNano()

	if bucket.Take() {
		t.Fatalf("expected no token at initial state")
	}

	clock.Advance(500 * time.Millisecond)
	if !bucket.Take() {
		t.Fatalf("expected one token after 500ms with rate=2/s")
	}
	if bucket.Take() {
		t.Fatalf("expected no more token without time advance")
	}

	clock.Advance(2 * time.Second)
	for i := 0; i < 3; i++ {
		if !bucket.Take() {
			t.Fatalf("expected capped token to be available, index=%d", i)
		}
	}
	if bucket.Take() {
		t.Fatalf("expected token bucket to be empty after consuming cap")
	}
}
