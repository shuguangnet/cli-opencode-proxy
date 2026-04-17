# OpenCode CLI Proxy

<p align="center">
  Turn your local <code>opencode</code> CLI into an OpenAI-compatible API.
</p>

<p align="center">
  <a href="#english">English</a> ·
  <a href="#简体中文">简体中文</a>
</p>

---

## English

### Overview

OpenCode CLI Proxy exposes your local `opencode` runtime through OpenAI-compatible endpoints, so tools like **Cherry Studio**, **Cursor**, **NextChat**, and OpenAI-compatible SDKs can use it without custom integrations.

Instead of assuming a remote OpenCode HTTP service, this project talks to the local `opencode` CLI directly:

- `opencode models` for model discovery
- `opencode run --format json --model ...` for completions
- optional `--attach` support for an existing `opencode serve` instance

### Features

- OpenAI-compatible endpoints
  - `GET /v1/models`
  - `POST /v1/chat/completions`
  - `POST /v1/completions`
- Streaming response support via SSE
- Gateway API key authentication
- Native desktop GUI built with Fyne
- Works with local `opencode` models
- Model alias mapping via config

### Architecture

```text
OpenAI-compatible client
        |
        v
+---------------------------+
| OpenCode CLI Proxy        |
| - auth                    |
| - request mapping         |
| - response mapping        |
| - SSE conversion          |
+---------------------------+
        |
        v
+---------------------------+
| local opencode CLI        |
| - opencode models         |
| - opencode run            |
| - optional attach server  |
+---------------------------+
```

### Requirements

Before using this project, make sure:

1. `opencode` is installed
2. `opencode models` works
3. `opencode run --model <model> "hello"` works
4. Go is installed

Example:

```bash
opencode models
opencode run --model opencode-go/glm-5.1 "reply with exactly: ok"
```

If `opencode` is not in your PATH, use an absolute path in config or in the desktop app.

### Quick Start

#### Run the server

```bash
make run
```

Default address:

```text
http://127.0.0.1:18080
```

#### Check service status

```bash
curl http://127.0.0.1:18080/
curl http://127.0.0.1:18080/health
curl http://127.0.0.1:18080/v1
```

#### List models

```bash
curl http://127.0.0.1:18080/v1/models \
  -H "Authorization: Bearer sk-gw-demo"
```

#### Chat completion

```bash
curl http://127.0.0.1:18080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-gw-demo" \
  -d '{
    "model": "opencode-go/glm-5.1",
    "messages": [
      { "role": "user", "content": "Introduce yourself briefly." }
    ]
  }'
```

#### Streaming chat completion

```bash
curl -N http://127.0.0.1:18080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-gw-demo" \
  -d '{
    "model": "opencode-go/glm-5.1",
    "stream": true,
    "messages": [
      { "role": "user", "content": "Explain Go in three sentences." }
    ]
  }'
```

### Use in Other Apps

Use the proxy as an OpenAI-compatible provider:

- **Base URL**: `http://127.0.0.1:18080/v1`
- **API Key**: `sk-gw-demo`
- **Model**: `opencode-go/glm-5.1`

#### Cherry Studio

- Provider Type: `OpenAI Compatible`
- Base URL: `http://127.0.0.1:18080/v1`
- API Key: `sk-gw-demo`
- Model: `opencode-go/glm-5.1`

#### NextChat

- Custom Endpoint: `http://127.0.0.1:18080/v1`
- API Key: `sk-gw-demo`
- Model: `opencode-go/glm-5.1`

#### Cursor / OpenAI-compatible clients

- OpenAI Base URL: `http://127.0.0.1:18080/v1`
- OpenAI API Key: `sk-gw-demo`
- Model Name: `opencode-go/glm-5.1`

### Desktop App

This project includes a native desktop app for macOS and Windows.

Run it with:

```bash
make desktop
```

Or:

```bash
go run -buildvcs=false ./cmd/desktop
```

Build it with:

```bash
make build-desktop
```

Desktop fields:

- Config file
- Listen host
- Listen port
- Opencode binary
- Attach server
- Gateway API key
- Default model
- Allowed models
- Client Base URL

### Configuration

Example config: `configs/config.example.yaml`

```yaml
server:
  host: 0.0.0.0
  port: 18080
  read_timeout: 15s
  write_timeout: 0s

upstream:
  binary: opencode
  attach: ""
  timeout: 120s

models:
  opencode-go/glm-5.1: opencode-go/glm-5.1
  opencode-go/glm-5: opencode-go/glm-5

accounts:
  default:
    auth_mode: local
    token: ""

keys:
  sk-gw-demo:
    account: default
    allowed_models:
      - opencode-go/glm-5.1
      - opencode-go/glm-5

mapping:
  temperature:
    target_min: 0
    target_max: 1

rate_limit:
  enabled: false
  rpm: 60
```

#### Config notes

- `server`: local bind address and timeouts
- `upstream.binary`: executable name or absolute path to `opencode`
- `upstream.attach`: optional `opencode serve` endpoint
- `models`: alias-to-real-model mapping
- `keys`: gateway API key to allowed models mapping

Example alias mapping:

```yaml
models:
  glm-latest: opencode-go/glm-5.1
```

Then clients can use `glm-latest` while the proxy runs `opencode-go/glm-5.1`.

### Supported Routes

Public routes:

- `GET /`
- `GET /health`
- `GET /v1`

Authenticated routes:

- `GET /v1/models`
- `POST /v1/chat/completions`
- `POST /v1/completions`

### How It Works

For `POST /v1/chat/completions`:

1. The client sends an OpenAI-style request
2. The proxy validates the gateway API key
3. The request is mapped into an internal chat request
4. The proxy calls local `opencode`
5. Output is parsed and converted back into OpenAI-style JSON or SSE

### Limitations

Current MVP limitations:

- Request mapping is still text-oriented for `opencode run`
- Message/role handling is intentionally simple for now
- `/v1/models` directly reflects local `opencode models`
- `accounts` is currently kept mainly for config structure compatibility
- No full audit, metrics, or production rate limiting yet
- No Docker release flow yet

### Development

```bash
make run
make desktop
make build
make build-desktop
make test
```

Core files:

- `cmd/server/main.go`
- `cmd/desktop/main.go`
- `internal/app/gateway.go`
- `internal/openai/handlers.go`
- `internal/upstream/client.go`
- `internal/adapter/chat_mapper.go`
- `internal/server/router.go`

### Roadmap

- Better message and role mapping
- Richer streaming event conversion
- Better error mapping
- Better model alias and default-model management
- Request logging and audit support
- Token usage reporting
- Docker packaging
- More OpenAI-compatible endpoints

### License

No `LICENSE` file is included yet. Add one before public open-source release.

---

## 简体中文

### 项目简介

OpenCode CLI Proxy 可以把本地 `opencode` CLI 包装成 OpenAI 兼容接口，让 **Cherry Studio**、**Cursor**、**NextChat**、OpenAI SDK 等客户端无需定制适配即可直接接入。

当前实现不是去对接一个假设中的远程 OpenCode HTTP 服务，而是直接调用本机 `opencode`：

- 用 `opencode models` 获取模型列表
- 用 `opencode run --format json --model ...` 发起对话
- 可选通过 `--attach` 连接已有的 `opencode serve`

### 功能特性

- OpenAI 兼容接口
  - `GET /v1/models`
  - `POST /v1/chat/completions`
  - `POST /v1/completions`
- 支持流式 SSE 输出
- 支持网关 API Key 鉴权
- 提供原生桌面 GUI
- 支持模型别名映射
- 基于本地 `opencode` 能力运行

### 快速开始

#### 启动服务

```bash
make run
```

默认地址：

```text
http://127.0.0.1:18080
```

#### 测试接口

```bash
curl http://127.0.0.1:18080/
curl http://127.0.0.1:18080/health
curl http://127.0.0.1:18080/v1
```

查看模型：

```bash
curl http://127.0.0.1:18080/v1/models \
  -H "Authorization: Bearer sk-gw-demo"
```

聊天测试：

```bash
curl http://127.0.0.1:18080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-gw-demo" \
  -d '{
    "model": "opencode-go/glm-5.1",
    "messages": [
      { "role": "user", "content": "你好，简单介绍一下你自己" }
    ]
  }'
```

流式测试：

```bash
curl -N http://127.0.0.1:18080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-gw-demo" \
  -d '{
    "model": "opencode-go/glm-5.1",
    "stream": true,
    "messages": [
      { "role": "user", "content": "用三句话介绍 Go 语言" }
    ]
  }'
```

### 在其他软件中使用

按 OpenAI Compatible 方式填写：

- Base URL：`http://127.0.0.1:18080/v1`
- API Key：`sk-gw-demo`
- Model：`opencode-go/glm-5.1`

#### Cherry Studio

- Provider Type：`OpenAI Compatible`
- Base URL：`http://127.0.0.1:18080/v1`
- API Key：`sk-gw-demo`
- Model：`opencode-go/glm-5.1`

#### NextChat

- Custom Endpoint：`http://127.0.0.1:18080/v1`
- API Key：`sk-gw-demo`
- Model：`opencode-go/glm-5.1`

#### Cursor / 通用 OpenAI 客户端

- OpenAI Base URL：`http://127.0.0.1:18080/v1`
- OpenAI API Key：`sk-gw-demo`
- Model Name：`opencode-go/glm-5.1`

### 桌面版

项目提供原生桌面 GUI，适用于 macOS / Windows。

启动：

```bash
make desktop
```

或：

```bash
go run -buildvcs=false ./cmd/desktop
```

构建：

```bash
make build-desktop
```

可配置项：

- Config file
- Listen host
- Listen port
- Opencode binary
- Attach server
- Gateway API key
- Default model
- Allowed models
- Client Base URL

### 配置示例

```yaml
server:
  host: 0.0.0.0
  port: 18080
  read_timeout: 15s
  write_timeout: 0s

upstream:
  binary: opencode
  attach: ""
  timeout: 120s

models:
  opencode-go/glm-5.1: opencode-go/glm-5.1
  opencode-go/glm-5: opencode-go/glm-5

accounts:
  default:
    auth_mode: local
    token: ""

keys:
  sk-gw-demo:
    account: default
    allowed_models:
      - opencode-go/glm-5.1
      - opencode-go/glm-5
```

### 当前支持的接口

无需鉴权：

- `GET /`
- `GET /health`
- `GET /v1`

需要鉴权：

- `GET /v1/models`
- `POST /v1/chat/completions`
- `POST /v1/completions`

### 已知限制

当前版本属于 MVP：

- 目前仍以 `opencode run` 的文本事件输出为核心
- message / role 映射还比较基础
- `/v1/models` 直接读取本地 `opencode models`
- 还没有完整的审计、指标、生产级限流能力
- 暂未提供 Docker 发布流程

### 后续规划

- 更准确的消息映射
- 更丰富的流式事件适配
- 更好的错误映射
- 模型别名和默认模型管理
- 请求日志与审计
- Token 使用统计
- Docker 化发布
- 更多 OpenAI 兼容接口

### License

当前仓库还没有 `LICENSE` 文件。公开发布前建议补充。
