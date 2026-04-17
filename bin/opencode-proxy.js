#!/usr/bin/env node

require('../lib/cli').main().catch((error) => {
  console.error(error.message || error)
  process.exit(1)
})
