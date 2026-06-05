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
		"docs/ai-collaboration.md",
		".codex/agents/",
		"Tell the agent your goal",
		"The collaboration runner will maintain PRD, architecture, tasks, handoff, review, and QA notes.",
	} {
		if !strings.Contains(output, snippet) {
			t.Fatalf("stdout missing %q:\n%s", snippet, output)
		}
	}

	projectRoot := filepath.Join(root, "demoapp")
	wantFiles := []string{
		"AGENTS.md",
		".codex/agents/collaboration-runner.md",
		".codex/agents/pm.md",
		".codex/agents/architect.md",
		".codex/agents/developer.md",
		".codex/agents/api-designer.md",
		".codex/agents/gin-architect.md",
		".codex/agents/gorm-expert.md",
		".codex/agents/quality-gate.md",
		"docs/handoff.md",
		"docs/ai-collaboration.md",
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

	agents, err := os.ReadFile(filepath.Join(projectRoot, "AGENTS.md"))
	if err != nil {
		t.Fatalf("read AGENTS.md: %v", err)
	}
	for _, snippet := range []string{
		"Users should only describe the goal",
		"Automatic Collaboration Protocol",
		"collaboration-runner.md",
	} {
		if !strings.Contains(string(agents), snippet) {
			t.Fatalf("AGENTS.md missing %q:\n%s", snippet, string(agents))
		}
	}

	guide, err := os.ReadFile(filepath.Join(projectRoot, "docs/ai-collaboration.md"))
	if err != nil {
		t.Fatalf("read docs/ai-collaboration.md: %v", err)
	}
	for _, snippet := range []string{
		"AI Collaboration Guide / AI 协作使用指南",
		"Implement registration and login",
		"实现注册登录",
		"PM -> Architect -> Developer -> Reviewer",
	} {
		if !strings.Contains(string(guide), snippet) {
			t.Fatalf("docs/ai-collaboration.md missing %q:\n%s", snippet, string(guide))
		}
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
