'use strict';

const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');

const installSh = fs.readFileSync(
  path.resolve(__dirname, '..', '..', '..', 'install', 'install.sh'),
  'utf8'
);

test('shell installer keeps release archive filename aligned with checksums', () => {
  assert.match(installSh, /ASSET="helmor_\$\{OS\}_\$\{ARCH\}\.tar\.gz"/);
  assert.match(installSh, /curl -fsSL "\$URL" -o "\$TMP_DIR\/\$ASSET"/);
  assert.match(installSh, /grep "\$ASSET" checksums\.txt \| sha256sum -c -/);
  assert.match(installSh, /grep "\$ASSET" checksums\.txt \| shasum -a 256 -c -/);
  assert.match(installSh, /tar -xzf "\$TMP_DIR\/\$ASSET" -C "\$TMP_DIR"/);
  assert.doesNotMatch(installSh, /helmor\.tar\.gz/);
});
