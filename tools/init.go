package tools

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type ProjectInitOptions struct {
	ModulePath   string
	SourceModule string
	Root         string
}

func InitProject(opts ProjectInitOptions) error {
	if strings.TrimSpace(opts.ModulePath) == "" {
		return fmt.Errorf("module path cannot be empty")
	}
	if strings.TrimSpace(opts.SourceModule) == "" {
		return fmt.Errorf("source module cannot be empty")
	}

	root := opts.Root
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	if err := writeGoMod(root, opts.ModulePath); err != nil {
		return err
	}
	return replaceGoImports(root, opts.SourceModule, opts.ModulePath)
}

func writeGoMod(root string, modulePath string) error {
	path := filepath.Join(root, "go.mod")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		content := fmt.Sprintf("module %s\n\ngo 1.26.3\n", modulePath)
		return os.WriteFile(path, []byte(content), 0644)
	}
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "module ") {
		return fmt.Errorf("go.mod has no module declaration")
	}
	lines[0] = "module " + modulePath
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}

func replaceGoImports(root string, sourceModule string, modulePath string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "vendor", "node_modules":
				return filepath.SkipDir
			default:
				return nil
			}
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		updated := strings.ReplaceAll(string(data), sourceModule, modulePath)
		if updated == string(data) {
			return nil
		}
		return os.WriteFile(path, []byte(updated), 0644)
	})
}
