// Package config defines the JSON schema used to describe a language that
// MkLang can generate, along with loading and validation helpers.
package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	tk "github.com/theawakener0/MkLang/token"
)

// Features toggles which parts of the language are available.
type Features struct {
	Floats              bool `json:"floats"`
	Strings             bool `json:"strings"`
	Arrays              bool `json:"arrays"`
	Hashes              bool `json:"hashes"`
	ForLoops            bool `json:"for_loops"`
	InfiniteLoops       bool `json:"infinite_loops"`
	Try                 bool `json:"try"`
	PrefixPostfixIncDec bool `json:"prefix_postfix_inc_dec"`
	CompoundAssignment  bool `json:"compound_assignment"`
	ElseIf              bool `json:"elseif"`
	Matrices            bool `json:"matrices"`
	AutoSemicolons      bool `json:"auto_semicolons"`
	Comments            bool `json:"comments"`
	BlockComments       bool `json:"block_comments"`
}

// Comments describes the comment markers of the language.
type Comments struct {
	Line       string `json:"line"`
	BlockStart string `json:"block_start"`
	BlockEnd   string `json:"block_end"`
}

// Identifier describes which characters may form identifiers.
type Identifier struct {
	StartLetters      bool   `json:"start_letters"`
	ExtraStart        string `json:"extra_start"`
	PartLettersDigits bool   `json:"part_letters_digits"`
	ExtraPart         string `json:"extra_part"`
}

// StringConfig describes how string literals are written.
type StringConfig struct {
	Delimiter       string `json:"delimiter"`
	EscapeSequences bool   `json:"escape_sequences"`
}

// Output describes how values are printed.
type Output struct {
	True       string `json:"true"`
	False      string `json:"false"`
	Null       string `json:"null"`
	Break      string `json:"break"`
	Continue   string `json:"continue"`
	Function   string `json:"function"`
	ArrayOpen  string `json:"array_open"`
	ArrayClose string `json:"array_close"`
	ArraySep   string `json:"array_separator"`
	HashOpen   string `json:"hash_open"`
	HashClose  string `json:"hash_close"`
	HashSep    string `json:"hash_separator"`
	HashKeyVal string `json:"hash_keyval"`
}

// Spec is the JSON schema for a language. Values left out of the JSON fall
// back to the built-in defaults (see Default).
type Spec struct {
	Name         string            `json:"name"`
	Module       string            `json:"module"`
	Version      string            `json:"version"`
	Features     Features          `json:"features"`
	Comments     Comments          `json:"comments"`
	Identifier   Identifier        `json:"identifier"`
	String       StringConfig      `json:"string"`
	Keywords     map[string]string `json:"keywords"`       // spelling -> canonical token type
	Symbols      map[string]string `json:"symbols"`        // spelling -> canonical token type
	Builtins     map[string]string `json:"builtins"`       // spelling -> canonical builtin
	StatementEnd []string          `json:"statement_end_tokens,omitempty"`
	Output       Output            `json:"output"`
}

