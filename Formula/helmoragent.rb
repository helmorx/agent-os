class Helmoragent < Formula
  desc "Agent watcher for AI-assisted product development"
  homepage "https://github.com/helmorx/agent-os"
  url "https://github.com/helmorx/agent-os/archive/refs/tags/v0.1.1.tar.gz"
  sha256 "7df75614949a7246e610ba089845737a8c100bd120fe57ef1fe922a5406fa4de"
  license "Apache-2.0"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(output: bin/"helmor"), "./cmd/helmor"
  end

  test do
    assert_match "0.1.1", shell_output("#{bin}/helmor version")
    assert_match "HELMOR Agent OS", shell_output("#{bin}/helmor help")
  end
end
