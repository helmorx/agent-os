package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/helmorx/devsuite/internal/config"
)

func TestInitCreatesObserveProfile(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "package.json"), `{"packageManager":"nub@0.1.14","scripts":{"test":"vitest run","lint":"eslint ."}}`)
	write(t, filepath.Join(root, "README.md"), "# Test\n")

	cfg, err := Init(root, config.ModeObserve, false)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Mode != config.ModeObserve {
		t.Fatalf("mode = %s", cfg.Mode)
	}
	if cfg.PackageRunner != "nub" {
		t.Fatalf("runner = %s", cfg.PackageRunner)
	}
	if len(cfg.Checks) != 2 {
		t.Fatalf("checks = %#v", cfg.Checks)
	}
	if _, err := os.Stat(filepath.Join(root, DirName, ConfigFileName)); err != nil {
		t.Fatal(err)
	}
}

func TestFindRootWalksUp(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "go.mod"), "module example.com/test\n")
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	found, err := FindRoot(nested)
	if err != nil {
		t.Fatal(err)
	}
	if found != root {
		t.Fatalf("found %s, want %s", found, root)
	}
}

func write(t *testing.T, path string, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}
