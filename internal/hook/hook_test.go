package hook

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/helmorx/devsuite/internal/config"
	"github.com/helmorx/devsuite/internal/project"
)

func TestObserveModeWarnsButAllowsRunnerBypass(t *testing.T) {
	root := setupProject(t, config.ModeObserve)
	input := map[string]any{
		"cwd":             root,
		"hook_event_name": "PreToolUse",
		"tool_name":       "Bash",
		"tool_input": map[string]any{
			"command": "npm install",
		},
	}

	out, code := Handle(jsonInput(t, input), "")
	if code != 0 {
		t.Fatalf("code = %d", code)
	}
	if !out.Continue {
		t.Fatal("observe should continue")
	}
	if !strings.Contains(out.SystemMessage, "package runner bypass") {
		t.Fatalf("message = %q", out.SystemMessage)
	}
}

func TestGuardModeBlocksSecretWrite(t *testing.T) {
	root := setupProject(t, config.ModeGuard)
	input := map[string]any{
		"cwd":             root,
		"hook_event_name": "PreToolUse",
		"tool_name":       "Write",
		"tool_input": map[string]any{
			"file_path": ".env",
		},
	}

	out, code := Handle(jsonInput(t, input), "")
	if code != 2 {
		t.Fatalf("code = %d", code)
	}
	if out.Continue {
		t.Fatal("guard should block")
	}
}

func TestStrictModeBlocksStopWhenSecurityReviewIncomplete(t *testing.T) {
	root := setupProject(t, config.ModeStrict)
	state := config.RuntimeState{SecurityReviewRequired: true}
	if err := project.SaveState(root, state); err != nil {
		t.Fatal(err)
	}

	out, code := Handle(jsonInput(t, map[string]any{
		"cwd":             root,
		"hook_event_name": "Stop",
	}), "")
	if code != 2 {
		t.Fatalf("code = %d", code)
	}
	if out.Decision != "block" {
		t.Fatalf("decision = %s", out.Decision)
	}
}

func TestEncodeEmitsJSON(t *testing.T) {
	var buffer bytes.Buffer
	if err := Encode(Output{Continue: true, SystemMessage: "ok"}, &buffer); err != nil {
		t.Fatal(err)
	}
	if !json.Valid(buffer.Bytes()) {
		t.Fatalf("invalid json: %s", buffer.String())
	}
}

func setupProject(t *testing.T, mode string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"packageManager":"nub@0.1.14"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := project.Init(root, mode, false)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Policies.RequireChecksOnStop = false
	cfg.Policies.RequireHandoffOnStop = false
	if err := project.Save(root, cfg); err != nil {
		t.Fatal(err)
	}
	return root
}

func jsonInput(t *testing.T, value any) *bytes.Reader {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(data)
}
