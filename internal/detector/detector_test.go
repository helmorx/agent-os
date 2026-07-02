package detector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/helmorx/agent-os/internal/config"
)

func TestSecretPathDetector(t *testing.T) {
	cases := map[string]bool{
		".env":                                 true,
		".env.local":                           true,
		".env.example":                         false,
		".secrets/key.json":                    true,
		"id_rsa":                               true,
		"mainnet-treasury-backup.json":         true,
		"r2-creds.local.env":                   true,
		"signer-auth-token.local":              true,
		"services/wallet-monitor/package.json": false,
		"scripts/check-secret-shaped-filenames.sh": false,
		"tests/no-user-private-key-schema.test.ts": false,
		"src/index.ts": false,
	}
	for path, want := range cases {
		if got := IsSecretPath(path); got != want {
			t.Fatalf("IsSecretPath(%q) = %v, want %v", path, got, want)
		}
	}
}

func TestRunSkipsGitignoredLocalSecrets(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "README.md"), "# Test\n")
	write(t, filepath.Join(root, ".gitignore"), ".env.local\n.secrets/\n")
	write(t, filepath.Join(root, ".env.local"), "TOKEN=x\n")
	write(t, filepath.Join(root, ".secrets", "mainnet-treasury-backup.json"), "{}\n")

	cfg := config.Project{
		ProjectName:   "test",
		Mode:          config.ModeObserve,
		PackageRunner: "npm",
		TruthFiles:    []string{"README.md"},
	}
	findings := Run(root, cfg)
	if hasRule(findings, "secret-shaped-filename") {
		t.Fatalf("ignored local secret files should not produce doctor findings: %#v", findings)
	}
}

func TestCommandDetectors(t *testing.T) {
	if !IsDestructiveGit("git reset --hard HEAD") {
		t.Fatal("expected destructive git")
	}
	if !IsUnsafeDeploy("vercel deploy --prod") {
		t.Fatal("expected unsafe deploy")
	}
	if !IsRunnerBypass("npm install", "nub") {
		t.Fatal("expected runner bypass")
	}
	if IsRunnerBypass("nub install", "nub") {
		t.Fatal("did not expect runner bypass")
	}
}

func TestRunFindsDesignAndSecretFindings(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "README.md"), "# Test\n")
	write(t, filepath.Join(root, ".env.local"), "TOKEN=x\n")
	write(t, filepath.Join(root, "src", "app.css"), ".hero{background-clip:text;letter-spacing:-0.04em}\n")

	cfg := config.Project{
		ProjectName:   "test",
		Mode:          config.ModeObserve,
		PackageRunner: "npm",
		TruthFiles:    []string{"README.md"},
		Tools:         config.Tools{},
		Policies: config.Policies{
			DesignDetectors: true,
		},
	}
	findings := Run(root, cfg)
	if !hasRule(findings, "secret-shaped-filename") {
		t.Fatalf("missing secret finding: %#v", findings)
	}
	if !hasRule(findings, "design.gradient-text") {
		t.Fatalf("missing design finding: %#v", findings)
	}
}

func hasRule(findings []Finding, rule string) bool {
	for _, finding := range findings {
		if finding.Rule == rule {
			return true
		}
	}
	return false
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
