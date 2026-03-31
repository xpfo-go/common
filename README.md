# common

[![CI](https://github.com/xpfo-go/common/actions/workflows/ci.yml/badge.svg)](https://github.com/xpfo-go/common/actions/workflows/ci.yml)
[![Release](https://github.com/xpfo-go/common/actions/workflows/release.yml/badge.svg)](https://github.com/xpfo-go/common/actions/workflows/release.yml)

`common` 是一个 Go 工具组件集合，包含加密、限流、数据结构、日志、重试与 WebSocket 管理等模块。

## Installation

> v2 使用 Go Modules 语义化导入路径（`/v2`）。

```bash
go get github.com/xpfo-go/common/v2@latest
```

## Quick Start

```go
import (
    "github.com/xpfo-go/common/v2/limiter"
    rsautil "github.com/xpfo-go/common/v2/rsa"
)

func example() error {
    tb := limiter.NewTokenBucket(10, 100)
    _ = tb.Take()

    if err := rsautil.Generate(2048); err != nil {
        return err
    }
    return nil
}
```

## Modules

| Module | Description | Doc |
| --- | --- | --- |
| `aes` | AES 加解密工具 | `aes/aes_test.go` |
| `gis` | GIS 相关工具 | [gis/README.md](gis/README.md) |
| `heap` | 堆结构 | [heap/README.md](heap/README.md) |
| `limiter` | 本地/分布式限流器（令牌桶、漏桶） | `limiter/*_test.go` |
| `logs` | zap + lumberjack 日志封装 | `logs/logs_test.go` |
| `queue` | 队列实现 | [queue/README.md](queue/README.md) |
| `retry` | 重试工具 | [retry/README.md](retry/README.md) |
| `rsa` | RSA 密钥生成与加解密 | `rsa/rsa_test.go` |
| `stack` | 栈结构 | [stack/README.md](stack/README.md) |
| `websockets` | WebSocket 连接管理 | [websockets/README.md](websockets/README.md) |

## Breaking Changes in v2.0.0

1. `limiter.NewDTokenBucket` 现在返回 `(*DTokenBucket, error)`。
2. `limiter.NewDLeakyBucket` 现在返回 `(*DLeakyBucket, error)`。
3. `rsa.Generate` 现在返回 `error`。
4. `rsa.Encrypt` 现在返回 `([]byte, error)`。
5. `rsa.Decrypt` 现在返回 `([]byte, error)`。

### Migration Example

```go
// v1
// b := rsa.Encrypt(data, "./public_key.pem")

// v2
b, err := rsa.Encrypt(data, "./public_key.pem")
if err != nil {
    return err
}
```

## Test

```bash
go test ./...
go vet ./...
```

Redis 集成测试默认关闭，开启方式：

```bash
COMMON_RUN_REDIS_INTEGRATION=1 \
COMMON_REDIS_ADDR=127.0.0.1:6379 \
COMMON_REDIS_PASSWORD= \
go test ./limiter -run TestNewD
```

## Contributing

请先阅读 [CONTRIBUTING.md](CONTRIBUTING.md)。

## Security

安全问题请参考 [SECURITY.md](SECURITY.md)。

## License

本项目遵循仓库根目录中的 License。
