package limiter

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"testing"
	"time"
)

func TestNewDLeakyBucketValidation(t *testing.T) {
	cli := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	if _, err := NewDLeakyBucket("", 1, 4, cli); err == nil {
		t.Fatalf("expected error when redisKey is empty")
	}
	if _, err := NewDLeakyBucket("bucket", 0, 4, cli); err == nil {
		t.Fatalf("expected error when rate <= 0")
	}
	if _, err := NewDLeakyBucket("bucket", 1, 0, cli); err == nil {
		t.Fatalf("expected error when peak <= 0")
	}
	if _, err := NewDLeakyBucket("bucket", 1, 4, nil); err == nil {
		t.Fatalf("expected error when redisClient is nil")
	}
}

func TestNewDLeakyBucket(t *testing.T) {
	if os.Getenv("COMMON_RUN_REDIS_INTEGRATION") != "1" {
		t.Skip("set COMMON_RUN_REDIS_INTEGRATION=1 to run redis integration tests")
	}

	addr := os.Getenv("COMMON_REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	password := os.Getenv("COMMON_REDIS_PASSWORD")

	cli := redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: 0})
	s, err := cli.Ping(context.TODO()).Result()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(s)

	dt, err := NewDLeakyBucket("dl_bucket_test1", 2, 4, cli)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		t.Log(time.Now().Format("2006-01-02 15:04:05"), dt.Take())
		time.Sleep(time.Second / 4)
	}
}
