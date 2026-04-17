const fs = require('fs')
const { execFileSync } = require('child_process')
const { ensureBinaryDownloaded } = require('./downloader')
const { renderConfig } = require('./config-template')
const { configPath, ensureAppDirs, baseDir, binaryPath } = require('./paths')
const { startGateway, stopGateway, getStatus } = require('./process-manager')
const { installService, uninstallService, serviceHelp } = require('./service')
const { CLI_NAME } = require('./constants')

function printHelp() {
  console.log(`${CLI_NAME} commands:
  init               Create default config
  setup              Download binary and create config if missing
  download           Download or refresh the platform binary
  start              Start the gateway in background
  stop               Stop the gateway
  restart            Restart the gateway
  status             Show gateway status
  install-service    Enable auto start
  uninstall-service  Remove auto start
  paths              Show important local paths
  open-config        Print config path
  help               Show this help
`)
}

function parseArgs(argv) {
  const [command = 'help', ...rest] = argv
  const options = {}
  for (let i = 0; i < rest.length; i += 1) {
    const arg = rest[i]
    if (arg.startsWith('--')) {
      const key = arg.slice(2)
      const value = rest[i + 1] && !rest[i + 1].startsWith('--') ? rest[++i] : 'true'
      options[key] = value
    }
  }
  return { command, options }
}

function ensureOpencodeExists(binary) {
  try {
    execFileSync(binary, ['models'], { stdio: 'ignore' })
  } catch {
    throw new Error(`Unable to run '${binary} models'. Install opencode first or pass --binary with an absolute path.`)
  }
}

function writeConfig(options = {}) {
  ensureAppDirs()
  const binary = options.binary || 'opencode'
  ensureOpencodeExists(binary)
  const content = renderConfig({
    host: options.host || '127.0.0.1',
    port: Number(options.port || 18080),
    binary,
    attach: options.attach || '',
    apiKey: options['api-key'] || 'sk-gw-demo',
    model: options.model || 'opencode-go/glm-5.1',
  })
  fs.writeFileSync(configPath, content)
  return configPath
}

async function cmdInit(options) {
  const file = writeConfig(options)
  console.log(`Config created: ${file}`)
}

async function cmdSetup(options) {
  const downloaded = await ensureBinaryDownloaded()
  const existed = fs.existsSync(configPath)
  if (!existed) {
    writeConfig(options)
    console.log(`Config created: ${configPath}`)
  }
  console.log(`Binary ready: ${downloaded}`)
  console.log(`Config ready: ${configPath}`)
}

async function cmdDownload() {
  const file = await ensureBinaryDownloaded({ force: true })
  console.log(`Binary downloaded: ${file}`)
}

async function cmdStart() {
  await ensureBinaryDownloaded()
  const result = startGateway()
  if (result.alreadyRunning) {
    console.log(`Gateway already running (pid ${result.pid})`)
    return
  }
  console.log(`Gateway started (pid ${result.pid})`)
}

async function cmdStop() {
  const stopped = stopGateway()
  console.log(stopped ? 'Gateway stopped' : 'Gateway is not running')
}

async function cmdRestart() {
  stopGateway()
  await cmdStart()
}

async function cmdStatus() {
  const status = getStatus()
  console.log(JSON.stringify(status, null, 2))
}

async function cmdInstallService() {
  await ensureBinaryDownloaded()
  const location = installService()
  console.log(`Auto start installed: ${location}`)
}

async function cmdUninstallService() {
  const location = uninstallService()
  console.log(`Auto start removed: ${location}`)
}

async function cmdPaths() {
  console.log(JSON.stringify({
    baseDir,
    configPath,
    binaryPath: binaryPath(),
  }, null, 2))
}

async function main() {
  const { command, options } = parseArgs(process.argv.slice(2))
  switch (command) {
    case 'init':
      return cmdInit(options)
    case 'setup':
      return cmdSetup(options)
    case 'download':
      return cmdDownload()
    case 'start':
      return cmdStart()
    case 'stop':
      return cmdStop()
    case 'restart':
      return cmdRestart()
    case 'status':
      return cmdStatus()
    case 'install-service':
      return cmdInstallService()
    case 'uninstall-service':
      return cmdUninstallService()
    case 'paths':
      return cmdPaths()
    case 'open-config':
      console.log(configPath)
      return
    case 'help':
    case '--help':
    case '-h':
      printHelp()
      console.log(serviceHelp())
      return
    default:
      throw new Error(`Unknown command: ${command}`)
  }
}

module.exports = { main }
