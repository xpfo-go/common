package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	rate     int64 // 速率
	cap      int64 // 容量
	tokens   int64 // 当前桶里得令牌数
	lastTime int64 // 上次添加令牌得时间戳
	mu       sync.Mutex
}

// NewTokenBucket 限流器 令牌桶算法
func NewTokenBucket(rate, cap int64) *TokenBucket {
	return &TokenBucket{rate: rate, cap: cap, tokens: 0, lastTime: time.Now().Unix()}
}

// Take 拿令牌
func (t *TokenBucket) Take() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now().Unix()
	// 添加令牌
	t.tokens = t.tokens + (now-t.lastTime)*t.rate
	if t.tokens > t.cap {
		t.tokens = t.cap
	}
	t.lastTime = now

	if t.tokens > 0 {
		t.tokens--
		return true
	}
	return false
}
