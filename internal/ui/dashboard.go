package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/helmorx/agent-os/internal/agent"
	"github.com/helmorx/agent-os/internal/config"
	"github.com/helmorx/agent-os/internal/detector"
)

func Dashboard(writer io.Writer, cfg config.Project, findings []detector.Finding) {
	fmt.Fprintln(writer, "HELMOR Agent OS")
	fmt.Fprintln(writer, strings.Repeat("=", 16))
	fmt.Fprintf(writer, "Project: %s\n", cfg.ProjectName)
	fmt.Fprintf(writer, "Mode: %s\n", cfg.Mode)
	fmt.Fprintf(writer, "Framework: %s\n", cfg.Framework)
	fmt.Fprintf(writer, "Package runner: %s\n", cfg.PackageRunner)
	fmt.Fprintln(writer)

	fmt.Fprintln(writer, "Agents")
	for _, line := range agent.AdapterSummary(cfg) {
		fmt.Fprintf(writer, "- %s\n", line)
	}
	fmt.Fprintln(writer)

	fmt.Fprintln(writer, "Checks")
	if len(cfg.Checks) == 0 {
		fmt.Fprintln(writer, "- none detected")
	} else {
		for _, check := range cfg.Checks {
			fmt.Fprintf(writer, "- %s: %s\n", check.Name, check.Command)
		}
	}
	fmt.Fprintln(writer)

	fmt.Fprintln(writer, "Findings")
	if len(findings) == 0 {
		fmt.Fprintln(writer, "- none")
	} else {
		for _, finding := range findings {
			path := ""
			if finding.Path != "" {
				path = " " + finding.Path
			}
			fmt.Fprintf(writer, "- [%s] %s%s: %s\n", finding.Severity, finding.Rule, path, finding.Message)
		}
	}
}
