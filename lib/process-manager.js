const fs = require('fs')
const { spawn } = require('child_process')
const { configPath, pidPath, stdoutPath, stderrPath, ensureAppDirs, binaryPath } = require('./paths')

function readPid() {
  if (!fs.existsSync(pidPath)) return null
  const value = Number(fs.readFileSync(pidPath, 'utf8').trim())
  return Number.isFinite(value) && value > 0 ? value : null
}

function isRunning(pid) {
  if (!pid) return false
  try {
    process.kill(pid, 0)
    return true
  } catch {
    return false
  }
}

function ensureConfigExists() {
  if (!fs.existsSync(configPath)) {
    throw new Error(`Missing config: ${configPath}. Run 'opencode-proxy init' first.`)
  }
}

function startGateway() {
  ensureAppDirs()
  ensureConfigExists()
  const pid = readPid()
  if (isRunning(pid)) {
    return { alreadyRunning: true, pid }
  }
  const out = fs.openSync(stdoutPath, 'a')
  const err = fs.openSync(stderrPath, 'a')
  const child = spawn(binaryPath(), ['-config', configPath], {
    detached: true,
    stdio: ['ignore', out, err],
  })
  child.unref()
  fs.writeFileSync(pidPath, String(child.pid))
  return { alreadyRunning: false, pid: child.pid }
}

function stopGateway() {
  const pid = readPid()
  if (!isRunning(pid)) {
    if (fs.existsSync(pidPath)) fs.rmSync(pidPath, { force: true })
    return false
  }
  process.kill(pid, 'SIGTERM')
  fs.rmSync(pidPath, { force: true })
  return true
}

function getStatus() {
  const pid = readPid()
  return {
    pid,
    running: isRunning(pid),
    configPath,
    binaryPath: binaryPath(),
    stdoutPath,
    stderrPath,
  }
}

module.exports = {
  readPid,
  isRunning,
  startGateway,
  stopGateway,
  getStatus,
}
