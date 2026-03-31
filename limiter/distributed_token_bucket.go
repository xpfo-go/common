package limiter

import (
	"context"
	_ "embed"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

//go:embed d_token_bucket.lua
var tokenBucketLua string

func NewDTokenBucket(redisKey string, rate, cap int, redisClient *redis.Client) (*DTokenBucket, error) {
	if redisKey == "" {
		return nil, errors.New("redisKey cannot be empty")
	}
	if rate <= 0 {
		return nil, errors.New("rate must be greater than 0")
	}
	if cap <= 0 {
		return nil, errors.New("cap must be greater than 0")
	}
	if redisClient == nil {
		return nil, errors.New("redisClient cannot be nil")
	}
	if tokenBucketLua == "" {
		return nil, errors.New("token bucket lua script is empty")
	}

	return &DTokenBucket{
		redisKey:    redisKey,
		rate:        rate,
		cap:         cap,
		redisClient: redisClient,
	}, nil
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
