package limiter

import (
	"context"
	"github.com/go-redis/redis/v8"
	"io"
	"os"
	"sync"
	"time"
)

var leakyBucketLua = ""
var dlOnce sync.Once

func init() {
	dlOnce.Do(func() {
		f, err := os.Open("./d_leaky_bucket.lua")
		if err != nil {
			panic(err)
		}
		bs, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		leakyBucketLua = string(bs)
	})
}

func NewDLeakyBucket(redisKey string, rate, peak int, redisClient *redis.Client) *DLeakyBucket {
	return &DLeakyBucket{redisKey: redisKey, rate: rate, peak: peak, redisClient: redisClient}
}

type DLeakyBucket struct {
	redisKey    string
	rate        int
	peak        int
	redisClient *redis.Client
}

func (dl *DLeakyBucket) Take() bool {
	// 执行lua
	res, err := dl.redisClient.Eval(context.TODO(), leakyBucketLua, []string{dl.redisKey}, dl.peak, dl.rate, time.Now().Unix()).Result()
	if err != nil {
		return false
	}

	if code, ok := res.(int64); ok {
		return code == 1
	}
	return false
}
