const fs = require('fs')
const os = require('os')
const path = require('path')
const { execFileSync } = require('child_process')
const { CLI_NAME, SERVICE_NAME } = require('./constants')
const { configPath, stdoutPath, stderrPath, ensureAppDirs, binaryPath, baseDir } = require('./paths')

function installMacService() {
  ensureAppDirs()
  const agentDir = path.join(os.homedir(), 'Library', 'LaunchAgents')
  fs.mkdirSync(agentDir, { recursive: true })
  const plistPath = path.join(agentDir, `${SERVICE_NAME}.plist`)
  const plist = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>${SERVICE_NAME}</string>
    <key>ProgramArguments</key>
    <array>
      <string>${binaryPath()}</string>
      <string>-config</string>
      <string>${configPath}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>${stdoutPath}</string>
    <key>StandardErrorPath</key>
    <string>${stderrPath}</string>
    <key>WorkingDirectory</key>
    <string>${baseDir}</string>
  </dict>
</plist>
`
  fs.writeFileSync(plistPath, plist)
  execFileSync('launchctl', ['unload', plistPath], { stdio: 'ignore' })
  execFileSync('launchctl', ['load', plistPath], { stdio: 'inherit' })
  return plistPath
}

function uninstallMacService() {
  const plistPath = path.join(os.homedir(), 'Library', 'LaunchAgents', `${SERVICE_NAME}.plist`)
  if (fs.existsSync(plistPath)) {
    try {
      execFileSync('launchctl', ['unload', plistPath], { stdio: 'ignore' })
    } catch {}
    fs.rmSync(plistPath, { force: true })
  }
  return plistPath
}

function installService() {
  if (process.platform === 'darwin') return installMacService()
  throw new Error(`install-service is not implemented for ${process.platform} yet`)
}

function uninstallService() {
  if (process.platform === 'darwin') return uninstallMacService()
  throw new Error(`uninstall-service is not implemented for ${process.platform} yet`)
}

function serviceHelp() {
  return `${CLI_NAME} install-service currently supports macOS LaunchAgent only.`
}

module.exports = {
  installService,
  uninstallService,
  serviceHelp,
}
