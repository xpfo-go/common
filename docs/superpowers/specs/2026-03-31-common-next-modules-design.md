# common v2 下一批常用模块设计

## 背景与目标

当前 `common` 已覆盖日志、限流、加密和数据结构，但业务工程中高频出现的配置加载、错误模型、缓存与 HTTP 客户端仍需重复开发。

目标是在不引入复杂框架绑定的前提下，增加 4 个可复用模块：
- `config`
- `xerrors`
- `cache`
- `httpx`

## 可选方案对比

### 方案 A：每个模块都做“全功能框架级实现”
- 优点：能力最全。
- 缺点：开发和维护成本高，API 面会快速膨胀，早期易过度设计。

### 方案 B：先做“最小可用核心 + 明确扩展点”（推荐）
- 优点：实现速度快，风险低，能覆盖大部分常见场景。
- 缺点：部分高级功能（如分布式 cache、多配置源优先级 DSL）后续再补。

### 方案 C：仅封装第三方库并暴露透传 API
- 优点：实现快。
- 缺点：业务强耦合第三方，升级/替换成本高，`common` 抽象价值偏低。

选择：方案 B。

## 模块设计

### 1) config

#### 目标
- 支持 JSON/YAML 文件加载。
- 支持 `default` 标签默认值。
- 支持环境变量覆盖（带 prefix）。
- 支持 `required:"true"` 校验。

#### API
- `LoadFile(path string, out any) error`
- `LoadJSON(path string, out any) error`
- `LoadYAML(path string, out any) error`
- `ApplyDefaults(cfg any) error`
- `LoadFromEnv(prefix string, cfg any) error`
- `ValidateRequired(cfg any) error`

#### 约束
- 仅支持导出字段。
- 优先支持基础类型（string/bool/int/uint/float/duration 和 []string）。

### 2) xerrors

#### 目标
- 统一错误码和包装风格。
- 支持 `errors.Is/As` 链路。
- 提供错误码到 HTTP 状态码映射。

#### API
- `New(code, message string) *Error`
- `Wrap(err error, code, message string) *Error`
- `IsCode(err error, code string) bool`
- `CodeOf(err error) string`
- `RegisterHTTPStatus(code string, status int)`
- `HTTPStatus(err error) int`

### 3) cache

#### 目标
- 内存缓存，支持 TTL + LRU。
- 并发安全。
- 支持 `GetOrLoad`（singleflight 防击穿）。

#### API
- `New(options Options) *Cache`
- `Set(key string, value any, ttl time.Duration)`
- `Get(key string) (any, bool)`
- `Delete(key string)`
- `Len() int`
- `GetOrLoad(ctx, key, ttl, loader)`
- `Close()`

### 4) httpx

#### 目标
- 统一超时、重试、JSON 编解码。
- 对 5xx/网络错误进行有限重试。

#### API
- `NewClient(options Options) *Client`
- `Do(ctx context.Context, req *http.Request) (*http.Response, error)`
- `JSON(ctx, method, url string, in any, headers map[string]string, out any) error`

## 测试与验收

### TDD 要求
- 每个模块先写失败测试，再写最小实现。
- 每个模块至少覆盖：成功路径 + 一个关键失败路径。

### 验收标准
- 新增 4 个模块均有测试。
- `go test ./...` 通过。
- `go vet ./...` 通过。
- `golangci-lint run` 通过。
- README 新增模块说明与最小示例。
