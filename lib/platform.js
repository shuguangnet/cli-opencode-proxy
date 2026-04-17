const { PACKAGE_VERSION, REPO_OWNER, REPO_NAME } = require('./constants')

function getTarget() {
  const key = `${process.platform}/${process.arch}`
  const targets = {
    'darwin/arm64': 'darwin-arm64',
    'darwin/x64': 'darwin-amd64',
    'linux/x64': 'linux-amd64',
    'win32/x64': 'windows-amd64',
  }
  const target = targets[key]
  if (!target) {
    throw new Error(`Unsupported platform: ${key}`)
  }
  return target
}

function getAssetName() {
  const target = getTarget()
  return process.platform === 'win32'
    ? `opencode-cli-proxy-${target}.exe`
    : `opencode-cli-proxy-${target}`
}

function getDownloadUrl() {
  if (process.env.OPENCODE_PROXY_BINARY_URL) {
    return process.env.OPENCODE_PROXY_BINARY_URL
  }
  return `https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/v${PACKAGE_VERSION}/${getAssetName()}`
}

module.exports = {
  getTarget,
  getAssetName,
  getDownloadUrl,
}
