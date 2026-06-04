package tools

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewProjectCopiesCollaborationFiles(t *testing.T) {
	root := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})
	if err := os.Chdir(root); err != nil {
		t.Fatalf("change directory: %v", err)
	}

	output := captureStdout(t, func() {
		if err := NewProject(ProjectNewOptions{
			ProjectName: "demoapp",
			ModulePath:  "example.com/demoapp",
		}); err != nil {
			t.Fatalf("new project: %v", err)
		}
	})

	for _, snippet := range []string{
		"Collaboration workspace initialized:",
		"AGENTS.md",
		"docs/handoff.md",
		".codex/agents/",
		"Start by reading AGENTS.md and docs/handoff.md.",
	} {
		if !strings.Contains(output, snippet) {
			t.Fatalf("stdout missing %q:\n%s", snippet, output)
		}
	}

	projectRoot := filepath.Join(root, "demoapp")
	wantFiles := []string{
		"AGENTS.md",
		".codex/agents/api-designer.md",
		".codex/agents/gin-architect.md",
		".codex/agents/gorm-expert.md",
		".codex/agents/quality-gate.md",
		"docs/handoff.md",
		"docs/decision-log.md",
		"docs/tasks.md",
		"docs/product/PRD.md",
		"docs/tech/ARCHITECTURE.md",
		"docs/review/quality-report.md",
		"docs/review/README.md",
		"docs/qa/test_cases.md",
		"docs/qa/README.md",
		".codex/agents/README.md",
	}
	for _, rel := range wantFiles {
		if _, err := os.Stat(filepath.Join(projectRoot, rel)); err != nil {
			t.Fatalf("missing %s: %v", rel, err)
		}
	}

	goMod, err := os.ReadFile(filepath.Join(projectRoot, "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if !strings.Contains(string(goMod), "module example.com/demoapp") {
		t.Fatalf("go.mod does not contain target module path:\n%s", string(goMod))
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	os.Stdout = writer
	t.Cleanup(func() {
		os.Stdout = original
	})

	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}
	os.Stdout = original

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(reader); err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close stdout reader: %v", err)
	}
	return buf.String()
}
