'use strict';

const SUPPORTED_PLATFORMS = {
  darwin: 'darwin',
  linux: 'linux',
  win32: 'windows'
};

function normalizeArch(arch) {
  switch (arch) {
    case 'arm64':
      return 'arm64';
    case 'x64':
      return 'amd64';
    default:
      throw new Error(`unsupported architecture: ${arch}`);
  }
}

function platformTarget(platform = process.platform, arch = process.arch) {
  const goos = SUPPORTED_PLATFORMS[platform];
  if (!goos) {
    throw new Error(`unsupported platform: ${platform}`);
  }

  const goarch = normalizeArch(arch);
  const archiveExt = platform === 'win32' ? 'zip' : 'tar.gz';
  const executable = platform === 'win32' ? 'helmor.exe' : 'helmor';

  return {
    goos,
    goarch,
    archiveExt,
    executable,
    asset: `helmor_${goos}_${goarch}.${archiveExt}`
  };
}

function tagFromVersion(version) {
  if (!version) {
    throw new Error('version is required');
  }

  if (version === 'latest' || version.startsWith('v')) {
    return version;
  }

  return `v${version}`;
}

function releaseBaseUrl(repo, tag) {
  if (!repo) {
    throw new Error('repo is required');
  }

  if (tag === 'latest') {
    return `https://github.com/${repo}/releases/latest/download`;
  }

  return `https://github.com/${repo}/releases/download/${tag}`;
}

function checksumFor(checksums, asset) {
  for (const line of checksums.split(/\r?\n/)) {
    const parts = line.trim().split(/\s+/);
    if (parts.length >= 2 && parts[parts.length - 1] === asset) {
      return parts[0].toLowerCase();
    }
  }

  throw new Error(`missing checksum for ${asset}`);
}

module.exports = {
  checksumFor,
  platformTarget,
  releaseBaseUrl,
  tagFromVersion
};
