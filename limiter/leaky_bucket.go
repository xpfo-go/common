package limiter

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	peak     int64 // 最高水位
	cur      int64 // 当前水位
	rate     int64 // 水流速度
	lastTime int64 // 上次放水时间戳
	now      func() time.Time
	mu       sync.Mutex
}

// NewLeakyBucket 限流器 漏桶算法
func NewLeakyBucket(rate, peak int64) *LeakyBucket {
	return &LeakyBucket{
		peak:     peak,
		cur:      0,
		rate:     rate,
		lastTime: time.Now().UnixNano(),
		now:      time.Now,
	}
}

func (l *LeakyBucket) Take() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now().UnixNano()
	elapsed := now - l.lastTime
	l.cur = l.cur - elapsed*l.rate/int64(time.Second)
	if l.cur < 0 {
		l.cur = 0
	}
	l.lastTime = now

	if l.cur >= l.peak {
		return false
	}

	l.cur++

	return true
}
