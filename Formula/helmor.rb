class Helmor < Formula
  desc "Local-first development engine for AI-assisted coding"
  homepage "https://github.com/helmorx/devsuite"
  version "0.1.0"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/helmorx/devsuite/releases/download/v0.1.0/helmor_darwin_arm64.tar.gz"
      sha256 "REPLACE_WITH_DARWIN_ARM64_SHA256"
    else
      url "https://github.com/helmorx/devsuite/releases/download/v0.1.0/helmor_darwin_amd64.tar.gz"
      sha256 "REPLACE_WITH_DARWIN_AMD64_SHA256"
    end
  end

  def install
    bin.install "helmor"
  end

  test do
    assert_match "0.1.0", shell_output("#{bin}/helmor version")
  end
end
