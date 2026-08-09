package token

// Features toggles which parts of the language are available.
type Features struct {
	Floats              bool
	Strings             bool
	Arrays              bool
	Hashes              bool
	ForLoops            bool
	InfiniteLoops       bool
	Try                 bool
	PrefixPostfixIncDec bool
	CompoundAssignment  bool
	ElseIf              bool
	Matrices            bool
	AutoSemicolons      bool
	Comments            bool
	BlockComments       bool
}

// Comments describes the comment markers of the language.
type Comments struct {
	Line       string
	BlockStart string
	BlockEnd   string
}

// Identifier describes which characters may form identifiers.
type Identifier struct {
	StartLetters      bool
	ExtraStart        string
	PartLettersDigits bool
	ExtraPart         string
}

// StringConfig describes how string literals are written.
type StringConfig struct {
	Delimiter       string
	EscapeSequences bool
}

// Output describes how values are printed.
type Output struct {
	True       string
	False      string
	Null       string
	Break      string
	Continue   string
	Function   string
	ArrayOpen  string
	ArrayClose string
	ArraySep   string
	HashOpen   string
	HashClose  string
	HashSep    string
	HashKeyVal string
}

// Spec describes the syntax, features, and identity of the current language.
type Spec struct {
	Name         string
	Module       string
	Version      string
	Features     Features
	Comments     Comments
	Identifier   Identifier
	String       StringConfig
	Builtins     map[string]string // configured spelling -> canonical builtin name
	StatementEnd []TokenType       // token types that can end a statement
	Output       Output
}

// defaultSpec returns the built-in Zod language.
func defaultSpec() Spec {
	return Spec{
		Name:    "Zod",
		Module:  "github.com/theawakener0/MkLang",
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
		Builtins: map[string]string{
			"len": "len", "println": "println", "printf": "printf", "input": "input",
			"int": "int", "float": "float", "string": "string", "type": "type",
			"first": "first", "last": "last", "pop": "pop", "push": "push",
			"insert": "insert", "remove": "remove", "keys": "keys", "vals": "vals",
			"contains": "contains", "random": "random", "matrix": "matrix",
			"make": "make", "color": "color", "sleep": "sleep", "exp": "exp", "pi": "pi",
			"try": "try",
		},
		StatementEnd: []TokenType{
			IDENT, INT, FLOAT, STRING, TRUE, FALSE, NULL,
			RPAREN, RBRACKET, RBRACE, INC, DEC,
			RETURN, BREAK, CONTINUE,
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

// Current describes the language currently loaded by this binary.
var Current = defaultSpec()
