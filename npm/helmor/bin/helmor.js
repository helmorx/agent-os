#!/usr/bin/env node
'use strict';

const fs = require('node:fs');
const path = require('node:path');
const { spawn, spawnSync } = require('node:child_process');

const packageRoot = path.resolve(__dirname, '..');
const executable = process.platform === 'win32' ? 'helmor.exe' : 'helmor';
const binaryPath = path.join(packageRoot, 'vendor', executable);

function installIfMissing() {
  if (fs.existsSync(binaryPath)) {
    return;
  }

  const installer = path.join(packageRoot, 'scripts', 'install.js');
  const result = spawnSync(process.execPath, [installer], {
    stdio: 'inherit',
    env: process.env
  });

  if (result.error) {
    console.error(`helmor: failed to run installer: ${result.error.message}`);
    process.exit(1);
  }

  if (result.status !== 0) {
    process.exit(result.status || 1);
  }
}

installIfMissing();

const child = spawn(binaryPath, process.argv.slice(2), {
  stdio: 'inherit'
});

child.on('error', (err) => {
  console.error(`helmor: failed to run native binary: ${err.message}`);
  process.exit(1);
});

child.on('exit', (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }

  process.exit(code == null ? 1 : code);
});