// Default returns the built-in Zod language as a config template.
func Default() *Spec {
	return &Spec{
		Name:    "Zod",
		Module:  "github.com/yourname/yourlang",
		Version: "0.1.0",
		Features: Features{
			Floats:              true,
			Strings:             true,
			Arrays:              true,
			Hashes:              true,
			ForLoops:            true,
			InfiniteLoops:       true,
			Try:                 true,
			PrefixPostfixIncDec: true,
			CompoundAssignment:  true,
			ElseIf:              true,
			Matrices:            true,
			AutoSemicolons:      true,
			Comments:            true,
			BlockComments:       true,
		},
		Comments:   Comments{Line: "//", BlockStart: "/*", BlockEnd: "*/"},
		Identifier: Identifier{StartLetters: true, ExtraStart: "_", PartLettersDigits: true, ExtraPart: "_"},
		String:     StringConfig{Delimiter: "\"", EscapeSequences: true},
		Keywords: map[string]string{
			"fn": "FUNCTION", "let": "LET", "true": "TRUE", "false": "FALSE",
			"if": "IF", "elseif": "ELSEIF", "else": "ELSE", "return": "RETURN",
			"for": "FOR", "loop": "LOOP", "break": "BREAK", "continue": "CONTINUE", "null": "NULL",
		},
		Symbols: map[string]string{
			"=": "ASSIGN", "==": "EQ", "!": "BANG", "!=": "NOTEQ",
			"+": "PLUS", "++": "INC", "+=": "INCASSIGN",
			"-": "MINUS", "--": "DEC", "-=": "DECDASSIGN",
			"*": "ASTERISK", "*=": "MLTASSIGN",
			"/": "SLASH", "/=": "DIVASSIGN",
			"<": "LT", "<=": "LTEQ", ">": "GT", ">=": "GTEQ",
			"&&": "LAND", "||": "LOR", ":=": "ASSIGNCHAR",
			",": "COMMA", ":": "COLOMN", ";": "SEMICOLON", ".": "DOT",
			"(": "LPAREN", ")": "RPAREN", "{": "LBRACE", "}": "RBRACE",
			"[": "LBRACKET", "]": "RBRACKET",
		},
		Builtins: map[string]string{
			"len": "len", "println": "println", "printf": "printf", "input": "input",
			"int": "int", "float": "float", "string": "string", "type": "type",
			"first": "first", "last": "last", "pop": "pop", "push": "push",
			"insert": "insert", "remove": "remove", "keys": "keys", "vals": "vals",
			"contains": "contains", "random": "random", "matrix": "matrix",
			"make": "make", "color": "color", "sleep": "sleep", "exp": "exp", "pi": "pi",
			"try": "try",
		},
		StatementEnd: []string{
			"IDENT", "INT", "FLOAT", "STRING", "TRUE", "FALSE", "NULL",
			"RPAREN", "RBRACKET", "RBRACE", "INC", "DEC",
			"RETURN", "BREAK", "CONTINUE",
		},
		Output: Output{
			True:       "true",
			False:      "false",
			Null:       "null",
			Break:      "break",
			Continue:   "continue",
			Function:   "fn",
			ArrayOpen:  "[",
			ArrayClose: "]",
			ArraySep:   ", ",
			HashOpen:   "{",
			HashClose:  "}",
			HashSep:    ", ",
			HashKeyVal: ": ",
		},
	}
}

