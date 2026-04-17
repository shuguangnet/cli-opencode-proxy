# OpenCode CLI Proxy

<p align="center">
  Turn your local <code>opencode</code> CLI into an OpenAI-compatible API.
</p>

<p align="center">
  <a href="./README.md">English</a> ·
  <a href="./readme_zh.md">简体中文</a>
</p>

---

## Overview

OpenCode CLI Proxy exposes your local `opencode` runtime through OpenAI-compatible endpoints, so tools like **Cherry Studio**, **Cursor**, **NextChat**, and OpenAI-compatible SDKs can use it without custom integrations.

Instead of assuming a remote OpenCode HTTP service, this project talks to the local `opencode` CLI directly:

- `opencode models` for model discovery
- `opencode run --format json --model ...` for completions
- optional `--attach` support for an existing `opencode serve` instance

## Features

- OpenAI-compatible endpoints
  - `GET /v1/models`
  - `POST /v1/chat/completions`
  - `POST /v1/completions`
- Streaming response support via SSE
- Gateway API key authentication
- Native desktop GUI built with Fyne
- Works with local `opencode` models
- Model alias mapping via config

## Architecture

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

## Requirements

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

## Installation

### Clone the repository

```bash
git clone <your-repo-url>
cd opencode-cli-proxy
```

### Start from source

```bash
make run
```

### Build binaries

```bash
make build
make build-desktop
```

## Quick Start

### Run the server

```bash
make run
```

Default address:

```text
http://127.0.0.1:18080
```

### Check service status

```bash
curl http://127.0.0.1:18080/
curl http://127.0.0.1:18080/health
curl http://127.0.0.1:18080/v1
```

### List models

```bash
curl http://127.0.0.1:18080/v1/models \
  -H "Authorization: Bearer sk-gw-demo"
```

### Chat completion

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

### Streaming chat completion

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

## Use in Other Apps

Use the proxy as an OpenAI-compatible provider:

- **Base URL**: `http://127.0.0.1:18080/v1`
- **API Key**: `sk-gw-demo`
- **Model**: `opencode-go/glm-5.1`

### Cherry Studio

- Provider Type: `OpenAI Compatible`
- Base URL: `http://127.0.0.1:18080/v1`
- API Key: `sk-gw-demo`
- Model: `opencode-go/glm-5.1`

### NextChat

- Custom Endpoint: `http://127.0.0.1:18080/v1`
- API Key: `sk-gw-demo`
- Model: `opencode-go/glm-5.1`

### Cursor / OpenAI-compatible clients

- OpenAI Base URL: `http://127.0.0.1:18080/v1`
- OpenAI API Key: `sk-gw-demo`
- Model Name: `opencode-go/glm-5.1`

## Desktop App

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

## Configuration

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

### Config notes

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

## Supported Routes

Public routes:

- `GET /`
- `GET /health`
- `GET /v1`

Authenticated routes:

- `GET /v1/models`
- `POST /v1/chat/completions`
- `POST /v1/completions`

## How It Works

For `POST /v1/chat/completions`:

1. The client sends an OpenAI-style request
2. The proxy validates the gateway API key
3. The request is mapped into an internal chat request
4. The proxy calls local `opencode`
5. Output is parsed and converted back into OpenAI-style JSON or SSE

## Project Structure

```text
cmd/
  server/          HTTP server entry
  desktop/         Fyne desktop app
configs/
  config.example.yaml
internal/
  adapter/         OpenAI request/response mapping
  app/             Gateway lifecycle
  config/          Config loading and validation
  domain/          Shared protocol types
  openai/          HTTP handlers
  server/          Router and middleware
  upstream/        Local opencode CLI integration
```

## FAQ

### Does this proxy call a remote OpenCode HTTP API?

No. The current implementation primarily uses the local `opencode` CLI.

### Can I use it with Cursor, Cherry Studio, or NextChat?

Yes. Point them to the proxy's OpenAI-compatible `/v1` endpoint.

### Does it support streaming?

Yes. Streaming responses are exposed as SSE in an OpenAI-compatible format.

### Do I need to expose my real upstream credentials to clients?

No. Clients only use the gateway API key configured in this proxy.

## Limitations

Current MVP limitations:

- Request mapping is still text-oriented for `opencode run`
- Message/role handling is intentionally simple for now
- `/v1/models` directly reflects local `opencode models`
- `accounts` is currently kept mainly for config structure compatibility
- No full audit, metrics, or production rate limiting yet
- No Docker release flow yet

## Development

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

## Roadmap

- [ ] Better message and role mapping
- [ ] Richer streaming event conversion
- [ ] Better error mapping
- [ ] Better model alias and default-model management
- [ ] Request logging and audit support
- [ ] Token usage reporting
- [ ] Docker packaging
- [ ] More OpenAI-compatible endpoints

## License

No `LICENSE` file is included yet. Add one before public open-source release.
