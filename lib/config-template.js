function renderConfig({ host = '127.0.0.1', port = 18080, binary = 'opencode', attach = '', apiKey = 'sk-gw-demo', model = 'opencode-go/glm-5.1' }) {
  return `server:
  host: ${host}
  port: ${port}
  read_timeout: 15s
  write_timeout: 0s

upstream:
  binary: ${binary}
  attach: ${attach ? JSON.stringify(attach) : '""'}
  timeout: 120s

models:
  ${model}: ${model}

accounts:
  default:
    auth_mode: local
    token: ""

keys:
  ${apiKey}:
    account: default
    allowed_models:
      - ${model}

mapping:
  temperature:
    target_min: 0
    target_max: 1

rate_limit:
  enabled: false
  rpm: 60
`
}

module.exports = { renderConfig }
