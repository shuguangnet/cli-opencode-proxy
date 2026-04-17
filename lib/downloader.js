const fs = require('fs')
const https = require('https')
const { pipeline } = require('stream')
const { promisify } = require('util')
const { ensureAppDirs, binaryPath } = require('./paths')
const { getDownloadUrl } = require('./platform')

const pipe = promisify(pipeline)

function request(url) {
  return new Promise((resolve, reject) => {
    https.get(url, (response) => {
      if (response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
        resolve(request(response.headers.location))
        return
      }
      if (response.statusCode !== 200) {
        reject(new Error(`Download failed with status ${response.statusCode}: ${url}`))
        return
      }
      resolve(response)
    }).on('error', reject)
  })
}

async function ensureBinaryDownloaded({ force = false } = {}) {
  ensureAppDirs()
  const targetPath = binaryPath()
  if (!force && fs.existsSync(targetPath)) {
    return targetPath
  }
  const response = await request(getDownloadUrl())
  const tmpPath = `${targetPath}.tmp`
  await pipe(response, fs.createWriteStream(tmpPath, { mode: 0o755 }))
  fs.renameSync(tmpPath, targetPath)
  fs.chmodSync(targetPath, 0o755)
  return targetPath
}

module.exports = { ensureBinaryDownloaded }
