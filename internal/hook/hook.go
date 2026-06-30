package hook

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/helmorx/agent-os/internal/config"
	"github.com/helmorx/agent-os/internal/detector"
	"github.com/helmorx/agent-os/internal/project"
)

type Input struct {
	CWD           string         `json:"cwd"`
	HookEventName string         `json:"hook_event_name"`
	ToolName      string         `json:"tool_name"`
	ToolInput     map[string]any `json:"tool_input"`
	UserPrompt    string         `json:"user_prompt"`
	Reason        string         `json:"reason"`
}

type Output struct {
	Continue           bool           `json:"continue"`
	SuppressOutput     bool           `json:"suppressOutput,omitempty"`
	SystemMessage      string         `json:"systemMessage,omitempty"`
	Decision           string         `json:"decision,omitempty"`
	Reason             string         `json:"reason,omitempty"`
	HookSpecificOutput map[string]any `json:"hookSpecificOutput,omitempty"`
}

func Handle(stdin io.Reader, explicitEvent string) (Output, int) {
	var input Input
	if err := json.NewDecoder(stdin).Decode(&input); err != nil && err != io.EOF {
		return Output{Continue: true, SuppressOutput: true}, 0
	}
	if input.HookEventName == "" {
		input.HookEventName = explicitEvent
	}
	root, err := project.FindRoot(input.CWD)
	if err != nil {
		return Output{Continue: true, SuppressOutput: true}, 0
	}
	cfg, installed, err := project.LoadOrDefault(root)
	if err != nil {
		return Output{Continue: true, SystemMessage: "HELMOR failed to load project config: " + err.Error()}, 0
	}
	state, _ := project.LoadState(root)

	switch input.HookEventName {
	case "SessionStart":
		msg := project.ContextCard(cfg)
		if !installed {
			msg += "\nHELMOR is observing from defaults. Run `helmor init` to install a project profile.\n"
		}
		_ = project.WriteContextCard(root, cfg)
		return Output{Continue: true, SystemMessage: msg}, 0
	case "UserPromptSubmit":
		return handlePrompt(root, cfg, state, input.UserPrompt), 0
	case "PreToolUse":
		return handlePreTool(root, cfg, state, input)
	case "PostToolUse":
		return handlePostTool(root, state, input), 0
	case "Stop":
		return handleStop(cfg, state)
	case "PreCompact", "SessionEnd":
		return handlePreserve(root, cfg, state, input.HookEventName), 0
	default:
		return Output{Continue: true, SuppressOutput: true}, 0
	}
}

func Encode(out Output, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(out)
}

func handlePrompt(root string, cfg config.Project, state config.RuntimeState, prompt string) Output {
	lower := strings.ToLower(prompt)
	state.LastStage = classifyStage(lower)
	if strings.Contains(lower, "approve helmor production action") {
		state.ProductionApproval = true
	}
	if strings.Contains(lower, "helmor security review complete") || strings.Contains(lower, "final security review complete") {
		state.SecurityReviewComplete = true
	}
	if state.LastStage == "launch" || state.LastStage == "security" || isSecuritySensitive(lower) {
		state.SecurityReviewRequired = true
	}
	_ = project.SaveState(root, state)

	var parts []string
	if state.LastStage != "" {
		parts = append(parts, "HELMOR stage: "+state.LastStage+".")
	}
	if state.SecurityReviewRequired && !state.SecurityReviewComplete {
		parts = append(parts, "Security review required before release/launch closeout.")
	}
	if cfg.Tools.CodeGraphFirst {
		parts = append(parts, "Use graph/code discovery before broad file search.")
	}
	if len(parts) == 0 {
		return Output{Continue: true, SuppressOutput: true}
	}
	return Output{Continue: true, SystemMessage: strings.Join(parts, " ")}
}

func handlePreTool(root string, cfg config.Project, state config.RuntimeState, input Input) (Output, int) {
	command := commandFrom(input.ToolInput)
	path := pathFrom(input.ToolInput)
	var reasons []string

	if input.ToolName == "mcp__codebase_memory_mcp" || strings.HasPrefix(input.ToolName, "mcp__codebase_memory_mcp") {
		state.CodeGraphSeen = true
		_ = project.SaveState(root, state)
	}
	if cfg.Policies.BlockSecretWrites && path != "" && detector.IsSecretPath(path) {
		reasons = append(reasons, "secret-shaped path: "+path)
	}
	if cfg.Policies.BlockDestructiveGit && detector.IsDestructiveGit(command) {
		reasons = append(reasons, "destructive git command")
	}
	if cfg.Tools.PackageRunnerGuard && detector.IsRunnerBypass(command, cfg.PackageRunner) {
		reasons = append(reasons, "package runner bypass; expected "+cfg.PackageRunner)
	}
	if cfg.Tools.CodeGraphFirst && !state.CodeGraphSeen && (input.ToolName == "Grep" || input.ToolName == "Glob" || detector.IsBroadSearch(command)) {
		reasons = append(reasons, "broad search before graph/code discovery")
	}
	if cfg.Policies.BlockUnsafeDeployActions && detector.IsUnsafeDeploy(command) && !state.ProductionApproval {
		reasons = append(reasons, "production/mainnet/provider action without approval")
	}

	if len(reasons) == 0 {
		return allow(), 0
	}
	message := "HELMOR " + cfg.Mode + " finding: " + strings.Join(reasons, "; ")
	if cfg.Mode == config.ModeGuard || cfg.Mode == config.ModeStrict {
		return deny(message), 2
	}
	return Output{Continue: true, SystemMessage: message, HookSpecificOutput: map[string]any{"permissionDecision": "allow"}}, 0
}

