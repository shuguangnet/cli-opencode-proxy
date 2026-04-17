#!/usr/bin/env node

const fs = require('fs')
const path = require('path')
const { spawnSync } = require('child_process')
const packageJson = require('../package.json')

const rootDir = path.resolve(__dirname, '..')
const distDir = path.join(rootDir, 'dist')
const expectedAssets = [
  'opencode-cli-proxy-darwin-arm64',
  'opencode-cli-proxy-darwin-amd64',
  'opencode-cli-proxy-linux-amd64',
  'opencode-cli-proxy-windows-amd64.exe',
  `checksums-v${packageJson.version}.txt`,
]

function run(command, args) {
  return spawnSync(command, args, {
    cwd: rootDir,
    encoding: 'utf8',
    stdio: 'pipe',
  })
}

function assert(condition, message) {
  if (!condition) {
    throw new Error(message)
  }
}

function checkCleanVersionTag() {
  assert(/^\d+\.\d+\.\d+$/.test(packageJson.version), `Invalid package version: ${packageJson.version}`)
}

function checkDist() {
  assert(fs.existsSync(distDir), 'dist directory is missing. Run make release-dist first.')
  for (const file of expectedAssets) {
    assert(fs.existsSync(path.join(distDir, file)), `Missing release artifact: dist/${file}`)
  }
}

function checkGitStatus() {
  const result = run('git', ['status', '--short'])
  assert(result.status === 0, result.stderr || 'git status failed')
}

function checkGoTests() {
  const result = run('go', ['test', './...'])
  assert(result.status === 0, result.stdout + result.stderr)
}

function checkNpmPack() {
  const result = run('npm', ['pack', '--dry-run'])
  assert(result.status === 0, result.stdout + result.stderr)
}

function main() {
  checkCleanVersionTag()
  checkDist()
  checkGitStatus()
  checkGoTests()
  checkNpmPack()
  console.log(`Preflight passed for v${packageJson.version}`)
}

main()
