package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/helmorx/agent-os/internal/agent"
	"github.com/helmorx/agent-os/internal/config"
	"github.com/helmorx/agent-os/internal/detector"
	"github.com/helmorx/agent-os/internal/hook"
	"github.com/helmorx/agent-os/internal/project"
	"github.com/helmorx/agent-os/internal/ui"
)

const version = "0.1.2"

func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stdout)
		return 0
	}
	switch args[0] {
	case "-h", "--help", "help":
		usage(stdout)
		return 0
	case "version":
		fmt.Fprintln(stdout, version)
		return 0
	case "init":
		return runInit(args[1:], stdout, stderr)
	case "install":
		return runInstall(args[1:], stdout, stderr)
	case "uninstall":
		return runUninstall(args[1:], stdout, stderr)
	case "status":
		return runStatus(args[1:], stdout, stderr)
	case "doctor":
		return runDoctor(args[1:], stdout, stderr)
	case "dashboard":
		return runDashboard(args[1:], stdout, stderr)
	case "checks":
		return runChecks(args[1:], stdout, stderr)
	case "handoff":
		return runHandoff(args[1:], stdout, stderr)
	case "reduce-tokens":
		return runReduceTokens(args[1:], stdout, stderr)
	case "verify":
		return runDoctor(args[1:], stdout, stderr)
	case "security":
		return runSecurity(args[1:], stdout, stderr)
	case "design":
		return runDesign(args[1:], stdout, stderr)
	case "task":
		return runTask(args[1:], stdout, stderr)
	case "hook":
		return runHook(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		usage(stderr)
		return 2
	}
}

func usage(writer io.Writer) {
	fmt.Fprintln(writer, "HELMOR Agent OS")
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "Usage: helmor <command> [options]")
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "Commands:")
	fmt.Fprintln(writer, "  init              create .helmor/project.json in the current project")
	fmt.Fprintln(writer, "  install           init project and generate agent adapters")
	fmt.Fprintln(writer, "  uninstall         remove .helmor from the current project")
	fmt.Fprintln(writer, "  status            show compact project status")
	fmt.Fprintln(writer, "  doctor            run deterministic checks")
	fmt.Fprintln(writer, "  dashboard         show terminal dashboard")
	fmt.Fprintln(writer, "  checks            list detected project checks")
	fmt.Fprintln(writer, "  handoff           write .helmor/handoff.md")
	fmt.Fprintln(writer, "  reduce-tokens     show token-saving recommendations")
	fmt.Fprintln(writer, "  security          run security detector pack")
	fmt.Fprintln(writer, "  design audit      run UI/design detector pack")
	fmt.Fprintln(writer, "  hook              agent hook entrypoint")
}

func runInit(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(stderr)
	mode := fs.String("mode", config.ModeObserve, "observe, guard, or strict")
	force := fs.Bool("force", false, "overwrite existing project profile")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	root, err := project.FindRoot(".")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	cfg, err := project.Init(root, *mode, *force)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Initialized HELMOR for %s in %s mode\n", cfg.ProjectName, cfg.Mode)
	return 0
}

func runInstall(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	fs.SetOutput(stderr)
	mode := fs.String("mode", config.ModeObserve, "observe, guard, or strict")
	force := fs.Bool("force", false, "overwrite existing project profile")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	root, err := project.FindRoot(".")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	cfg, err := project.Init(root, *mode, *force)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			cfg, err = project.Load(root)
		}
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	if err := agent.InstallAdapters(root, cfg); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Installed HELMOR adapters for %s\n", cfg.ProjectName)
	return 0
}

func runUninstall(_ []string, stdout io.Writer, stderr io.Writer) int {
	root, err := project.FindRoot(".")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := os.RemoveAll(filepath.Join(root, project.DirName)); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintln(stdout, "Removed .helmor")
	return 0
}

func runStatus(_ []string, stdout io.Writer, stderr io.Writer) int {
	root, cfg, err := loadProject(".")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Project: %s\nRoot: %s\nMode: %s\nFramework: %s\nPackage runner: %s\n", cfg.ProjectName, root, cfg.Mode, cfg.Framework, cfg.PackageRunner)
	if len(cfg.TruthFiles) > 0 {
		fmt.Fprintf(stdout, "Truth files: %s\n", strings.Join(cfg.TruthFiles, ", "))
	}
	fmt.Fprintf(stdout, "Checks: %d\n", len(cfg.Checks))
	return 0
}

func runDoctor(_ []string, stdout io.Writer, stderr io.Writer) int {
	root, cfg, err := loadProject(".")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	findings := detector.Run(root, cfg)
	if len(findings) == 0 {
		fmt.Fprintln(stdout, "HELMOR doctor: ok")
		return 0
	}
	for _, finding := range findings {
		fmt.Fprintf(stdout, "[%s] %s %s %s\n", finding.Severity, finding.Rule, finding.Path, finding.Message)
	}
	return 0
}

