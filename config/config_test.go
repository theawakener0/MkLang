package config

import (
	"strings"
	"testing"
)

func TestDefault(t *testing.T) {
	spec := Default()
	if spec.Name != "Zod" {
		t.Errorf("default name = %q, want Zod", spec.Name)
	}
	if !spec.Features.Floats || !spec.Features.Arrays {
		t.Errorf("default features missing floats/arrays")
	}
	if err := spec.Validate(); err != nil {
		t.Errorf("default config should validate: %v", err)
	}
}

func TestParseMergesDefaults(t *testing.T) {
	data := []byte(`{"name": "Foo", "features": {"arrays": false}}`)
	spec, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if spec.Name != "Foo" {
		t.Errorf("name = %q, want Foo", spec.Name)
	}
	if spec.Features.Arrays {
		t.Error("arrays should be false (overridden)")
	}
	if !spec.Features.Floats {
		t.Error("floats should remain true (from defaults)")
	}
	if spec.Module != Default().Module {
		t.Errorf("module should fall back to default, got %q", spec.Module)
	}
}

func TestParseInvalidJSON(t *testing.T) {
	if _, err := Parse([]byte("{not json")); err == nil {
		t.Fatal("expected an error for invalid JSON")
	}
}

func TestValidateErrors(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Spec)
		want   string
	}{
		{
			name:   "empty name",
			mutate: func(s *Spec) { s.Name = " " },
			want:   "name must not be empty",
		},
		{
			name:   "bad module",
			mutate: func(s *Spec) { s.Module = "not a path!" },
			want:   "is not a valid Go module path",
		},
		{
			name:   "unknown keyword type",
			mutate: func(s *Spec) { s.Keywords["foo"] = "BOGUS" },
			want:   "unknown token type",
		},
		{
			name:   "unknown symbol type",
			mutate: func(s *Spec) { s.Symbols["~"] = "BOGUS" },
			want:   "unknown token type",
		},
		{
			name:   "unknown builtin",
			mutate: func(s *Spec) { s.Builtins["bogus"] = "bogus" },
			want:   "unknown builtin",
		},
		{
			name: "keyword and symbol share spelling",
			mutate: func(s *Spec) {
				s.Keywords["let"] = "LET"
				s.Symbols["let"] = "LET"
			},
			want: "both a keyword and a symbol",
		},
		{
			name: "missing required symbol",
			mutate: func(s *Spec) {
				delete(s.Symbols, "(")
			},
			want: `symbols must define "LPAREN"`,
		},
		{
			name: "strings without delimiter",
			mutate: func(s *Spec) { s.String.Delimiter = "" },
			want:   "requires a non-empty string.delimiter",
		},
		{
			name: "no declaration mechanism",
			mutate: func(s *Spec) {
				delete(s.Keywords, "let")
				delete(s.Symbols, ":=")
			},
			want: "no way to declare variables",
		},
		{
			name: "bad statement end token",
			mutate: func(s *Spec) { s.StatementEnd = append(s.StatementEnd, "BOGUS") },
			want:   "unknown token type",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			spec := Default()
			tc.mutate(spec)
			err := spec.Validate()
			if err == nil {
				t.Fatalf("expected validation error containing %q", tc.want)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error %q does not contain %q", err, tc.want)
			}
		})
	}
}

func TestNoAutoSemicolonsRequiresSemicolonSymbol(t *testing.T) {
	spec := Default()
	spec.Features.AutoSemicolons = false
	if err := spec.Validate(); err != nil {
		t.Fatalf("default has a ; symbol, should validate: %v", err)
	}

	delete(spec.Symbols, ";")
	err := spec.Validate()
	if err == nil || !strings.Contains(err.Error(), "auto_semicolons disabled") {
		t.Fatalf("expected auto_semicolons error, got %v", err)
	}
}

func TestWriteAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/lang.json"

	spec := Default()
	spec.Name = "Wibble"
	spec.Module = "github.com/example/wibble"
	if err := spec.WriteTo(path); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Name != "Wibble" {
		t.Errorf("loaded name = %q, want Wibble", loaded.Name)
	}
	if loaded.Features.Matrices != spec.Features.Matrices {
		t.Errorf("loaded features mismatch")
	}
}