// Load reads and parses a config file.
func Load(path string) (*Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

// Parse parses JSON config data. Any top-level field that is absent falls back
// to the built-in defaults; any field that is present replaces them. In
// particular, providing "keywords", "symbols", or "builtins" replaces the
// default table entirely, so deleting a line from an edited template removes
// that token from the generated language.
func Parse(data []byte) (*Spec, error) {
	spec := Default()
	if err := json.Unmarshal(data, spec); err != nil {
		return nil, fmt.Errorf("invalid config JSON: %w", err)
	}
	return spec, nil
}

// UnmarshalJSON implements the field-level fallback rules of Parse. The three
// map fields are shadowed with json.RawMessage so that json.Unmarshal's
// default map-merging behavior (which would resurrect deleted entries) is
// replaced with a wholesale swap whenever the key is present.
func (s *Spec) UnmarshalJSON(data []byte) error {
	type shadow Spec
	var raw struct {
		*shadow
		Keywords json.RawMessage `json:"keywords"`
		Symbols  json.RawMessage `json:"symbols"`
		Builtins json.RawMessage `json:"builtins"`
	}
	raw.shadow = (*shadow)(Default())
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Spec(*raw.shadow)

	for _, f := range []struct {
		name string
		raw  json.RawMessage
		dst  *map[string]string
	}{
		{"keywords", raw.Keywords, &s.Keywords},
		{"symbols", raw.Symbols, &s.Symbols},
		{"builtins", raw.Builtins, &s.Builtins},
	} {
		if f.raw == nil || bytes.Equal(f.raw, []byte("null")) {
			continue
		}
		var m map[string]string
		if err := json.Unmarshal(f.raw, &m); err != nil {
			return fmt.Errorf("invalid %s: %w", f.name, err)
		}
		*f.dst = m
	}
	return nil
}

// WriteTo writes the config as indented JSON to a file.
func (s *Spec) WriteTo(path string) error {
	data, err := s.MarshalJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// MarshalJSON returns the indented JSON form of the config. The plain type
// alias prevents infinite recursion through json.Marshaler.
func (s *Spec) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent((*plainSpec)(s), "", "  ")
}

// plainSpec is Spec without its MarshalJSON method.
type plainSpec Spec

var (
	modulePathRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._/-]*$`)
	validBuiltin = map[string]bool{
		"len": true, "println": true, "printf": true, "input": true,
		"int": true, "float": true, "string": true, "type": true,
		"first": true, "last": true, "pop": true, "push": true,
		"insert": true, "remove": true, "keys": true, "vals": true,
		"contains": true, "random": true, "matrix": true, "make": true,
		"color": true, "sleep": true, "exp": true, "pi": true,
		"try": true,
	}
)

// Validate checks the config for common mistakes and returns a descriptive
// error if any are found.
func (s *Spec) Validate() error {
	var errs []string

	if strings.TrimSpace(s.Name) == "" {
		errs = append(errs, "name must not be empty")
	}
	if strings.TrimSpace(s.Module) == "" {
		errs = append(errs, "module must not be empty")
	} else if !modulePathRe.MatchString(s.Module) {
		errs = append(errs, fmt.Sprintf("module %q is not a valid Go module path", s.Module))
	}

	for spelling, typeName := range s.Keywords {
		if !tk.IsTypeName(typeName) {
			errs = append(errs, fmt.Sprintf("keyword %q maps to unknown token type %q", spelling, typeName))
		}
	}
	for spelling, typeName := range s.Symbols {
		if !tk.IsTypeName(typeName) {
			errs = append(errs, fmt.Sprintf("symbol %q maps to unknown token type %q", spelling, typeName))
		}
	}
	for spelling := range s.Keywords {
		if _, ok := s.Symbols[spelling]; ok {
			errs = append(errs, fmt.Sprintf("spelling %q is used by both a keyword and a symbol", spelling))
		}
	}
	for spelling, canonical := range s.Builtins {
		if !validBuiltin[canonical] {
			errs = append(errs, fmt.Sprintf("builtin %q maps to unknown builtin %q", spelling, canonical))
		}
	}

	for _, required := range []string{"LPAREN", "RPAREN", "LBRACE", "RBRACE", "COMMA"} {
		if !hasSymbol(s.Symbols, required) {
			errs = append(errs, fmt.Sprintf("symbols must define %q", required))
		}
	}

	// A language must be able to declare a variable: either a let keyword or
	// the := assignchar symbol (bare x = 5 only reassigns existing variables).
	if !hasKeyword(s.Keywords, "LET") && !hasSymbol(s.Symbols, "ASSIGNCHAR") {
		errs = append(errs, "no way to declare variables: add a LET keyword or a \":=\" (ASSIGNCHAR) symbol")
	}

	if s.Features.Strings && s.String.Delimiter == "" {
		errs = append(errs, "strings feature requires a non-empty string.delimiter")
	}
	if s.Features.Comments && s.Comments.Line == "" && !s.Features.BlockComments {
		errs = append(errs, "comments feature enabled but no line comment and no block comments configured")
	}
	if s.Features.BlockComments && (s.Comments.BlockStart == "" || s.Comments.BlockEnd == "") {
		errs = append(errs, "block_comments feature requires comment.block_start and comment.block_end")
	}
	if !s.Features.AutoSemicolons && !hasSymbol(s.Symbols, "SEMICOLON") {
		errs = append(errs, "auto_semicolons disabled but no \";\" symbol defined (statements could not be separated)")
	}
	if s.Features.ForLoops && !hasSymbol(s.Symbols, "SEMICOLON") {
		errs = append(errs, "for_loops feature requires a \";\" symbol (for clauses are separated with semicolons)")
	}

	for _, name := range s.StatementEnd {
		if !tk.IsTypeName(name) {
			errs = append(errs, fmt.Sprintf("statement_end_tokens contains unknown token type %q", name))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("invalid language config:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}

func hasSymbol(symbols map[string]string, typeName string) bool {
	for _, v := range symbols {
		if v == typeName {
			return true
		}
	}
	return false
}

func hasKeyword(keywords map[string]string, typeName string) bool {
	for _, v := range keywords {
		if v == typeName {
			return true
		}
	}
	return false
}

// SortedKeys returns the keys of a string map sorted alphabetically. It is
// used by the generator to emit deterministic Go source.
func SortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
