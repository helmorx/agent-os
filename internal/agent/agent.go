package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/helmorx/agent-os/internal/config"
	"github.com/helmorx/agent-os/internal/project"
)

func InstallAdapters(root string, cfg config.Project) error {
	if err := os.MkdirAll(filepath.Join(root, project.DirName, "adapters"), 0o755); err != nil {
		return err
	}
	for _, adapter := range cfg.Agents {
		if !adapter.Enabled {
			continue
		}
		switch adapter.Name {
		case "codex":
			if err := writeCodex(root); err != nil {
				return err
			}
		case "claude":
			if err := writeClaude(root); err != nil {
				return err
			}
		case "cursor":
			if err := writeRules(root, "cursor", ".cursor/rules/helmor.mdc", cfg); err != nil {
				return err
			}
		case "windsurf":
			if err := writeRules(root, "windsurf", ".windsurf/rules/helmor.md", cfg); err != nil {
				return err
			}
		}
	}
	return nil
}

func AdapterSummary(cfg config.Project) []string {
	var lines []string
	for _, adapter := range cfg.Agents {
		status := "disabled"
		if adapter.Enabled {
			status = "enabled"
		}
		lines = append(lines, fmt.Sprintf("%s: %s (%s)", adapter.Name, status, adapter.Adapter))
	}
	return lines
}

func writeCodex(root string) error {
	path := filepath.Join(root, project.DirName, "adapters", "codex-hooks.json")
	data := map[string]any{
		"hooks": hookMap(),
	}
	return writeJSON(path, data)
}

func writeClaude(root string) error {
	path := filepath.Join(root, project.DirName, "adapters", "claude-settings.local.json")
	data := map[string]any{
		"hooks": hookMap(),
	}
	return writeJSON(path, data)
}

func writeRules(root string, name string, relPath string, cfg config.Project) error {
	path := filepath.Join(root, relPath)
	content := strings.Join([]string{
		"---",
		"description: HELMOR Agent project rules",
		"---",
		"",
		"# HELMOR Agent",
		"",
		"Project: " + cfg.ProjectName,
		"Mode: " + cfg.Mode,
		"Package runner: " + cfg.PackageRunner,
		"Framework: " + cfg.Framework,
		"",
		"Use compact context first, avoid broad repo reads, prefer declared truth files, and run checks before closeout.",
		"Use `helmor status`, `helmor doctor`, and `helmor handoff` when available.",
		"",
		"Adapter: " + name,
	}, "\n")
	return writeFile(path, content+"\n")
}

func hookMap() map[string]any {
	command := "helmor hook --event "
	return map[string]any{
		"SessionStart":     []any{hookEntry("startup|resume|clear|compact", command+"SessionStart")},
		"UserPromptSubmit": []any{hookEntry("*", command+"UserPromptSubmit")},
		"PreToolUse":       []any{hookEntry("*", command+"PreToolUse")},
		"PostToolUse":      []any{hookEntry("*", command+"PostToolUse")},
		"Stop":             []any{hookEntry("*", command+"Stop")},
		"PreCompact":       []any{hookEntry("*", command+"PreCompact")},
		"SessionEnd":       []any{hookEntry("*", command+"SessionEnd")},
	}
}

func hookEntry(matcher string, command string) map[string]any {
	return map[string]any{
		"matcher": matcher,
		"hooks": []any{
			map[string]any{
				"type":    "command",
				"command": command,
				"timeout": 10,
			},
		},
	}
}

func writeJSON(path string, data any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(path, raw, 0o644)
}

func writeFile(path string, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
