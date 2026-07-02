'use strict';

const test = require('node:test');
const assert = require('node:assert/strict');

const {
  checksumFor,
  platformTarget,
  releaseBaseUrl,
  tagFromVersion
} = require('../scripts/lib');

test('platformTarget maps macOS arm64 release asset', () => {
  assert.deepEqual(platformTarget('darwin', 'arm64'), {
    goos: 'darwin',
    goarch: 'arm64',
    archiveExt: 'tar.gz',
    executable: 'helmor',
    asset: 'helmor_darwin_arm64.tar.gz'
  });
});

test('platformTarget maps Windows x64 release asset', () => {
  assert.deepEqual(platformTarget('win32', 'x64'), {
    goos: 'windows',
    goarch: 'amd64',
    archiveExt: 'zip',
    executable: 'helmor.exe',
    asset: 'helmor_windows_amd64.zip'
  });
});

test('platformTarget rejects unsupported platforms and architectures', () => {
  assert.throws(() => platformTarget('freebsd', 'x64'), /unsupported platform/);
  assert.throws(() => platformTarget('linux', 'ia32'), /unsupported architecture/);
});

test('tagFromVersion normalizes npm versions to git tags', () => {
  assert.equal(tagFromVersion('0.1.2'), 'v0.1.2');
  assert.equal(tagFromVersion('v0.1.2'), 'v0.1.2');
  assert.equal(tagFromVersion('latest'), 'latest');
});

test('releaseBaseUrl builds GitHub release download URLs', () => {
  assert.equal(
    releaseBaseUrl('helmorx/helmoragent', 'latest'),
    'https://github.com/helmorx/helmoragent/releases/latest/download'
  );
  assert.equal(
    releaseBaseUrl('helmorx/helmoragent', 'v0.1.2'),
    'https://github.com/helmorx/helmoragent/releases/download/v0.1.2'
  );
});

test('checksumFor reads sha256sum output', () => {
  const checksums = [
    'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa  helmor_darwin_amd64.tar.gz',
    'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb  helmor_linux_arm64.tar.gz'
  ].join('\n');

  assert.equal(
    checksumFor(checksums, 'helmor_linux_arm64.tar.gz'),
    'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb'
  );
  assert.throws(() => checksumFor(checksums, 'missing.tar.gz'), /missing checksum/);
});
