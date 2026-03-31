# Common Next Modules Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add four high-frequency reusable modules (`config`, `xerrors`, `cache`, `httpx`) to `github.com/xpfo-go/common/v2` with tests and docs.

**Architecture:** Keep each module independent and lightweight. Build around explicit APIs, deterministic unit tests, and minimal dependencies. Use test-first cycles for each module, then integrate docs updates and full repository verification.

**Tech Stack:** Go 1.22+, stdlib, `gopkg.in/yaml.v3`, `golang.org/x/sync/singleflight`, `go test`, `go vet`, `golangci-lint`.

---

### Task 1: Add `xerrors` module

**Files:**
- Create: `xerrors/error.go`
- Create: `xerrors/error_test.go`

- [ ] **Step 1: Write failing tests for error wrapping, code lookup, and HTTP status mapping**
- [ ] **Step 2: Run `go test ./xerrors` and confirm RED state**
- [ ] **Step 3: Implement minimal `xerrors` API**
- [ ] **Step 4: Run `go test ./xerrors` and confirm GREEN state**

### Task 2: Add `config` module

**Files:**
- Create: `config/file.go`
- Create: `config/env.go`
- Create: `config/config_test.go`

- [ ] **Step 1: Write failing tests for file load, default tags, env override, required validation**
- [ ] **Step 2: Run `go test ./config` and confirm RED state**
- [ ] **Step 3: Implement `LoadFile/LoadJSON/LoadYAML/ApplyDefaults/LoadFromEnv/ValidateRequired`**
- [ ] **Step 4: Run `go test ./config` and confirm GREEN state**

### Task 3: Add `cache` module

**Files:**
- Create: `cache/cache.go`
- Create: `cache/cache_test.go`

- [ ] **Step 1: Write failing tests for TTL, LRU eviction, and `GetOrLoad` dedup behavior**
- [ ] **Step 2: Run `go test ./cache` and confirm RED state**
- [ ] **Step 3: Implement thread-safe cache with cleanup loop and singleflight loader**
- [ ] **Step 4: Run `go test ./cache` and confirm GREEN state**

### Task 4: Add `httpx` module

**Files:**
- Create: `httpx/client.go`
- Create: `httpx/client_test.go`

- [ ] **Step 1: Write failing tests for retry policy and JSON request/response handling**
- [ ] **Step 2: Run `go test ./httpx` and confirm RED state**
- [ ] **Step 3: Implement minimal retryable client and JSON helper**
- [ ] **Step 4: Run `go test ./httpx` and confirm GREEN state**

### Task 5: Documentation and integration updates

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Add the 4 new modules to module table**
- [ ] **Step 2: Add concise usage snippets and notes on module intent**

### Task 6: Full verification

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`

- [ ] **Step 1: Run `go test ./...`**
- [ ] **Step 2: Run `go vet ./...`**
- [ ] **Step 3: Run `golangci-lint run`**
- [ ] **Step 4: Summarize final API additions and migration impact**
