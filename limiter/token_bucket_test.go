package limiter

import (
	"testing"
	"time"
)

func TestNewTokenBucket(t *testing.T) {
	l := NewTokenBucket(2, 10)
	for i := 0; i < 10; i++ {
		t.Log(time.Now().Format("2006-01-02 15:04:05"), l.Take())
		time.Sleep(time.Second / 4)
	}
}
