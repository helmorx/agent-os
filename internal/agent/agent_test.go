package agent

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/helmorx/agent-os/internal/config"
	"github.com/helmorx/agent-os/internal/project"
)

func TestInstallAdaptersWritesCoreFourOutputs(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"packageManager":"npm@10.0.0"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := project.DefaultConfig(root)
	cfg.Mode = config.ModeObserve

	if err := InstallAdapters(root, cfg); err != nil {
		t.Fatal(err)
	}

	required := []string{
		filepath.Join(root, project.DirName, "adapters", "codex-hooks.json"),
		filepath.Join(root, project.DirName, "adapters", "claude-settings.local.json"),
		filepath.Join(root, ".cursor", "rules", "helmor.mdc"),
		filepath.Join(root, ".windsurf", "rules", "helmor.md"),
	}
	for _, path := range required {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing adapter output %s: %v", path, err)
		}
	}
}
