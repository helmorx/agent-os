class Helmor < Formula
  desc "Local-first development engine for AI-assisted coding"
  homepage "https://github.com/helmorx/devsuite"
  version "0.1.0"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/helmorx/devsuite/releases/download/v0.1.0/helmor_darwin_arm64.tar.gz"
      sha256 "85ad2f1e948379c344da0f5963aee78cf57ab3b07ed915362436ee723462b749"
    else
      url "https://github.com/helmorx/devsuite/releases/download/v0.1.0/helmor_darwin_amd64.tar.gz"
      sha256 "69170944ac74d0e4000334024793c004af82546226d9b5db73be47c4ca2042a7"
    end
  end

  def install
    bin.install "helmor"
  end

  test do
    assert_match "0.1.0", shell_output("#{bin}/helmor version")
  end
end
