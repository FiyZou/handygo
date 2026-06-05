package tools

import (
	"bytes"
	"io"
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

	var tidyCalls []string
	output := captureStdout(t, func() {
		if err := NewProject(ProjectNewOptions{
			ProjectName: "demoapp",
			ModulePath:  "example.com/demoapp",
			CommandRunner: func(dir string, name string, args ...string) error {
				tidyCalls = append(tidyCalls, strings.Join(append([]string{name}, args...), " "))
				if !strings.HasSuffix(dir, "demoapp") {
					t.Fatalf("go mod tidy dir = %q", dir)
				}
				return nil
			},
		}); err != nil {
			t.Fatalf("new project: %v", err)
		}
	})

	for _, snippet := range []string{
		"Collaboration workspace initialized:",
		"AGENTS.md",
		"docs/handoff.md",
		"docs/ai-collaboration.md",
		"docs/collaboration-config.yaml",
		".codex/agents/",
		"Running go mod tidy...",
		"Tell the agent your goal",
		"The collaboration runner will maintain PRD, architecture, tasks, handoff, review, and QA notes.",
	} {
		if !strings.Contains(output, snippet) {
			t.Fatalf("stdout missing %q:\n%s", snippet, output)
		}
	}
	if len(tidyCalls) != 1 || tidyCalls[0] != "go mod tidy" {
		t.Fatalf("go mod tidy calls = %#v", tidyCalls)
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
		"docs/collaboration-config.yaml",
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
		"docs/collaboration-config.yaml",
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
		"Frontend workflow",
		"前端工作流",
	} {
		if !strings.Contains(string(guide), snippet) {
			t.Fatalf("docs/ai-collaboration.md missing %q:\n%s", snippet, string(guide))
		}
	}

	config, err := os.ReadFile(filepath.Join(projectRoot, "manifest/config.yaml"))
	if err != nil {
		t.Fatalf("read manifest/config.yaml: %v", err)
	}
	for _, snippet := range []string{
		"cache:",
		"enabled: false",
		"timeFormat: \"2006-01-02 15:04:05\"",
		"fileOutputPath: logs/app.log",
		"errorFileOutputPath: logs/error.log",
		"rotation:",
		"maxSizeMB: 100",
		"maxAgeDays: 30",
		"# MySQL example:",
		"mysql://user:password@127.0.0.1:3306/handygo",
		"# PostgreSQL example:",
		"postgres://postgres:postgres@127.0.0.1:5432/handygo",
		"driver: sqlite",
	} {
		if !strings.Contains(string(config), snippet) {
			t.Fatalf("manifest/config.yaml missing %q:\n%s", snippet, string(config))
		}
	}
	for _, snippet := range []string{
		"client:\n    name: asynq-client\n    redis:",
		"server:\n    name: asynq-server\n    concurrency: 10\n    queues:\n      critical: 6\n      default: 3\n      low: 1\n    redis:",
		"scheduler:\n    name: asynq-scheduler\n    location: Asia/Shanghai\n    redis:",
	} {
		if strings.Contains(string(config), snippet) {
			t.Fatalf("manifest/config.yaml should not contain separate asynq redis block %q:\n%s", snippet, string(config))
		}
	}
	if strings.Contains(string(config), "host=127.0.0.1 user=postgres") {
		t.Fatalf("manifest/config.yaml should use PostgreSQL URL DSN, not key/value DSN:\n%s", string(config))
	}

	collaborationConfig, err := os.ReadFile(filepath.Join(projectRoot, "docs/collaboration-config.yaml"))
	if err != nil {
		t.Fatalf("read docs/collaboration-config.yaml: %v", err)
	}
	for _, snippet := range []string{
		"# Collaboration Config / 协作配置",
		"enabled: false",
		"styleSkill: \"\"",
		"是否启用前端 Agent 工作流",
	} {
		if !strings.Contains(string(collaborationConfig), snippet) {
			t.Fatalf("docs/collaboration-config.yaml missing %q:\n%s", snippet, string(collaborationConfig))
		}
	}
}

func TestNewProjectWritesFrontendCollaborationConfig(t *testing.T) {
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

	enabled := true
	var stdout bytes.Buffer
	if err := NewProject(ProjectNewOptions{
		ProjectName:        "frontendapp",
		ModulePath:         "example.com/frontendapp",
		Frontend:           &enabled,
		FrontendStyleSkill: "my-ui-style",
		Stdout:             &stdout,
		SkipTidy:           true,
	}); err != nil {
		t.Fatalf("new project: %v", err)
	}

	collaborationConfig, err := os.ReadFile(filepath.Join(root, "frontendapp", "docs/collaboration-config.yaml"))
	if err != nil {
		t.Fatalf("read docs/collaboration-config.yaml: %v", err)
	}
	for _, snippet := range []string{
		"enabled: true",
		"styleSkill: \"my-ui-style\"",
	} {
		if !strings.Contains(string(collaborationConfig), snippet) {
			t.Fatalf("docs/collaboration-config.yaml missing %q:\n%s", snippet, string(collaborationConfig))
		}
	}
	if !strings.Contains(stdout.String(), "Frontend workflow is enabled") {
		t.Fatalf("stdout missing frontend enabled message:\n%s", stdout.String())
	}
}

func TestNewProjectPromptsForFrontendWorkflow(t *testing.T) {
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

	var stdout bytes.Buffer
	if err := NewProject(ProjectNewOptions{
		ProjectName: "promptapp",
		ModulePath:  "example.com/promptapp",
		Interactive: true,
		Stdin:       strings.NewReader("y\nprompt-style\n"),
		Stdout:      &stdout,
		SkipTidy:    true,
	}); err != nil {
		t.Fatalf("new project: %v", err)
	}

	collaborationConfig, err := os.ReadFile(filepath.Join(root, "promptapp", "docs/collaboration-config.yaml"))
	if err != nil {
		t.Fatalf("read docs/collaboration-config.yaml: %v", err)
	}
	if !strings.Contains(string(collaborationConfig), "enabled: true") || !strings.Contains(string(collaborationConfig), "styleSkill: \"prompt-style\"") {
		t.Fatalf("unexpected docs/collaboration-config.yaml:\n%s", string(collaborationConfig))
	}
	for _, snippet := range []string{
		"Enable frontend agent workflow? / 是否启用前端 Agent 工作流？",
		"Frontend style skill name, leave empty to skip / 前端风格 skill 名称，可留空跳过",
	} {
		if !strings.Contains(stdout.String(), snippet) {
			t.Fatalf("stdout missing %q:\n%s", snippet, stdout.String())
		}
	}
}

func TestNewProjectCanSkipGoModTidy(t *testing.T) {
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

	called := false
	if err := NewProject(ProjectNewOptions{
		ProjectName: "skipapp",
		ModulePath:  "example.com/skipapp",
		SkipTidy:    true,
		CommandRunner: func(dir string, name string, args ...string) error {
			called = true
			return nil
		},
		Stdout: io.Discard,
	}); err != nil {
		t.Fatalf("new project: %v", err)
	}
	if called {
		t.Fatal("expected go mod tidy to be skipped")
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