func handlePostTool(root string, state config.RuntimeState, input Input) Output {
	if path := pathFrom(input.ToolInput); path != "" {
		state.TouchedFiles = appendUnique(state.TouchedFiles, path)
		state.PendingChecks = appendUnique(state.PendingChecks, "test")
	}
	if strings.HasPrefix(input.ToolName, "mcp__codebase_memory_mcp") {
		state.CodeGraphSeen = true
	}
	_ = project.SaveState(root, state)
	return Output{Continue: true, SuppressOutput: true}
}

func handleStop(cfg config.Project, state config.RuntimeState) (Output, int) {
	var blockers []string
	if cfg.Mode == config.ModeStrict {
		if cfg.Policies.RequireHandoffOnStop && !state.HandoffSaved {
			blockers = append(blockers, "handoff not saved")
		}
		if cfg.Policies.RequireChecksOnStop && len(state.PendingChecks) > 0 {
			blockers = append(blockers, "pending checks: "+strings.Join(state.PendingChecks, ", "))
		}
		if state.SecurityReviewRequired && !state.SecurityReviewComplete {
			blockers = append(blockers, "security review incomplete")
		}
	}
	if len(blockers) > 0 {
		reason := strings.Join(blockers, "; ")
		return Output{Continue: false, Decision: "block", Reason: reason, SystemMessage: "HELMOR strict closeout blocked: " + reason}, 2
	}
	return Output{Continue: true, Decision: "approve", SystemMessage: "HELMOR closeout approved. Save memory or handoff if useful."}, 0
}

func handlePreserve(root string, cfg config.Project, state config.RuntimeState, event string) Output {
	path := filepath.Join(root, project.DirName, project.HandoffFile)
	summary := fmt.Sprintf("# HELMOR Handoff\n\n- Project: %s\n- Mode: %s\n- Event: %s\n- Last stage: %s\n- Pending checks: %s\n",
		cfg.ProjectName,
		cfg.Mode,
		event,
		state.LastStage,
		strings.Join(state.PendingChecks, ", "),
	)
	_ = project.SaveState(root, withHandoff(state, path))
	_ = writeFile(path, summary)
	return Output{Continue: true, SystemMessage: "HELMOR handoff preserved at " + path}
}

func allow() Output {
	return Output{Continue: true, SuppressOutput: true, HookSpecificOutput: map[string]any{"permissionDecision": "allow"}}
}

func deny(reason string) Output {
	return Output{Continue: false, SystemMessage: reason, HookSpecificOutput: map[string]any{"permissionDecision": "deny", "permissionDecisionReason": reason}}
}

func commandFrom(input map[string]any) string {
	if input == nil {
		return ""
	}
	for _, key := range []string{"command", "cmd"} {
		if value, ok := input[key].(string); ok {
			return value
		}
	}
	return ""
}

func pathFrom(input map[string]any) string {
	if input == nil {
		return ""
	}
	for _, key := range []string{"file_path", "path", "notebook_path"} {
		if value, ok := input[key].(string); ok {
			return value
		}
	}
	return ""
}

func classifyStage(prompt string) string {
	switch {
	case strings.Contains(prompt, "launch") || strings.Contains(prompt, "pilot") || strings.Contains(prompt, "go live"):
		return "launch"
	case strings.Contains(prompt, "security") || strings.Contains(prompt, "secret") || strings.Contains(prompt, "auth"):
		return "security"
	case strings.Contains(prompt, "test") || strings.Contains(prompt, "lint") || strings.Contains(prompt, "build"):
		return "verification"
	default:
		return ""
	}
}

func isSecuritySensitive(prompt string) bool {
	terms := []string{"production", "mainnet", "wallet", "signer", "kyc", "pii", "secret", "credential", "deploy"}
	for _, term := range terms {
		if strings.Contains(prompt, term) {
			return true
		}
	}
	return false
}

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func withHandoff(state config.RuntimeState, path string) config.RuntimeState {
	state.HandoffSaved = true
	state.LastCompactionSummaryPath = path
	return state
}

func writeFile(path string, data string) error {
	return projectWriteFile(path, []byte(data))
}

var projectWriteFile = func(path string, data []byte) error {
	return writeOSFile(path, data)
}
