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
	mu       sync.Mutex
}

// NewLeakyBucket 限流器 漏桶算法
func NewLeakyBucket(rate, peak int64) *LeakyBucket {
	return &LeakyBucket{peak: peak, cur: 0, rate: rate, lastTime: time.Now().Unix()}
}

func (l *LeakyBucket) Take() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().Unix()
	l.cur = l.cur - (now-l.lastTime)*l.rate
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
