package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitStatusAndDashboard(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "go.mod"), "module example.com/test\n")
	oldWD, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := Run([]string{"init"}, &stdout, &stderr); code != 0 {
		t.Fatalf("init code=%d stderr=%s", code, stderr.String())
	}
	stdout.Reset()
	if code := Run([]string{"status"}, &stdout, &stderr); code != 0 {
		t.Fatalf("status code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Package runner: go") {
		t.Fatalf("status output = %s", stdout.String())
	}
	stdout.Reset()
	if code := Run([]string{"dashboard"}, &stdout, &stderr); code != 0 {
		t.Fatalf("dashboard code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "HELMOR Agent OS") {
		t.Fatalf("dashboard output = %s", stdout.String())
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
