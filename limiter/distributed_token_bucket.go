package limiter

import (
	"context"
	"github.com/go-redis/redis/v8"
	"io"
	"os"
	"sync"
	"time"
)

// args RedisKey 生成速度 容量 当前时间戳
var tokenBucketLua = ""
var dtOnce sync.Once

func init() {
	dtOnce.Do(func() {
		f, err := os.Open("./d_token_bucket.lua")
		if err != nil {
			panic(err)
		}
		bs, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		tokenBucketLua = string(bs)
	})
}

func NewDTokenBucket(redisKey string, rate, cap int, redisClient *redis.Client) *DTokenBucket {
	return &DTokenBucket{redisKey: redisKey, rate: rate, cap: cap, redisClient: redisClient}
}

type DTokenBucket struct {
	redisKey    string
	rate        int
	cap         int
	redisClient *redis.Client
}

func (dt *DTokenBucket) Take() bool {
	// 执行lua
	res, err := dt.redisClient.Eval(context.TODO(), tokenBucketLua, []string{dt.redisKey}, dt.rate, dt.cap, time.Now().Unix()).Result()
	if err != nil {
		return false
	}

	if code, ok := res.(int64); ok {
		return code == 1
	}
	return false
}
