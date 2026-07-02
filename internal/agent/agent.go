package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/helmorx/agent-os/internal/config"
	"github.com/helmorx/agent-os/internal/project"
)

// legacyWatchScript is a home-relative path fragment, not an absolute path:
// matching on it via strings.Contains lets migration detection work on any
// user's machine, not just the one whose home directory it was written for.
const legacyWatchScript = ".helmor/bin/helmor-watch.py"

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

func InstallGlobalHooks(cfg config.Project) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return InstallGlobalHooksAt(home, cfg)
}

func InstallGlobalHooksAt(home string, cfg config.Project) error {
	for _, adapter := range cfg.Agents {
		if !adapter.Enabled || adapter.Adapter != "hooks" {
			continue
		}
		switch adapter.Name {
		case "codex":
			if err := mergeGlobalHookFile(filepath.Join(home, ".codex", "hooks.json")); err != nil {
				return err
			}
		case "claude":
			// Claude Code only loads ~/.claude/settings.json at the user level;
			// settings.local.json is a project-scoped override file and is never
			// read from the home directory, so hooks written there are inert.
			if err := mergeGlobalHookFile(filepath.Join(home, ".claude", "settings.json")); err != nil {
				return err
			}
		}
	}
	return nil
}

type GlobalHookStatus struct {
	Agent   string
	Path    string
	OK      bool
	Message string
}

func CheckGlobalHooks(cfg config.Project) []GlobalHookStatus {
	home, err := os.UserHomeDir()
	if err != nil {
		return []GlobalHookStatus{{Agent: "global", OK: false, Message: err.Error()}}
	}
	return CheckGlobalHooksAt(home, cfg)
}

func CheckGlobalHooksAt(home string, cfg config.Project) []GlobalHookStatus {
	var statuses []GlobalHookStatus
	for _, adapter := range cfg.Agents {
		if !adapter.Enabled || adapter.Adapter != "hooks" {
			continue
		}
		var path string
		switch adapter.Name {
		case "codex":
			path = filepath.Join(home, ".codex", "hooks.json")
		case "claude":
			path = filepath.Join(home, ".claude", "settings.json")
		default:
			continue
		}
		ok, legacy, err := inspectGlobalHookFile(path)
		status := GlobalHookStatus{Agent: adapter.Name, Path: path, OK: ok && !legacy}
		switch {
		case err != nil:
			status.Message = err.Error()
		case legacy:
			status.Message = "legacy HELMOR Python hook entry remains"
		case !ok:
			status.Message = "official HELMOR hook entries missing"
		default:
			status.Message = "official HELMOR hook entries installed"
		}
		statuses = append(statuses, status)
	}
	return statuses
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
	command := hookCommandPrefix()
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

func hookCommandPrefix() string {
	exe, err := os.Executable()
	if err != nil {
		return "helmor hook --event "
	}
	return hookCommandPrefixForExecutable(exe)
}

func hookCommandPrefixForExecutable(exe string) string {
	if executableBase(exe) != "helmor" && executableBase(exe) != "helmor.exe" {
		return "helmor hook --event "
	}
	if abs, err := filepath.Abs(exe); err == nil {
		exe = abs
	}
	return quoteCommandPath(exe) + " hook --event "
}

func quoteCommandPath(exe string) string {
	if runtime.GOOS == "windows" {
		return `"` + strings.ReplaceAll(exe, `"`, `\"`) + `"`
	}
	return "'" + strings.ReplaceAll(exe, "'", "'\"'\"'") + "'"
}

func mergeGlobalHookFile(path string) error {
	data := map[string]any{}
	if raw, err := os.ReadFile(path); err == nil && len(strings.TrimSpace(string(raw))) > 0 {
		if err := json.Unmarshal(raw, &data); err != nil {
			return fmt.Errorf("load %s: %w", path, err)
		}
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}

	hooks, _ := data["hooks"].(map[string]any)
	if hooks == nil {
		hooks = map[string]any{}
		data["hooks"] = hooks
	}

	for event, entries := range hooks {
		hooks[event] = filterNonHelmorEntries(entries)
	}
	for event, entries := range hookMap() {
		hooks[event] = appendEntries(hooks[event], entries)
	}
	return writeJSON(path, data)
}

func inspectGlobalHookFile(path string) (bool, bool, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false, false, err
	}
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return false, false, err
	}
	hooks, _ := data["hooks"].(map[string]any)
	if hooks == nil {
		return false, false, nil
	}
	allOfficial := true
	legacy := false
	for event := range hookMap() {
		entries, ok := hooks[event].([]any)
		if !ok || !entriesContainOfficial(entries, event) {
			allOfficial = false
		}
		if entriesContainLegacy(entries) {
			legacy = true
		}
	}
	return allOfficial, legacy, nil
}

func filterNonHelmorEntries(entries any) []any {
	list, ok := entries.([]any)
	if !ok {
		return nil
	}
	var kept []any
	for _, entry := range list {
		if !entryHasHelmorCommand(entry) {
			kept = append(kept, entry)
		}
	}
	return kept
}

func appendEntries(existing any, additions any) []any {
	out, _ := existing.([]any)
	newEntries, _ := additions.([]any)
	return append(out, newEntries...)
}

func entriesContainOfficial(entries []any, event string) bool {
	for _, entry := range entries {
		if entryHasCommandMatching(entry, func(value string) bool {
			return isOfficialHelmorHookCommand(value, event)
		}) {
			return true
		}
	}
	return false
}

func entriesContainLegacy(entries []any) bool {
	for _, entry := range entries {
		if entryHasCommandContaining(entry, legacyWatchScript) {
			return true
		}
	}
	return false
}

func entryHasHelmorCommand(entry any) bool {
	return entryHasCommandMatching(entry, func(value string) bool {
		return isHelmorHookCommand(value) || strings.Contains(value, legacyWatchScript)
	})
}

func isOfficialHelmorHookCommand(command string, event string) bool {
	command = strings.TrimSpace(command)
	return !strings.Contains(command, legacyWatchScript) &&
		isHelmorHookCommand(command) &&
		strings.HasSuffix(command, " hook --event "+event)
}

func isHelmorHookCommand(command string) bool {
	before, _, ok := strings.Cut(strings.TrimSpace(command), " hook --event ")
	if !ok {
		return false
	}
	return executableBase(strings.Trim(before, `'"`)) == "helmor" ||
		executableBase(strings.Trim(before, `'"`)) == "helmor.exe"
}

func executableBase(exe string) string {
	return path.Base(strings.ReplaceAll(exe, `\`, `/`))
}

func entryHasCommandContaining(entry any, needle string) bool {
	return entryHasCommandMatching(entry, func(value string) bool { return strings.Contains(value, needle) })
}

func entryHasCommandMatching(entry any, match func(string) bool) bool {
	entryMap, ok := entry.(map[string]any)
	if !ok {
		return false
	}
	rawHooks, ok := entryMap["hooks"].([]any)
	if !ok {
		return false
	}
	for _, hook := range rawHooks {
		hookMap, ok := hook.(map[string]any)
		if !ok {
			continue
		}
		command, _ := hookMap["command"].(string)
		if match(command) {
			return true
		}
	}
	return false
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
