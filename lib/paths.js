const fs = require('fs')
const os = require('os')
const path = require('path')

const baseDir = path.join(os.homedir(), '.opencode-cli-proxy')
const binDir = path.join(baseDir, 'bin')
const logDir = path.join(baseDir, 'logs')
const runDir = path.join(baseDir, 'run')
const configPath = path.join(baseDir, 'config.yaml')
const pidPath = path.join(runDir, 'gateway.pid')
const stdoutPath = path.join(logDir, 'gateway.stdout.log')
const stderrPath = path.join(logDir, 'gateway.stderr.log')

function ensureAppDirs() {
  for (const dir of [baseDir, binDir, logDir, runDir]) {
    fs.mkdirSync(dir, { recursive: true })
  }
}

function binaryPath() {
  return path.join(binDir, process.platform === 'win32' ? 'opencode-cli-proxy.exe' : 'opencode-cli-proxy')
}

module.exports = {
  baseDir,
  binDir,
  logDir,
  runDir,
  configPath,
  pidPath,
  stdoutPath,
  stderrPath,
  ensureAppDirs,
  binaryPath,
}
