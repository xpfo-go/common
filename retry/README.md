# Retry

## 一 介绍

### 一个实现重试机制得组件

#### retry 包中设置了比较合理得 interval,timeout,maxRetryCount。如需修改可以 复制代码 或者 Fork本仓库。

``` go
// 执行时间间隔
const interval = 10 * time.Millisecond

// 超时时间
const timeout = 2 * time.Second

// 最大重试次数
const maxRetryCount = 3
```

## 二 使用

``` go
package main

import "github.com/xpfo-go/common/retry"

func main(){
    if err := retry.ExecFuncWithRetry(func() (err error) {
        // 调用方法
        err = http.request(url)
        return err
    }); err ! = nil {
    
    }
}
```