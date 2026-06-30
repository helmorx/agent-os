package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/helmorx/devsuite/internal/config"
)

const (
	DirName         = ".helmor"
	ConfigFileName  = "project.json"
	StateFileName   = "state.json"
	HandoffFile     = "handoff.md"
	ContextCardFile = "context-card.md"
)

var markerFiles = []string{
	".git",
	"package.json",
	"lock.yaml",
	"pnpm-lock.yaml",
	"yarn.lock",
	"bun.lock",
	"bun.lockb",
	"pyproject.toml",
	"go.mod",
	"Cargo.toml",
	"deno.json",
	"next.config.js",
	"next.config.mjs",
	"vite.config.ts",
	"vite.config.js",
	"apps",
	"services",
	"src",
}

var truthFileCandidates = []string{
	"DOMAIN_TRUTH.md",
	"PRD.md",
	"TRD.md",
	"PRODUCT.md",
	"DESIGN.md",
	"ARCHITECTURE.md",
	"AGENTS.md",
	"CLAUDE.md",
	"README.md",
}

func FindRoot(start string) (string, error) {
	if start == "" {
		start = "."
	}
	abs, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err == nil && !info.IsDir() {
		abs = filepath.Dir(abs)
	}
	for {
		if hasMarker(abs) {
			return abs, nil
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			return "", fmt.Errorf("no project markers found from %s", start)
		}
		abs = parent
	}
}

func ConfigPath(root string) string {
	return filepath.Join(root, DirName, ConfigFileName)
}

func StatePath(root string) string {
	return filepath.Join(root, DirName, StateFileName)
}

func Load(root string) (config.Project, error) {
	path := ConfigPath(root)
	data, err := os.ReadFile(path)
	if err != nil {
		return config.Project{}, err
	}
	var cfg config.Project
	if err := json.Unmarshal(data, &cfg); err != nil {
		return config.Project{}, err
	}
	if cfg.Mode == "" {
		cfg.Mode = config.ModeObserve
	}
	return cfg, nil
}

func Save(root string, cfg config.Project) error {
	if err := os.MkdirAll(filepath.Join(root, DirName), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(ConfigPath(root), data, 0o644)
}

func LoadOrDefault(root string) (config.Project, bool, error) {
	cfg, err := Load(root)
	if err == nil {
		return cfg, true, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return config.Project{}, false, err
	}
	return DefaultConfig(root), false, nil
}

func Init(root string, mode string, force bool) (config.Project, error) {
	if mode == "" {
		mode = config.ModeObserve
	}
	cfgPath := ConfigPath(root)
	if !force {
		if _, err := os.Stat(cfgPath); err == nil {
			return config.Project{}, fmt.Errorf("%s already exists; use --force to overwrite", cfgPath)
		}
	}
	cfg := DefaultConfig(root)
	cfg.Mode = mode
	if err := Save(root, cfg); err != nil {
		return config.Project{}, err
	}
	if err := writeDefaultFiles(root, cfg); err != nil {
		return config.Project{}, err
	}
	return cfg, nil
}

func LoadState(root string) (config.RuntimeState, error) {
	data, err := os.ReadFile(StatePath(root))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return config.RuntimeState{}, nil
		}
		return config.RuntimeState{}, err
	}
	var state config.RuntimeState
	if err := json.Unmarshal(data, &state); err != nil {
		return config.RuntimeState{}, err
	}
	return state, nil
}

