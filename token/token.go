package token

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

// keywords maps the configured spelling of each keyword to its canonical type.
var keywords = map[string]TokenType{
	"fn":       FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"elseif":   ELSEIF,
	"else":     ELSE,
	"return":   RETURN,
	"for":      FOR,
	"loop":     LOOP,
	"break":    BREAK,
	"continue": CONTINUE,
	"null":     NULL,
}

// symbols maps the configured spelling of each operator/delimiter to its
// canonical type. Longer spellings are matched first by the lexer.
var symbols = map[string]TokenType{
	"=":  ASSIGN,
	"+":  PLUS,
	"-":  MINUS,
	"!":  BANG,
	"*":  ASTERISK,
	"/":  SLASH,
	"<":  LT,
	">":  GT,
	"==": EQ,
	"!=": NOTEQ,
	"<=": LTEQ,
	">=": GTEQ,
	"+=": INCASSIGN,
	"-=": DECDASSIGN,
	"*=": MLTASSIGN,
	"/=": DIVASSIGN,
	"&&": LAND,
	"||": LOR,
	"++": INC,
	"--": DEC,
	":=": ASSIGNCHAR,
	",":  COMMA,
	":":  COLOMN,
	";":  SEMICOLON,
	".":  DOT,
	"(":  LPAREN,
	")":  RPAREN,
	"{":  LBRACE,
	"}":  RBRACE,
	"[":  LBRACKET,
	"]":  RBRACKET,
}