func runDashboard(_ []string, stdout io.Writer, stderr io.Writer) int {
	root, cfg, err := loadProject(".")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	ui.Dashboard(stdout, cfg, detector.Run(root, cfg))
	return 0
}

func runChecks(_ []string, stdout io.Writer, stderr io.Writer) int {
	_, cfg, err := loadProject(".")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if len(cfg.Checks) == 0 {
		fmt.Fprintln(stdout, "No checks detected.")
		return 0
	}
	for _, check := range cfg.Checks {
		fmt.Fprintf(stdout, "%s\t%s\n", check.Name, check.Command)
	}
	return 0
}

func runHandoff(args []string, stdout io.Writer, stderr io.Writer) int {
	root, cfg, err := loadProject(".")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	state, _ := project.LoadState(root)
	path := filepath.Join(root, project.DirName, project.HandoffFile)
	body := fmt.Sprintf("# HELMOR Handoff\n\n- Project: %s\n- Mode: %s\n- Stage: %s\n- Pending checks: %s\n", cfg.ProjectName, cfg.Mode, state.LastStage, strings.Join(state.PendingChecks, ", "))
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	state.HandoffSaved = true
	state.LastCompactionSummaryPath = path
	_ = project.SaveState(root, state)
	fmt.Fprintln(stdout, path)
	return 0
}

func runReduceTokens(_ []string, stdout io.Writer, stderr io.Writer) int {
	_, cfg, err := loadProject(".")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintln(stdout, "Token reduction recommendations:")
	if cfg.Tools.RTKFirst {
		fmt.Fprintln(stdout, "- Use rtk before shell/git/test commands when available.")
	}
	if cfg.Tools.SQZFallback {
		fmt.Fprintln(stdout, "- Use sqz as fallback for compressed outputs and handoffs.")
	}
	if cfg.Tools.CodeGraphFirst {
		fmt.Fprintln(stdout, "- Use code graph/MCP discovery before broad grep/read-all behavior.")
	}
	fmt.Fprintf(stdout, "- Use `%s` consistently for package commands.\n", cfg.PackageRunner)
	fmt.Fprintln(stdout, "- Keep .helmor/context-card.md current to avoid repeated repo rediscovery.")
	return 0
}

func runSecurity(_ []string, stdout io.Writer, stderr io.Writer) int {
	root, cfg, err := loadProject(".")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	for _, finding := range detector.Run(root, cfg) {
		if finding.Rule == "secret-shaped-filename" || finding.Severity == detector.Block {
			fmt.Fprintf(stdout, "[%s] %s %s %s\n", finding.Severity, finding.Rule, finding.Path, finding.Message)
		}
	}
	return 0
}

func runDesign(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: helmor design <init|audit|polish>")
		return 2
	}
	root, cfg, err := loadProject(".")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	switch args[0] {
	case "init":
		path := filepath.Join(root, "DESIGN.md")
		if _, err := os.Stat(path); err == nil {
			fmt.Fprintln(stdout, "DESIGN.md already exists")
			return 0
		}
		content := "# Design System\n\nAudience:\n\nBrand voice:\n\nColor rules:\n\nTypography rules:\n\nInteraction rules:\n"
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "Created DESIGN.md")
	case "audit", "polish":
		for _, finding := range detector.Run(root, cfg) {
			if strings.HasPrefix(finding.Rule, "design.") {
				fmt.Fprintf(stdout, "[%s] %s %s %s\n", finding.Severity, finding.Rule, finding.Path, finding.Message)
			}
		}
	default:
		fmt.Fprintln(stderr, "usage: helmor design <init|audit|polish>")
		return 2
	}
	return 0
}

func runTask(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: helmor task <start|finish>")
		return 2
	}
	root, _, err := loadProject(".")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	state, _ := project.LoadState(root)
	switch args[0] {
	case "start":
		stage := "task"
		if len(args) > 1 {
			stage = strings.Join(args[1:], " ")
		}
		state.LastStage = stage
		if err := project.SaveState(root, state); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "Started HELMOR task:", stage)
	case "finish":
		return runHandoff(nil, stdout, stderr)
	default:
		fmt.Fprintln(stderr, "usage: helmor task <start|finish>")
		return 2
	}
	return 0
}

func runHook(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("hook", flag.ContinueOnError)
	fs.SetOutput(stderr)
	event := fs.String("event", "", "hook event name")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	output, code := hook.Handle(os.Stdin, *event)
	if err := hook.Encode(output, stdout); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return code
}

func loadProject(start string) (string, config.Project, error) {
	root, err := project.FindRoot(start)
	if err != nil {
		return "", config.Project{}, err
	}
	cfg, _, err := project.LoadOrDefault(root)
	return root, cfg, err
}

func checksum(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

var _ = checksum
var _ = writeJSON