func SaveState(root string, state config.RuntimeState) error {
	if err := os.MkdirAll(filepath.Join(root, DirName), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(StatePath(root), data, 0o644)
}

func DefaultConfig(root string) config.Project {
	name := filepath.Base(root)
	runner := DetectPackageRunner(root)
	framework := DetectFramework(root)
	return config.Project{
		SchemaVersion: config.SchemaVersion,
		ProjectName:   name,
		ProjectRoot:   root,
		Mode:          config.ModeObserve,
		PackageRunner: runner,
		Framework:     framework,
		TruthFiles:    DetectTruthFiles(root),
		Checks:        DetectChecks(root, runner),
		Agents: []config.Agent{
			{Name: "codex", Adapter: "hooks", Enabled: true},
			{Name: "claude", Adapter: "hooks", Enabled: true},
			{Name: "cursor", Adapter: "rules", Enabled: true},
			{Name: "windsurf", Adapter: "rules", Enabled: true},
		},
		Tools: config.Tools{
			RTKFirst:           true,
			SQZFallback:        true,
			CodeGraphFirst:     true,
			PackageRunnerGuard: true,
		},
		Policies: config.Policies{
			BlockSecretWrites:        true,
			BlockDestructiveGit:      true,
			BlockUnsafeDeployActions: true,
			RequireHandoffOnStop:     false,
			RequireChecksOnStop:      false,
			DesignDetectors:          true,
		},
		Skills: []string{
			"project-memory",
			"token-reduction",
			"architecture",
			"api-contracts",
			"testing",
			"security",
			"launch-readiness",
			"ui-design",
			"frontend",
			"backend",
			"data",
			"infra",
			"docs-handoff",
			"product-planning",
		},
		Notes: []string{
			"Default mode is observe: HELMOR warns and routes without blocking most actions.",
			"Switch to guard or strict when the project is ready for enforcement.",
		},
	}
}

func DetectPackageRunner(root string) string {
	if packageManager := packageManagerField(root); packageManager != "" {
		switch {
		case strings.HasPrefix(packageManager, "nub@"):
			return "nub"
		case strings.HasPrefix(packageManager, "pnpm@"):
			return "pnpm"
		case strings.HasPrefix(packageManager, "yarn@"):
			return "yarn"
		case strings.HasPrefix(packageManager, "bun@"):
			return "bun"
		case strings.HasPrefix(packageManager, "npm@"):
			return "npm"
		}
	}
	switch {
	case exists(filepath.Join(root, "lock.yaml")):
		return "nub"
	case exists(filepath.Join(root, "pnpm-lock.yaml")):
		return "pnpm"
	case exists(filepath.Join(root, "yarn.lock")):
		return "yarn"
	case exists(filepath.Join(root, "bun.lock")) || exists(filepath.Join(root, "bun.lockb")):
		return "bun"
	case exists(filepath.Join(root, "package-lock.json")) || exists(filepath.Join(root, "package.json")):
		return "npm"
	case exists(filepath.Join(root, "pyproject.toml")):
		return "python"
	case exists(filepath.Join(root, "go.mod")):
		return "go"
	case exists(filepath.Join(root, "Cargo.toml")):
		return "cargo"
	default:
		return "unknown"
	}
}

func DetectFramework(root string) string {
	switch {
	case exists(filepath.Join(root, "next.config.js")) || exists(filepath.Join(root, "next.config.mjs")) || exists(filepath.Join(root, "next.config.ts")):
		return "next"
	case exists(filepath.Join(root, "vite.config.ts")) || exists(filepath.Join(root, "vite.config.js")):
		return "vite"
	case exists(filepath.Join(root, "go.mod")):
		return "go"
	case exists(filepath.Join(root, "Cargo.toml")):
		return "rust"
	case exists(filepath.Join(root, "pyproject.toml")):
		return "python"
	case exists(filepath.Join(root, "apps")) && exists(filepath.Join(root, "services")):
		return "monorepo"
	default:
		return "generic"
	}
}

func DetectTruthFiles(root string) []string {
	var files []string
	for _, candidate := range truthFileCandidates {
		if exists(filepath.Join(root, candidate)) {
			files = append(files, candidate)
		}
	}
	return files
}

func DetectChecks(root string, runner string) []config.Check {
	var checks []config.Check
	if exists(filepath.Join(root, "package.json")) {
		scripts := packageScripts(root)
		for _, name := range []string{"lint", "typecheck", "test", "build", "audit"} {
			if scripts[name] {
				checks = append(checks, config.Check{Name: name, Command: runnerCommand(runner, name)})
			}
		}
	}
	if exists(filepath.Join(root, "go.mod")) {
		checks = append(checks, config.Check{Name: "test", Command: "go test ./..."})
	}
	if exists(filepath.Join(root, "Cargo.toml")) {
		checks = append(checks, config.Check{Name: "test", Command: "cargo test"})
	}
	if exists(filepath.Join(root, "pyproject.toml")) {
		checks = append(checks, config.Check{Name: "test", Command: "python -m pytest"})
	}
	return dedupeChecks(checks)
}

func ContextCard(cfg config.Project) string {
	lines := []string{
		"# HELMOR Context Card",
		"",
		fmt.Sprintf("- Project: %s", cfg.ProjectName),
		fmt.Sprintf("- Mode: %s", cfg.Mode),
		fmt.Sprintf("- Framework: %s", cfg.Framework),
		fmt.Sprintf("- Package runner: %s", cfg.PackageRunner),
	}
	if len(cfg.TruthFiles) > 0 {
		lines = append(lines, fmt.Sprintf("- Truth files: %s", strings.Join(cfg.TruthFiles, " > ")))
	}
	if cfg.Tools.RTKFirst {
		lines = append(lines, "- Prefer rtk for compressed shell/git/test output.")
	}
	if cfg.Tools.SQZFallback {
		lines = append(lines, "- Use sqz as fallback compression/handoff support when available.")
	}
	if cfg.Tools.CodeGraphFirst {
		lines = append(lines, "- Use code graph/MCP discovery before broad repo scans.")
	}
	return strings.Join(lines, "\n") + "\n"
}

func WriteContextCard(root string, cfg config.Project) error {
	return os.WriteFile(filepath.Join(root, DirName, ContextCardFile), []byte(ContextCard(cfg)), 0o644)
}

func writeDefaultFiles(root string, cfg config.Project) error {
	if err := os.MkdirAll(filepath.Join(root, DirName), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(root, DirName, ".gitignore"), []byte("state.json\nhandoff.md\ncontext-card.md\nsession-log.jsonl\nlast-hook-output.json\n"), 0o644); err != nil {
		return err
	}
	return WriteContextCard(root, cfg)
}

func hasMarker(root string) bool {
	for _, marker := range markerFiles {
		if exists(filepath.Join(root, marker)) {
			return true
		}
	}
	return false
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func runnerCommand(runner string, script string) string {
	switch runner {
	case "nub":
		return "nub run " + script
	case "pnpm":
		return "pnpm run " + script
	case "yarn":
		return "yarn " + script
	case "bun":
		return "bun run " + script
	case "npm":
		return "npm run " + script
	default:
		return script
	}
}

func packageManagerField(root string) string {
	var data map[string]any
	if err := readPackageJSON(root, &data); err != nil {
		return ""
	}
	if value, ok := data["packageManager"].(string); ok {
		return value
	}
	return ""
}

func packageScripts(root string) map[string]bool {
	var data map[string]any
	scripts := map[string]bool{}
	if err := readPackageJSON(root, &data); err != nil {
		return scripts
	}
	raw, ok := data["scripts"].(map[string]any)
	if !ok {
		return scripts
	}
	for name := range raw {
		scripts[name] = true
	}
	return scripts
}

func readPackageJSON(root string, target any) error {
	data, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func dedupeChecks(checks []config.Check) []config.Check {
	seen := map[string]bool{}
	var out []config.Check
	for _, check := range checks {
		if seen[check.Name] {
			continue
		}
		seen[check.Name] = true
		out = append(out, check)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
