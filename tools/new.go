package tools

import (
	"embed"
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
	ProjectName string
	ModulePath  string
	Force       bool
}

func NewProject(opts ProjectNewOptions) error {
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

	fmt.Printf("Created HandyGo project in %s\n", targetDir)
	fmt.Println("Collaboration workspace initialized:")
	fmt.Println("  AGENTS.md")
	fmt.Println("  docs/handoff.md")
	fmt.Println("  docs/ai-collaboration.md")
	fmt.Println("  docs/tasks.md")
	fmt.Println("  docs/decision-log.md")
	fmt.Println("  .codex/agents/")
	fmt.Println("Tell the agent your goal, for example: \"实现注册登录\".")
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  make install-tools")
	fmt.Println("  go mod tidy")
	fmt.Println("  make generate")
	fmt.Println("  make dev")
	fmt.Println("The collaboration runner will maintain PRD, architecture, tasks, handoff, review, and QA notes.")
	return nil
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
