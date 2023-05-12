package retry

import (
	"context"
	"errors"
	"runtime"
	"time"
)

// Interval 执行时间间隔
var Interval = 10 * time.Millisecond

// Timeout 超时时间
var Timeout = 2 * time.Second

// MaxRetryCount 最大重试次数
var MaxRetryCount = 3

var ErrWaitTimeout = errors.New("timed out waiting for the condition")
var ErrEmptyFunc = errors.New("empty func")
var ErrExceedMaxLimit = errors.New("exceed max limit")

type WaitWithContextFunc func(ctx context.Context) <-chan struct{}
type ConditionWithContextFunc func(context.Context) (done bool, err error)

func InitRetryConfig(interval, timeout time.Duration, maxRetryCount int) {
	Interval, Timeout, MaxRetryCount = interval, timeout, maxRetryCount
}

func ExecFuncWithRetry(f func() (err error)) error {
	if f == nil {
		return ErrEmptyFunc
	}
	var lastErr error
	retryCount := 0

	wrapFunc := ConditionWithContextFunc(func(context.Context) (done bool, err error) {
		if retryCount >= MaxRetryCount {
			return true, ErrExceedMaxLimit
		}
		if err := f(); err != nil {
			lastErr = err
			// 如果想碰到指定错误时运行一次结束运行, 在这里判断err类型 返回 true,err 即可

			retryCount++
			return false, nil
		}

		return true, nil
	})

	if err := PollImmediate(Interval, Timeout, wrapFunc); err != nil {
		// 如果是runtime error 这里返回这个
		if lastErr == nil {
			return err
		}
		return lastErr
	}

	return nil
}

func PollImmediate(interval, timeout time.Duration, condition ConditionWithContextFunc) error {
	ctx := context.Background()
	done, err := ExecFunc(ctx, condition)
	if err != nil {
		return err
	}
	if done {
		return nil
	}

	wait := BuildTickerChanFunc(interval, timeout)
	select {
	case <-ctx.Done():
		return ErrWaitTimeout
	default:
		return WaitForWithContext(ctx, wait, condition)
	}

}

func WaitForWithContext(ctx context.Context, wait WaitWithContextFunc, f ConditionWithContextFunc) error {
	waitCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := wait(waitCtx)
	for {
		select {
		case _, open := <-c:
			ok, err := ExecFunc(ctx, f)
			if err != nil {
				return err
			}
			if ok {
				return nil
			}
			if !open {
				return ErrWaitTimeout
			}
		case <-ctx.Done():
			// 父ctx 关闭了
			return ErrWaitTimeout
		}
	}
}

// ExecFunc 执行方法,防止运行时错误
func ExecFunc(ctx context.Context, f ConditionWithContextFunc) (ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case runtime.Error:
				err = r.(error)
			default:
				err = errors.New("unknown error")
			}
		}
	}()
	return f(ctx)
}

// BuildTickerChanFunc 构建chan 不设置timeout(timeout = 0)得话 会一直调用 直到父context 关闭
func BuildTickerChanFunc(interval, timeout time.Duration) WaitWithContextFunc {
	return WaitWithContextFunc(func(ctx context.Context) <-chan struct{} {
		ch := make(chan struct{})

		go func() {
			defer close(ch)

			tick := time.NewTicker(interval)
			defer tick.Stop()

			var after <-chan time.Time
			if timeout != 0 {
				timer := time.NewTicker(timeout)
				after = timer.C
				defer timer.Stop()
			}

			for {
				select {
				case <-tick.C:
					ch <- struct{}{}
				case <-after:
					return
				case <-ctx.Done():
					return
				}
			}
		}()

		return ch
	})
}
