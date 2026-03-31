package limiter

import (
	"context"
	_ "embed"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

//go:embed d_leaky_bucket.lua
var leakyBucketLua string

func NewDLeakyBucket(redisKey string, rate, peak int, redisClient *redis.Client) (*DLeakyBucket, error) {
	if redisKey == "" {
		return nil, errors.New("redisKey cannot be empty")
	}
	if rate <= 0 {
		return nil, errors.New("rate must be greater than 0")
	}
	if peak <= 0 {
		return nil, errors.New("peak must be greater than 0")
	}
	if redisClient == nil {
		return nil, errors.New("redisClient cannot be nil")
	}
	if leakyBucketLua == "" {
		return nil, errors.New("leaky bucket lua script is empty")
	}

	return &DLeakyBucket{
		redisKey:    redisKey,
		rate:        rate,
		peak:        peak,
		redisClient: redisClient,
	}, nil
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
