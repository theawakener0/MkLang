// Package generator writes a complete, self-contained language project from a
// config.Spec. The engine packages (lexer, parser, evaluator, ...) are generic
// and copied as-is; the config-specific parts (keyword/symbol tables, the spec
// defaults, main.go, go.mod, README) are rendered from the spec.
package generator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/theawakener0/MkLang/config"
)

// engineModule is the module path baked into the embedded engine files. It is
// rewritten to spec.Module when copying.
const engineModule = "github.com/theawakener0/MkLang"

// Generate writes a complete language project described by spec into outDir.
// fsys should contain the embedded engine packages. If outDir exists and is
// non-empty, an error is returned unless force is set.
func Generate(spec *config.Spec, fsys fs.FS, outDir string, force bool) error {
	if err := checkOutDir(outDir, force); err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	rendered := map[string][]byte{
		"token/token.go": nil,
		"token/spec.go":  nil,
		"main.go":        nil,
		"go.mod":         renderGoMod(spec),
		"README.md":      renderReadme(spec),
		"lang.json":      renderLangJSON(spec),
	}
	var err error
	if rendered["token/token.go"], err = renderTokenGo(spec); err != nil {
		return fmt.Errorf("render token/token.go: %w", err)
	}
	if rendered["token/spec.go"], err = renderSpecGo(spec); err != nil {
		return fmt.Errorf("render token/spec.go: %w", err)
	}
	if rendered["main.go"], err = renderMainGo(spec); err != nil {
		return fmt.Errorf("render main.go: %w", err)
	}

	// Copy the generic engine files, skipping generated ones and test files,
	// rewriting the module path to the user's module.
	if err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if _, generated := rendered[path]; generated {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}
		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		content = []byte(strings.ReplaceAll(string(content), engineModule, spec.Module))
		return writeProjectFile(outDir, path, content)
	}); err != nil {
		return fmt.Errorf("copy engine files: %w", err)
	}

	for path, content := range rendered {
		if err := writeProjectFile(outDir, path, content); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
	}
	return nil
}

// checkOutDir rejects a non-empty output directory unless force is set.
func checkOutDir(outDir string, force bool) error {
	entries, err := os.ReadDir(outDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(entries) > 0 && !force {
		return fmt.Errorf("output directory %q is not empty (use -force to overwrite)", outDir)
	}
	return nil
}

// writeProjectFile writes content to outDir/path, creating parent directories.
func writeProjectFile(outDir, path string, content []byte) error {
	out := filepath.Join(outDir, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	return os.WriteFile(out, content, 0o644)
}
