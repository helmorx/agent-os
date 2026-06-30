package detector

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/helmorx/agent-os/internal/config"
)

type Severity string

const (
	Info  Severity = "info"
	Warn  Severity = "warn"
	Block Severity = "block"
)

type Finding struct {
	Rule     string   `json:"rule"`
	Severity Severity `json:"severity"`
	Path     string   `json:"path,omitempty"`
	Message  string   `json:"message"`
}

var secretNamePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(^|/)\.env($|\.)`),
	regexp.MustCompile(`(?i)\.secrets(/|$)`),
	regexp.MustCompile(`(?i)(secret|credential|private.*key|auth.*token|access.*token|refresh.*token|mnemonic)`),
	regexp.MustCompile(`(?i)(backup.*\.json|keystore.*\.(json|txt|env|key|pem|bak|enc)|wallet.*\.(key|pem|json))`),
	regexp.MustCompile(`(?i)(^|/)(id_rsa|id_dsa|id_ecdsa|id_ed25519)$`),
	regexp.MustCompile(`(?i)\.(pem|key|p12|pfx)$`),
}

var designPatterns = []struct {
	rule    string
	pattern *regexp.Regexp
}{
	{"design.gradient-text", regexp.MustCompile(`(?i)gradient.*text|background-clip:\s*text|-webkit-background-clip:\s*text`)},
	{"design.dark-glow", regexp.MustCompile(`(?i)box-shadow:\s*0\s+0|drop-shadow|glow`)},
	{"design.overused-rounded", regexp.MustCompile(`(?i)border-radius:\s*(2[4-9]|[3-9][0-9])px|rounded-full|rounded-\[`)},
	{"design.negative-tracking", regexp.MustCompile(`(?i)letter-spacing:\s*-|tracking-\[-`)},
}

func Run(root string, cfg config.Project) []Finding {
	var findings []Finding
	findings = append(findings, detectTruthFiles(root, cfg)...)
	findings = append(findings, detectTooling(cfg)...)
	findings = append(findings, scanFiles(root, cfg)...)
	return findings
}

func IsSecretPath(path string) bool {
	normalized := filepath.ToSlash(path)
	for _, pattern := range secretNamePatterns {
		if pattern.MatchString(normalized) {
			if strings.HasSuffix(strings.ToLower(normalized), ".env.example") || strings.HasSuffix(strings.ToLower(normalized), ".env.staging.example") {
				return false
			}
			return true
		}
	}
	return false
}

func IsDestructiveGit(command string) bool {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\bgit\s+reset\s+--hard\b`),
		regexp.MustCompile(`\bgit\s+checkout\s+--\b`),
		regexp.MustCompile(`\bgit\s+clean\s+-[^\n;|&]*[fd]`),
		regexp.MustCompile(`\bgit\s+push\b[^\n;|&]*\s--force(?:-with-lease)?\b`),
	}
	for _, pattern := range patterns {
		if pattern.MatchString(command) {
			return true
		}
	}
	return false
}

func IsUnsafeDeploy(command string) bool {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)\bmainnet\b`),
		regexp.MustCompile(`(?i)\bproduction\b`),
		regexp.MustCompile(`(?i)--prod\b`),
		regexp.MustCompile(`(?i)\brailway\s+(up|deploy)\b`),
		regexp.MustCompile(`(?i)\bvercel\s+(deploy\s+)?--prod\b`),
		regexp.MustCompile(`(?i)\bterraform\s+apply\b`),
		regexp.MustCompile(`(?i)\bkubectl\s+apply\b`),
	}
	for _, pattern := range patterns {
		if pattern.MatchString(command) {
			return true
		}
	}
	return false
}

func IsRunnerBypass(command string, runner string) bool {
	if runner == "" || runner == "unknown" {
		return false
	}
	pattern := regexp.MustCompile(`(^|[;&|]\s*)(npm|pnpm|yarn|bun|npx|nub)\s+(install|run|test|exec|dlx|add|remove)\b`)
	match := pattern.FindStringSubmatch(command)
	if len(match) < 3 {
		return false
	}
	used := match[2]
	if runner == "yarn" && used == "yarn" {
		return false
	}
	return used != runner
}

func IsBroadSearch(command string) bool {
	if strings.Contains(command, "codebase-memory") || strings.Contains(command, "mcp__codebase_memory") {
		return false
	}
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(^|[;&|]\s*)(rg|grep)\b`),
		regexp.MustCompile(`\bgrep\b[^\n;|&]*\s-R\b`),
		regexp.MustCompile(`(^|[;&|]\s*)find\s+(\.|apps|packages|services|src|tests)\b`),
	}
	for _, pattern := range patterns {
		if pattern.MatchString(command) {
			return true
		}
	}
	return false
}

func detectTruthFiles(root string, cfg config.Project) []Finding {
	if len(cfg.TruthFiles) > 0 {
		return nil
	}
	return []Finding{{
		Rule:     "truth-files.missing",
		Severity: Warn,
		Message:  "No truth files detected. Add PRODUCT.md, DESIGN.md, ARCHITECTURE.md, AGENTS.md, or README.md to reduce drift.",
	}}
}

func detectTooling(cfg config.Project) []Finding {
	var findings []Finding
	if cfg.Tools.RTKFirst && !binaryAvailable("rtk") {
		findings = append(findings, Finding{Rule: "tool.rtk.missing", Severity: Info, Message: "rtk is not available; shell/test output may consume more tokens."})
	}
	if cfg.Tools.SQZFallback && !binaryAvailable("sqz") {
		findings = append(findings, Finding{Rule: "tool.sqz.missing", Severity: Info, Message: "sqz is not available as fallback compression."})
	}
	if cfg.PackageRunner == "nub" && !binaryAvailable("nub") {
		findings = append(findings, Finding{Rule: "tool.nub.missing", Severity: Warn, Message: "project declares nub, but nub is not available on PATH."})
	}
	return findings
}

func scanFiles(root string, cfg config.Project) []Finding {
	var findings []Finding
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := entry.Name()
		if entry.IsDir() {
			switch name {
			case ".git", "node_modules", ".next", "dist", "coverage", "target", ".helmor":
				return filepath.SkipDir
			}
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			rel = path
		}
		if IsSecretPath(rel) {
			findings = append(findings, Finding{Rule: "secret-shaped-filename", Severity: Block, Path: rel, Message: "secret-shaped filename should not be committed or written by agents."})
			return nil
		}
		if cfg.Policies.DesignDetectors && isSourceLike(rel) {
			data, readErr := os.ReadFile(path)
			if readErr == nil {
				source := string(data)
				for _, rule := range designPatterns {
					if rule.pattern.MatchString(source) {
						findings = append(findings, Finding{Rule: rule.rule, Severity: Warn, Path: rel, Message: "design detector matched; review for generic AI UI drift."})
					}
				}
			}
		}
		return nil
	})
	return findings
}

func isSourceLike(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasSuffix(lower, ".css") ||
		strings.HasSuffix(lower, ".tsx") ||
		strings.HasSuffix(lower, ".jsx") ||
		strings.HasSuffix(lower, ".html") ||
		strings.HasSuffix(lower, ".vue") ||
		strings.HasSuffix(lower, ".svelte")
}

func binaryAvailable(name string) bool {
	pathEnv := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(pathEnv) {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			return true
		}
	}
	return false
}
