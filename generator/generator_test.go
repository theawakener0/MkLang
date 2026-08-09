package generator

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/theawakener0/MkLang/config"
)

// engineDirs mirrors the go:embed patterns in main.go. Tests reuse the real
// engine by reading it from the repository root.
var engineDirs = []string{"token", "lexer", "parser", "ast", "evaluator", "object", "repl"}

// engineDirFS exposes just the engine packages (plus LICENSE/.gitignore) from
// the repo root, mimicking the embedded FS the CLI uses.
type engineDirFS struct{ inner fs.FS }

func (e engineDirFS) Open(name string) (fs.File, error) {
	if name == "." {
		return e.inner.Open(name)
	}
	if !isEnginePath(name) {
		return nil, fs.ErrNotExist
	}
	return e.inner.Open(name)
}

func (e engineDirFS) ReadDir(name string) ([]fs.DirEntry, error) {
	entries, err := fs.ReadDir(e.inner, name)
	if err != nil {
		return nil, err
	}
	out := entries[:0]
	for _, en := range entries {
		full := name + "/" + en.Name()
		if isEnginePath(full) {
			out = append(out, en)
		}
	}
	return out, nil
}

func (e engineDirFS) ReadFile(name string) ([]byte, error) {
	if !isEnginePath(name) {
		return nil, fs.ErrNotExist
	}
	return fs.ReadFile(e.inner, name)
}

func isEnginePath(name string) bool {
	name = strings.TrimPrefix(name, "./")
	for _, dir := range engineDirs {
		if name == dir || strings.HasPrefix(name, dir+"/") {
			return true
		}
	}
	return name == "LICENSE" || name == ".gitignore"
}

func loadSpec(t *testing.T, path string) *config.Spec {
	t.Helper()
	spec, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := spec.Validate(); err != nil {
		t.Fatal(err)
	}
	return spec
}

func TestGenerateWritesExpectedFiles(t *testing.T) {
	spec := loadSpec(t, "../examples/configs/pythonish.json")
	out := t.TempDir()

	if err := Generate(spec, engineDirFS{os.DirFS("..")}, out, false); err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"token/token.go", "token/spec.go", "token/lookup.go",
		"lexer/lexer.go", "parser/parser.go", "ast/ast.go",
		"evaluator/evaluator.go", "evaluator/builtins.go", "object/object.go",
		"repl/repl.go", "main.go", "go.mod", "README.md", "lang.json",
		".gitignore", "LICENSE",
	}
	for _, f := range expected {
		if _, err := os.Stat(filepath.Join(out, f)); err != nil {
			t.Errorf("expected generated file %s: %v", f, err)
		}
	}

	// Test files must not leak into generated projects.
	if _, err := os.Stat(filepath.Join(out, "lexer/lexer_test.go")); !os.IsNotExist(err) {
		t.Error("generated project should not contain _test.go files")
	}
}

func TestGenerateUsesConfigSpellings(t *testing.T) {
	spec := loadSpec(t, "../examples/configs/pythonish.json")
	out := t.TempDir()
	if err := Generate(spec, engineDirFS{os.DirFS("..")}, out, false); err != nil {
		t.Fatal(err)
	}

	tokenGo, err := os.ReadFile(filepath.Join(out, "token/token.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"def":\s+FUNCTION`, `"not":\s+BANG`, `"and":\s+LAND`} {
		if !regexp.MustCompile(want).Match(tokenGo) {
			t.Errorf("token/token.go missing %q", want)
		}
	}
	for _, banned := range []string{`"fn":`, `"let":`, `"&&":`} {
		if strings.Contains(string(tokenGo), banned) {
			t.Errorf("token/token.go should not contain %q for pythonish", banned)
		}
	}

	specGo, err := os.ReadFile(filepath.Join(out, "token/spec.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`Name:    "Pythonish"`, `Delimiter: "\""`} {
		if !strings.Contains(string(specGo), want) {
			t.Errorf("token/spec.go missing %q", want)
		}
	}

	mainGo, err := os.ReadFile(filepath.Join(out, "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(mainGo), spec.Module) {
		t.Error("main.go should import the configured module")
	}
}

func TestGenerateRewritesModulePath(t *testing.T) {
	spec := loadSpec(t, "../examples/configs/pythonish.json")
	out := t.TempDir()
	if err := Generate(spec, engineDirFS{os.DirFS("..")}, out, false); err != nil {
		t.Fatal(err)
	}

	// No generated file may still reference the tool's module path.
	err := filepath.WalkDir(out, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(content), "github.com/theawakener0/MkLang") {
			t.Errorf("%s still references the engine module", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateRejectsNonEmptyDir(t *testing.T) {
	spec := loadSpec(t, "../examples/configs/pythonish.json")
	out := t.TempDir()
	if err := os.WriteFile(filepath.Join(out, "existing.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Generate(spec, engineDirFS{os.DirFS("..")}, out, false); err == nil {
		t.Error("expected an error for a non-empty output directory")
	}
	if err := Generate(spec, engineDirFS{os.DirFS("..")}, out, true); err != nil {
		t.Errorf("force should allow overwriting: %v", err)
	}
}

func TestGeneratedProjectBuildsAndRuns(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go toolchain not available")
	}

	cases := []struct {
		configPath string
		script     string
		want       []string
	}{
		{
			configPath: "../examples/configs/pythonish.json",
			script: `print("hello from pythonish")
x = 10
if (x > 5 and x < 20) {
    print("big-ish")
}
arr = [1, 2, 3]
print(len(arr))
`,
			want: []string{"hello from pythonish", "big-ish", "3"},
		},
	}

	for _, tc := range cases {
		t.Run(filepath.Base(tc.configPath), func(t *testing.T) {
			spec := loadSpec(t, tc.configPath)
			out := t.TempDir()
			if err := Generate(spec, engineDirFS{os.DirFS("..")}, out, false); err != nil {
				t.Fatal(err)
			}

			bin := filepath.Join(out, strings.ToLower(spec.Name))
			cmd := exec.Command("go", "build", "-o", bin, ".")
			cmd.Dir = out
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("go build failed: %v\n%s", err, output)
			}

			script := filepath.Join(out, "script.txt")
			if err := os.WriteFile(script, []byte(tc.script), 0o644); err != nil {
				t.Fatal(err)
			}
			cmd = exec.Command(bin, "script.txt")
			cmd.Dir = out
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("running generated language failed: %v\n%s", err, output)
			}
			for _, want := range tc.want {
				if !strings.Contains(string(output), want) {
					t.Errorf("output %q missing %q", output, want)
				}
			}
		})
	}
}
