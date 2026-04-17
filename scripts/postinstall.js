const { ensureBinaryDownloaded } = require('../lib/downloader')

ensureBinaryDownloaded().catch((error) => {
  console.warn(`[opencode-cli-proxy] postinstall skipped: ${error.message}`)
})
