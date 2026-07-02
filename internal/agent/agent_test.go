package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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

func TestInstallGlobalHooksMergesMigratesAndIsIdempotent(t *testing.T) {
	home := t.TempDir()
	codexPath := filepath.Join(home, ".codex", "hooks.json")
	claudePath := filepath.Join(home, ".claude", "settings.json")
	write(t, codexPath, `{
  "hooks": {
    "SessionStart": [
      {
        "matcher": "startup",
        "hooks": [
          {
            "type": "command",
            "command": "echo keep-me"
          }
        ]
      }
    ],
    "UserPromptSubmit": [
      {
        "matcher": "*",
        "hooks": [
          {
            "type": "command",
            "command": "python3 /Users/x/Documents/.helmor/bin/helmor-watch.py hook --event UserPromptSubmit"
          }
        ]
      }
    ]
  }
}`)
	write(t, claudePath, `{"hooks":{}}`)

	root := t.TempDir()
	write(t, filepath.Join(root, "package.json"), `{"packageManager":"nub@0.1.14"}`)
	cfg := project.DefaultConfig(root)

	if err := InstallGlobalHooksAt(home, cfg); err != nil {
		t.Fatal(err)
	}
	if err := InstallGlobalHooksAt(home, cfg); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{codexPath, claudePath} {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		source := string(raw)
		if strings.Contains(source, "/Users/x/Documents/.helmor/bin/helmor-watch.py") {
			t.Fatalf("legacy hook remained in %s: %s", path, source)
		}
		for event := range hookMap() {
			if countOfficialCommands(t, raw, event) != 1 {
				t.Fatalf("%s official command count for %s = %d", path, event, countOfficialCommands(t, raw, event))
			}
		}
	}

	codexRaw, err := os.ReadFile(codexPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(codexRaw), "echo keep-me") {
		t.Fatalf("unrelated hook was not preserved: %s", string(codexRaw))
	}
}

func TestCheckGlobalHooksReportsLegacy(t *testing.T) {
	home := t.TempDir()
	path := filepath.Join(home, ".codex", "hooks.json")
	write(t, path, `{"hooks":{"SessionStart":[{"matcher":"startup","hooks":[{"type":"command","command":"python3 /Users/x/Documents/.helmor/bin/helmor-watch.py hook --event SessionStart"}]}]}}`)
	root := t.TempDir()
	write(t, filepath.Join(root, "package.json"), `{"packageManager":"nub@0.1.14"}`)
	cfg := project.DefaultConfig(root)
	cfg.Agents = []config.Agent{{Name: "codex", Adapter: "hooks", Enabled: true}}

	statuses := CheckGlobalHooksAt(home, cfg)
	if len(statuses) != 1 {
		t.Fatalf("statuses = %#v", statuses)
	}
	if statuses[0].OK || !strings.Contains(statuses[0].Message, "legacy") {
		t.Fatalf("status = %#v", statuses[0])
	}
}

func countOfficialCommands(t *testing.T, raw []byte, event string) int {
	t.Helper()
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		t.Fatal(err)
	}
	hooks, _ := data["hooks"].(map[string]any)
	entries, _ := hooks[event].([]any)
	want := "helmor hook --event " + event
	count := 0
	for _, entry := range entries {
		entryMap, _ := entry.(map[string]any)
		rawHooks, _ := entryMap["hooks"].([]any)
		for _, hook := range rawHooks {
			hookMap, _ := hook.(map[string]any)
			if hookMap["command"] == want {
				count++
			}
		}
	}
	return count
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
