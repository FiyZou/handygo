package tools

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const scaffoldSourceModule = "github.com/FiyZou/handygo/examples"

//go:embed all:testdata/scaffold/**
var scaffoldFS embed.FS

type ProjectNewOptions struct {
	ProjectName        string
	ModulePath         string
	Force              bool
	Interactive        bool
	Frontend           *bool
	FrontendStyleSkill string
	Stdin              io.Reader
	Stdout             io.Writer
}

func NewProject(opts ProjectNewOptions) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	projectName := strings.TrimSpace(opts.ProjectName)
	if err := validateProjectName(projectName); err != nil {
		return err
	}
	modulePath := strings.TrimSpace(opts.ModulePath)
	if modulePath == "" {
		modulePath = projectName
	}

	targetDir, err := filepath.Abs(projectName)
	if err != nil {
		return err
	}
	if err := ensureTargetDir(targetDir, opts.Force); err != nil {
		return err
	}
	if err := copyScaffold(targetDir, opts.Force); err != nil {
		return err
	}
	if err := InitProject(ProjectInitOptions{
		Root:         targetDir,
		ModulePath:   modulePath,
		SourceModule: scaffoldSourceModule,
	}); err != nil {
		return err
	}
	frontendCfg, err := resolveFrontendConfig(opts, stdout)
	if err != nil {
		return err
	}
	if err := writeCollaborationConfig(targetDir, frontendCfg); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "Created HandyGo project in %s\n", targetDir)
	fmt.Fprintln(stdout, "Collaboration workspace initialized:")
	fmt.Fprintln(stdout, "  AGENTS.md")
	fmt.Fprintln(stdout, "  docs/handoff.md")
	fmt.Fprintln(stdout, "  docs/ai-collaboration.md")
	fmt.Fprintln(stdout, "  docs/collaboration-config.yaml")
	fmt.Fprintln(stdout, "  docs/tasks.md")
	fmt.Fprintln(stdout, "  docs/decision-log.md")
	fmt.Fprintln(stdout, "  .codex/agents/")
	printFrontendSummary(stdout, frontendCfg)
	fmt.Fprintln(stdout, "Tell the agent your goal, for example: \"实现注册登录\".")
	fmt.Fprintln(stdout, "Next steps:")
	fmt.Fprintf(stdout, "  cd %s\n", projectName)
	fmt.Fprintln(stdout, "  make install-tools")
	fmt.Fprintln(stdout, "  go mod tidy")
	fmt.Fprintln(stdout, "  make generate")
	fmt.Fprintln(stdout, "  make dev")
	fmt.Fprintln(stdout, "The collaboration runner will maintain PRD, architecture, tasks, handoff, review, and QA notes.")
	return nil
}

type frontendConfig struct {
	Enabled    bool
	StyleSkill string
}

func resolveFrontendConfig(opts ProjectNewOptions, stdout io.Writer) (frontendConfig, error) {
	cfg := frontendConfig{StyleSkill: strings.TrimSpace(opts.FrontendStyleSkill)}
	if opts.Frontend != nil {
		cfg.Enabled = *opts.Frontend
		return cfg, nil
	}
	if !opts.Interactive {
		return cfg, nil
	}
	stdin := opts.Stdin
	if stdin == nil {
		stdin = os.Stdin
	}
	reader := bufio.NewReader(stdin)
	fmt.Fprint(stdout, "Enable frontend agent workflow? / 是否启用前端 Agent 工作流？ [y/N]: ")
	answer, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return cfg, err
	}
	answer = strings.ToLower(strings.TrimSpace(answer))
	cfg.Enabled = answer == "y" || answer == "yes"
	if !cfg.Enabled {
		return cfg, nil
	}
	fmt.Fprint(stdout, "Frontend style skill name, leave empty to skip / 前端风格 skill 名称，可留空跳过: ")
	styleSkill, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return cfg, err
	}
	cfg.StyleSkill = strings.TrimSpace(styleSkill)
	return cfg, nil
}

func writeCollaborationConfig(targetDir string, cfg frontendConfig) error {
	content := fmt.Sprintf(`# Collaboration Config / 协作配置
frontend:
  # Enable frontend agent workflow.
  # 是否启用前端 Agent 工作流。
  enabled: %t

  # Default frontend style skill. Leave empty to use existing project style.
  # 默认前端风格 skill。留空时使用项目已有风格。
  styleSkill: %q
`, cfg.Enabled, cfg.StyleSkill)
	return os.WriteFile(filepath.Join(targetDir, "docs", "collaboration-config.yaml"), []byte(content), 0644)
}

func printFrontendSummary(stdout io.Writer, cfg frontendConfig) {
	if cfg.Enabled {
		fmt.Fprintln(stdout, "Frontend workflow is enabled. / 前端工作流已开启。")
		if cfg.StyleSkill != "" {
			fmt.Fprintf(stdout, "Frontend style skill: %s\n", cfg.StyleSkill)
		}
		return
	}
	fmt.Fprintln(stdout, "Frontend workflow is disabled. You can enable it later in docs/collaboration-config.yaml.")
	fmt.Fprintln(stdout, "前端工作流已关闭。你可以稍后在 docs/collaboration-config.yaml 中开启。")
}

func validateProjectName(projectName string) error {
	switch projectName {
	case "", ".", "..", string(os.PathSeparator):
		return fmt.Errorf("invalid project name: %q", projectName)
	}
	if filepath.IsAbs(projectName) {
		return fmt.Errorf("project name must be a relative directory: %s", projectName)
	}
	if strings.Contains(projectName, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("project name cannot traverse parent directories: %s", projectName)
	}
	return nil
}

func ensureTargetDir(targetDir string, force bool) error {
	info, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		return os.MkdirAll(targetDir, 0755)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("target exists and is not a directory: %s", targetDir)
	}
	if !force {
		return fmt.Errorf("target directory already exists: %s", targetDir)
	}
	return nil
}

func copyScaffold(targetDir string, force bool) error {
	const scaffoldRoot = "testdata/scaffold"
	return fs.WalkDir(scaffoldFS, scaffoldRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == scaffoldRoot {
			return nil
		}
		rel, err := filepath.Rel(scaffoldRoot, path)
		if err != nil {
			return err
		}
		if shouldSkipScaffoldPath(rel, d) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		targetPath := filepath.Join(targetDir, rel)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}
		if !force {
			if _, err := os.Stat(targetPath); err == nil {
				return fmt.Errorf("target file already exists: %s", targetPath)
			} else if !os.IsNotExist(err) {
				return err
			}
		}
		return copyScaffoldFile(path, targetPath)
	})
}

func shouldSkipScaffoldPath(rel string, d fs.DirEntry) bool {
	name := d.Name()
	if d.IsDir() {
		switch name {
		case ".git", ".cache", "vendor":
			return true
		}
		return false
	}
	switch {
	case name == ".genmanifest.json":
		return true
	case strings.HasSuffix(name, ".db"):
		return true
	case strings.HasSuffix(name, ".db-shm"):
		return true
	case strings.HasSuffix(name, ".db-wal"):
		return true
	case name == "handygo-example":
		return true
	default:
		return false
	}
}

func copyScaffoldFile(sourcePath string, targetPath string) error {
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}
	source, err := scaffoldFS.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := os.OpenFile(targetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = io.Copy(target, source)
	return err
}
