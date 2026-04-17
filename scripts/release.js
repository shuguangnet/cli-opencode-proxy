#!/usr/bin/env node

const fs = require('fs')
const path = require('path')
const { spawnSync } = require('child_process')
const packageJson = require('../package.json')

const rootDir = path.resolve(__dirname, '..')
const distDir = path.join(rootDir, 'dist')
const version = packageJson.version

const targets = [
  { goos: 'darwin', goarch: 'arm64', output: 'opencode-cli-proxy-darwin-arm64' },
  { goos: 'darwin', goarch: 'amd64', output: 'opencode-cli-proxy-darwin-amd64' },
  { goos: 'linux', goarch: 'amd64', output: 'opencode-cli-proxy-linux-amd64' },
  { goos: 'windows', goarch: 'amd64', output: 'opencode-cli-proxy-windows-amd64.exe' },
]

function run(command, args, options = {}) {
  const result = spawnSync(command, args, {
    cwd: rootDir,
    stdio: 'inherit',
    env: { ...process.env, ...(options.env || {}) },
  })
  if (result.status !== 0) {
    process.exit(result.status || 1)
  }
}

function ensureDir(dir) {
  fs.mkdirSync(dir, { recursive: true })
}

function cleanDist() {
  fs.rmSync(distDir, { recursive: true, force: true })
  ensureDir(distDir)
}

function buildTarget(target) {
  const outputPath = path.join(distDir, target.output)
  run('go', ['build', '-buildvcs=false', '-o', outputPath, './cmd/server'], {
    env: {
      GOOS: target.goos,
      GOARCH: target.goarch,
      CGO_ENABLED: '0',
    },
  })
}

function writeChecksums() {
  const files = fs.readdirSync(distDir).filter((name) => !name.endsWith('.sha256'))
  const lines = []
  for (const file of files) {
    const fullPath = path.join(distDir, file)
    const hash = spawnSync('shasum', ['-a', '256', fullPath], { encoding: 'utf8' })
    if (hash.status !== 0) {
      process.exit(hash.status || 1)
    }
    lines.push(hash.stdout.trim())
  }
  fs.writeFileSync(path.join(distDir, `checksums-v${version}.txt`), `${lines.join('\n')}\n`)
}

function main() {
  cleanDist()
  for (const target of targets) {
    buildTarget(target)
  }
  writeChecksums()
  console.log(`Release artifacts written to ${distDir}`)
}

main()
