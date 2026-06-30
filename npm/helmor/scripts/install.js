#!/usr/bin/env node
'use strict';

const crypto = require('node:crypto');
const fs = require('node:fs');
const https = require('node:https');
const os = require('node:os');
const path = require('node:path');
const { spawnSync } = require('node:child_process');

const pkg = require('../package.json');
const {
  checksumFor,
  platformTarget,
  releaseBaseUrl,
  tagFromVersion
} = require('./lib');

const packageRoot = path.resolve(__dirname, '..');
const vendorDir = path.join(packageRoot, 'vendor');

function log(message) {
  if (!process.env.HELMOR_NPM_QUIET) {
    console.log(message);
  }
}

function downloadFile(url, destination, redirects = 0) {
  if (redirects > 5) {
    return Promise.reject(new Error(`too many redirects for ${url}`));
  }

  return new Promise((resolve, reject) => {
    const request = https.get(url, {
      headers: {
        'user-agent': `helmor-npm/${pkg.version}`
      }
    }, (response) => {
      if (response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
        response.resume();
        const redirectUrl = new URL(response.headers.location, url).toString();
        downloadFile(redirectUrl, destination, redirects + 1).then(resolve, reject);
        return;
      }

      if (response.statusCode !== 200) {
        response.resume();
        reject(new Error(`download failed with HTTP ${response.statusCode}: ${url}`));
        return;
      }

      const file = fs.createWriteStream(destination, { mode: 0o600 });
      response.pipe(file);

      file.on('finish', () => {
        file.close(resolve);
      });

      file.on('error', (err) => {
        fs.rm(destination, { force: true }, () => reject(err));
      });
    });

    request.on('error', reject);
  });
}

function sha256(filePath) {
  const hash = crypto.createHash('sha256');
  hash.update(fs.readFileSync(filePath));
  return hash.digest('hex');
}

function run(command, args) {
  const result = spawnSync(command, args, {
    stdio: 'inherit'
  });

  if (result.error) {
    throw result.error;
  }

  if (result.status !== 0) {
    throw new Error(`${command} exited with status ${result.status}`);
  }
}

function quotePowerShell(value) {
  return `'${value.replace(/'/g, "''")}'`;
}

function extractArchive(archivePath, destination, target) {
  if (target.archiveExt === 'zip') {
    const command = `Expand-Archive -LiteralPath ${quotePowerShell(archivePath)} -DestinationPath ${quotePowerShell(destination)} -Force`;
    run('powershell.exe', ['-NoProfile', '-ExecutionPolicy', 'Bypass', '-Command', command]);
    return;
  }

  run('tar', ['-xzf', archivePath, '-C', destination]);
}

function verifyInstalledVersion(binaryPath, expectedVersion) {
  const result = spawnSync(binaryPath, ['version'], {
    encoding: 'utf8'
  });

  if (result.error) {
    throw result.error;
  }

  if (result.status !== 0) {
    throw new Error(`helmor version exited with status ${result.status}`);
  }

  const actualVersion = result.stdout.trim();
  if (actualVersion !== expectedVersion) {
    throw new Error(`expected helmor ${expectedVersion}, got ${actualVersion}`);
  }
}

async function main() {
  const repo = process.env.HELMOR_REPO || 'helmorx/agent-os';
  const tag = tagFromVersion(process.env.HELMOR_VERSION || pkg.version);
  const target = platformTarget();
  const baseUrl = releaseBaseUrl(repo, tag);
  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'helmor-npm-'));
  const archivePath = path.join(tmpDir, target.asset);
  const checksumsPath = path.join(tmpDir, 'checksums.txt');

  try {
    log(`Downloading ${target.asset} from ${repo}@${tag}`);
    await downloadFile(`${baseUrl}/${target.asset}`, archivePath);
    await downloadFile(`${baseUrl}/checksums.txt`, checksumsPath);

    const expectedChecksum = checksumFor(fs.readFileSync(checksumsPath, 'utf8'), target.asset);
    const actualChecksum = sha256(archivePath);
    if (expectedChecksum !== actualChecksum) {
      throw new Error(`checksum mismatch for ${target.asset}`);
    }

    extractArchive(archivePath, tmpDir, target);

    fs.rmSync(vendorDir, { recursive: true, force: true });
    fs.mkdirSync(vendorDir, { recursive: true });

    const sourceBinary = path.join(tmpDir, target.executable);
    const installedBinary = path.join(vendorDir, target.executable);
    fs.copyFileSync(sourceBinary, installedBinary);

    if (process.platform !== 'win32') {
      fs.chmodSync(installedBinary, 0o755);
    }

    if (tag !== 'latest') {
      verifyInstalledVersion(installedBinary, tag.replace(/^v/, ''));
    }

    log(`Installed helmor to ${installedBinary}`);
  } finally {
    fs.rmSync(tmpDir, { recursive: true, force: true });
  }
}

main().catch((err) => {
  console.error(`helmor: ${err.message}`);
  process.exit(1);
});
