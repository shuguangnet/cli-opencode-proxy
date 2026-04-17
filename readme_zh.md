# OpenCode CLI Proxy

<p align="center">
  将本地 <code>opencode</code> CLI 转换为 OpenAI 兼容 API。
</p>

<p align="center">
  <a href="./README.md">English</a> ·
  <a href="./readme_zh.md">简体中文</a>
</p>

---

## 项目简介

OpenCode CLI Proxy 可以把本地 `opencode` CLI 包装成 OpenAI 兼容接口，让 **Cherry Studio**、**Cursor**、**NextChat**、OpenAI SDK 等客户端无需定制适配即可直接接入。

当前实现不是去对接一个假设中的远程 OpenCode HTTP 服务，而是直接调用本机 `opencode`：

- 用 `opencode models` 获取模型列表
- 用 `opencode run --format json --model ...` 发起对话
- 可选通过 `--attach` 连接已有的 `opencode serve`

## 功能特性

- OpenAI 兼容接口
  - `GET /v1/models`
  - `POST /v1/chat/completions`
  - `POST /v1/completions`
- 支持流式 SSE 输出
- 支持网关 API Key 鉴权
- 提供原生桌面 GUI
- 支持模型别名映射
- 基于本地 `opencode` 能力运行

## 环境要求

使用前请确认：

1. 已安装 `opencode`
2. `opencode models` 可正常执行
3. `opencode run --model <model> "hello"` 可正常执行
4. 已安装 Go

示例：

```bash
opencode models
opencode run --model opencode-go/glm-5.1 "reply with exactly: ok"
```

如果 `opencode` 不在 PATH 中，可以在配置文件或桌面版中填写其绝对路径。

## 安装

### 克隆仓库

```bash
git clone <your-repo-url>
cd opencode-cli-proxy
```

### 通过 npm 安装

```bash
npm install -g opencode-cli-proxy
opencode-proxy setup
opencode-proxy start
```

### 直接源码运行

```bash
make run
```

### 构建二进制

```bash
make build
make build-desktop
```

### 后台服务与自启动

```bash
opencode-proxy install-service
opencode-proxy status
```

当前支持：

- macOS LaunchAgent
- Linux / Windows 服务化：后续补充

## 快速开始

### 启动服务

```bash
make run
```

默认地址：

```text
http://127.0.0.1:18080
```

### 测试接口

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

## 在其他软件中使用

按 OpenAI Compatible 方式填写：

- **Base URL**：`http://127.0.0.1:18080/v1`
- **API Key**：`sk-gw-demo`
- **Model**：`opencode-go/glm-5.1`

### Cherry Studio

- Provider Type：`OpenAI Compatible`
- Base URL：`http://127.0.0.1:18080/v1`
- API Key：`sk-gw-demo`
- Model：`opencode-go/glm-5.1`

### NextChat

- Custom Endpoint：`http://127.0.0.1:18080/v1`
- API Key：`sk-gw-demo`
- Model：`opencode-go/glm-5.1`

### Cursor / 通用 OpenAI 客户端

- OpenAI Base URL：`http://127.0.0.1:18080/v1`
- OpenAI API Key：`sk-gw-demo`
- Model Name：`opencode-go/glm-5.1`

## 桌面版

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

## 配置示例

示例配置：`configs/config.example.yaml`

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

### 配置说明

- `server`：本地监听地址和超时配置
- `upstream.binary`：`opencode` 可执行文件名或绝对路径
- `upstream.attach`：可选的 `opencode serve` 地址
- `models`：对外模型名到真实模型名的映射
- `keys`：网关 API Key 与允许模型的映射关系

模型别名示例：

```yaml
models:
  glm-latest: opencode-go/glm-5.1
```

这样客户端就可以使用 `glm-latest`，而代理内部实际调用 `opencode-go/glm-5.1`。

## 当前支持的接口

无需鉴权：

- `GET /`
- `GET /health`
- `GET /v1`

需要鉴权：

- `GET /v1/models`
- `POST /v1/chat/completions`
- `POST /v1/completions`

## 工作原理

以 `POST /v1/chat/completions` 为例：

1. 客户端发送 OpenAI 风格请求
2. 网关校验 API Key
3. 将请求转换为内部 chat 请求结构
4. 调用本地 `opencode`
5. 将输出转换回 OpenAI JSON 或 SSE

## 项目结构

```text
cmd/
  server/          HTTP 服务入口
  desktop/         Fyne 桌面应用
configs/
  config.example.yaml
internal/
  adapter/         OpenAI 请求与响应映射
  app/             Gateway 生命周期管理
  config/          配置加载与校验
  domain/          共享协议类型
  openai/          HTTP 处理器
  server/          路由与中间件
  upstream/        本地 opencode CLI 集成
```

## FAQ

### 这个项目会调用远程 OpenCode HTTP API 吗？

不会。当前实现主要基于本地 `opencode` CLI。

### 可以接入 Cursor、Cherry Studio、NextChat 吗？

可以。把这些客户端指向代理提供的 OpenAI 兼容 `/v1` 地址即可。

### 支持流式输出吗？

支持。流式响应会以 OpenAI 兼容的 SSE 格式返回。

### 客户端需要知道真实上游凭证吗？

不需要。客户端只使用本代理配置的 gateway API key。

## 已知限制

当前版本属于 MVP：

- 目前仍以 `opencode run` 的文本事件输出为核心
- message / role 映射还比较基础
- `/v1/models` 直接读取本地 `opencode models`
- `accounts` 目前主要为配置结构兼容而保留
- 还没有完整的审计、指标、生产级限流能力
- 暂未提供 Docker 发布流程

## 开发命令

```bash
make run
make desktop
make build
make build-desktop
make test
```

核心文件：

- `cmd/server/main.go`
- `cmd/desktop/main.go`
- `internal/app/gateway.go`
- `internal/openai/handlers.go`
- `internal/upstream/client.go`
- `internal/adapter/chat_mapper.go`
- `internal/server/router.go`

## 后续规划

- [ ] 更准确的消息映射
- [ ] 更丰富的流式事件适配
- [ ] 更好的错误映射
- [ ] 模型别名和默认模型管理
- [ ] 请求日志与审计
- [ ] Token 使用统计
- [ ] Docker 化发布
- [ ] 更多 OpenAI 兼容接口

## 发布

为 npm 分发构建 release 二进制：

```bash
make release-dist
```

执行本地发布前检查：

```bash
npm run preflight
```

本地打 npm 包：

```bash
make npm-pack
```

构建产物会输出到 `dist/`，文件名如下：

- `opencode-cli-proxy-darwin-arm64`
- `opencode-cli-proxy-darwin-amd64`
- `opencode-cli-proxy-linux-amd64`
- `opencode-cli-proxy-windows-amd64.exe`
- `checksums-v<version>.txt`

### GitHub Actions 自动发布

仓库已包含 `.github/workflows/release.yml`。

使用方式：

1. 在仓库 Secrets 中添加 `NPM_TOKEN`
2. 提交代码
3. 创建与 `package.json` 一致的 tag，例如 `v0.1.0`
4. 推送该 tag

工作流会自动：

- 执行测试
- 构建 `dist/*`
- 创建 GitHub Release
- 上传 release 产物
- 发布 npm 包

如果手动发布，也请把这些文件上传到与 `package.json` 版本一致的 GitHub Release，例如 `v0.1.0`。

## License

当前仓库还没有 `LICENSE` 文件。公开发布前建议补充。
