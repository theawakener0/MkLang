package generator

// tokenSkeleton is the fixed part of every generated token/token.go. Only the
// keywords and symbols maps below it vary per config. It must stay in sync
// with the tool's own token/token.go.
const tokenSkeleton = `package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING"

	ASSIGN   = "ASSIGN"
	PLUS     = "PLUS"
	MINUS    = "MINUS"
	BANG     = "BANG"
	ASTERISK = "ASTERISK"
	SLASH    = "SLASH"

	LT = "LT"
	GT = "GT"

	EQ         = "EQ"
	NOTEQ      = "NOTEQ"
	LTEQ       = "LTEQ"
	GTEQ       = "GTEQ"
	INCASSIGN  = "INCASSIGN"
	DECDASSIGN = "DECDASSIGN"
	MLTASSIGN  = "MLTASSIGN"
	DIVASSIGN  = "DIVASSIGN"
	LAND       = "LAND"
	LOR        = "LOR"
	INC        = "INC"
	DEC        = "DEC"
	ASSIGNCHAR = "ASSIGNCHAR"

	COMMA     = "COMMA"
	COLOMN    = "COLOMN"
	SEMICOLON = "SEMICOLON"
	DOT       = "DOT"

	LPAREN   = "LPAREN"
	RPAREN   = "RPAREN"
	LBRACE   = "LBRACE"
	RBRACE   = "RBRACE"
	LBRACKET = "LBRACKET"
	RBRACKET = "RBRACKET"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSEIF   = "ELSEIF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	FOR      = "FOR"
	LOOP     = "LOOP"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	NULL     = "NULL"
)

// allTypes lists every canonical token type. It is used for config validation.
var allTypes = []TokenType{
	ILLEGAL, EOF, IDENT, INT, FLOAT, STRING,
	ASSIGN, PLUS, MINUS, BANG, ASTERISK, SLASH,
	LT, GT, EQ, NOTEQ, LTEQ, GTEQ,
	INCASSIGN, DECDASSIGN, MLTASSIGN, DIVASSIGN,
	LAND, LOR, INC, DEC, ASSIGNCHAR,
	COMMA, COLOMN, SEMICOLON, DOT,
	LPAREN, RPAREN, LBRACE, RBRACE, LBRACKET, RBRACKET,
	FUNCTION, LET, TRUE, FALSE, IF, ELSEIF, ELSE,
	RETURN, FOR, LOOP, BREAK, CONTINUE, NULL,
}

`

// specSkeleton is the fixed part of every generated token/spec.go (the Spec
// types). The defaultSpec() body is rendered per config. It must stay in sync
// with the tool's own token/spec.go.
const specSkeleton = `package token

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

func defaultSpec() Spec {
	return Spec{
`

// specTail closes defaultSpec and declares the active spec.
const specTail = `	}
}

// Current describes the language currently loaded by this binary.
var Current = defaultSpec()
`
