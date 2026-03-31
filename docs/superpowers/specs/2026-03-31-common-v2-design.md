# common v2.0.0 设计说明

## 目标
将 `github.com/xpfo-go/common/v2` 升级为可稳定被外部项目依赖的 v2 版本，重点解决运行时 panic、错误不可观测、测试不稳定、仓库规范缺失的问题。

## 主要改造范围

### 1. limiter（分布式令牌桶/漏桶）
- 现状问题
  - `init()` 使用相对路径读取 Lua 脚本，导入包时在不同工作目录可能 panic。
  - 构造函数不返回错误，调用方无法感知脚本或参数错误。
- v2 设计
  - 使用 `go:embed` 内嵌 Lua 文件，移除启动时文件 IO 依赖。
  - 新增参数校验和 `error` 返回：
    - `NewDTokenBucket(...) (*DTokenBucket, error)`
    - `NewDLeakyBucket(...) (*DLeakyBucket, error)`
  - 保留 `Take() bool` 行为，避免扩大 API 破坏面。

### 2. rsa（加解密）
- 现状问题
  - 大量忽略错误，错误路径会 panic 或返回无意义结果。
- v2 设计
  - API 改为显式错误返回：
    - `Generate(bits int) error`
    - `Encrypt(plainText []byte, publicKeyPath string) ([]byte, error)`
    - `Decrypt(cipherText []byte, privateKeyPath string) ([]byte, error)`
  - 对 PEM 解析、类型断言、文件读取、加解密失败都返回明确错误。

### 3. 测试策略
- 将当前依赖本地 Redis 的测试改为“可选集成测试”。
- 默认 `go test ./...` 不依赖外部 Redis。
- Redis 集成测试通过环境变量开启：`COMMON_RUN_REDIS_INTEGRATION=1`。

### 4. CI/CD 与 GitHub 规范
- 增加 GitHub Actions：
  - `ci.yml`：`go test ./...` + `go vet ./...` + `golangci-lint`
  - `release.yml`：tag 推送时自动发布 release（附 changelog）
- 增加仓库规范文件：
  - `README.md`（完善安装、模块说明、快速开始、兼容性说明）
  - `CHANGELOG.md`（Keep a Changelog 风格）
  - `CONTRIBUTING.md`
  - `SECURITY.md`
  - Issue/PR 模板

## Breaking Changes
1. `limiter.NewDTokenBucket` 和 `limiter.NewDLeakyBucket` 现在返回 `error`。
2. `rsa.Generate/Encrypt/Decrypt` 现在返回 `error`。
3. Redis 集成测试默认不会执行，需显式开启。

## 兼容性和版本策略
- 模块路径升级为 `github.com/xpfo-go/common/v2`，符合 Go module 的大版本语义化导入要求。
- 通过语义化版本发布 `v2.0.0`，在 release note 明确标注 breaking changes 和迁移示例。

## 验收标准
- `go test ./...` 在无 Redis 环境稳定通过。
- `go vet ./...` 通过。
- CI 在 PR 和 `master` 推送自动运行。
- 仓库文档与模板齐全。
- 成功创建 `v2.0.0` tag 和 GitHub Release。
