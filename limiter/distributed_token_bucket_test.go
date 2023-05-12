package limiter

import (
	"context"
	"github.com/go-redis/redis/v8"
	"io"
	"os"
	"testing"
	"time"
)

func TestLua(t *testing.T) {
	f, _ := os.Open("./d_token_bucket.lua")
	bs, _ := io.ReadAll(f)
	t.Log(string(bs))
}

func TestNewDTokenBucket(t *testing.T) {
	cli := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", Password: "123456", DB: 0})
	s, err := cli.Ping(context.TODO()).Result()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(s)

	dt := NewDTokenBucket("dt_bucket", 2, 4, cli)
	for i := 0; i < 20; i++ {
		t.Log(time.Now().Format("2006-01-02 15:04:05"), dt.Take())
		time.Sleep(time.Second / 4)
	}
}
