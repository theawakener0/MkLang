package lexer

import (
	"strings"

	tk "github.com/theawakener0/MkLang/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	lastToken    tk.TokenType
	spec         tk.Spec
	stmtEnd      map[tk.TokenType]bool
	symbols      []string
}

func New(input string) *Lexer {
	return NewWithSpec(input, tk.Current)
}

func NewWithSpec(input string, spec tk.Spec) *Lexer {
	l := &Lexer{input: input, spec: spec}
	l.symbols = tk.SymbolSpellings()
	l.stmtEnd = make(map[tk.TokenType]bool, len(spec.StatementEnd))
	for _, t := range spec.StatementEnd {
		l.stmtEnd[t] = true
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition

	l.readPosition += 1
}

func (l *Lexer) readChars(n int) {
	for i := 0; i < n; i++ {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func newToken(tokenType tk.TokenType, ch byte) tk.Token {
	return tk.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) NextToken() tk.Token {
	tok := l.nextToken()
	l.lastToken = tok.Type
	return tok
}

func (l *Lexer) nextToken() tk.Token {
	var tok tk.Token

	if l.skipWhitespaceAndComments() && l.spec.Features.AutoSemicolons && l.stmtEnd[l.lastToken] {
		return tk.Token{Type: tk.SEMICOLON, Literal: "\n"}
	}

	if l.spec.Features.Strings && l.spec.String.Delimiter != "" && l.ch == l.spec.String.Delimiter[0] {
		tok.Type = tk.STRING
		tok.Literal = l.readString()
		l.readChar()
		return tok
	}

	if l.isIdentStart(l.ch) {
		tok.Literal = l.readIdentifier()
		tok.Type = tk.LookupIdent(tok.Literal)
		return tok
	}

	if isDigit(l.ch) {
		tok.Literal = l.readNumber()
		if strings.Contains(tok.Literal, ".") {
			tok.Type = tk.FLOAT
		} else {
			tok.Type = tk.INT
		}
		return tok
	}

	if l.spec.Features.Floats && l.ch == '.' && isDigit(l.peekChar()) {
		tok.Type = tk.FLOAT
		tok.Literal = l.readDotNumber()
		return tok
	}

	if spelling := l.matchSymbol(); spelling != "" {
		tt, _ := tk.LookupSymbol(spelling)
		tok = tk.Token{Type: tt, Literal: spelling}
		l.readChars(len(spelling))
		return tok
	}

	if l.ch == 0 {
		tok.Literal = ""
		tok.Type = tk.EOF
		return tok
	}

	tok = newToken(tk.ILLEGAL, l.ch)
	l.readChar()
	return tok
}

// matchSymbol returns the longest configured symbol spelling that starts at
// the current position, or "" if none matches.
func (l *Lexer) matchSymbol() string {
	rest := l.input[l.position:]
	for _, s := range l.symbols {
		if len(s) <= len(rest) && rest[:len(s)] == s {
			return s
		}
	}
	return ""
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for l.isIdentPart(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.spec.Features.Floats && l.ch == '.' {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readDotNumber() string {
	position := l.position

	l.readChar()

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) isIdentStart(ch byte) bool {
	if l.spec.Identifier.StartLetters && isLetter(ch) {
		return true
	}
	if l.spec.Identifier.ExtraStart != "" && strings.ContainsRune(l.spec.Identifier.ExtraStart, rune(ch)) {
		return true
	}
	return false
}

func (l *Lexer) isIdentPart(ch byte) bool {
	if l.isIdentStart(ch) {
		return true
	}
	if l.spec.Identifier.PartLettersDigits && isDigit(ch) {
		return true
	}
	if l.spec.Identifier.ExtraPart != "" && strings.ContainsRune(l.spec.Identifier.ExtraPart, rune(ch)) {
		return true
	}
	return false
}

func (l *Lexer) skipWhitespace() bool {
	crossedNewline := false
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			crossedNewline = true
		}
		l.readChar()
	}
	return crossedNewline
}

func (l *Lexer) skipWhitespaceAndComments() bool {
	crossedNewline := false

	for {
		if l.skipWhitespace() {
			crossedNewline = true
		}

		if !l.atComment() {
			break
		}

		l.skipComment()
	}

	return crossedNewline
}

func (l *Lexer) atComment() bool {
	if !l.spec.Features.Comments {
		return false
	}
	if l.spec.Comments.Line != "" && l.hasPrefix(l.spec.Comments.Line) {
		return true
	}
	if l.spec.Features.BlockComments && l.spec.Comments.BlockStart != "" && l.hasPrefix(l.spec.Comments.BlockStart) {
		return true
	}
	return false
}

func (l *Lexer) skipComment() {
	if l.spec.Comments.Line != "" && l.hasPrefix(l.spec.Comments.Line) {
		l.skipLineComment()
		return
	}
	l.skipBlockComment()
}

func (l *Lexer) skipLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipBlockComment() {
	l.readChars(len(l.spec.Comments.BlockStart))
	for l.ch != 0 {
		if l.hasPrefix(l.spec.Comments.BlockEnd) {
			l.readChars(len(l.spec.Comments.BlockEnd))
			return
		}
		l.readChar()
	}
}

func (l *Lexer) hasPrefix(prefix string) bool {
	if prefix == "" || l.position+len(prefix) > len(l.input) {
		return false
	}
	return l.input[l.position:l.position+len(prefix)] == prefix
}

func (l *Lexer) readString() string {
	delim := l.spec.String.Delimiter[0]
	l.readChar()

	var sb strings.Builder
	for l.ch != delim && l.ch != 0 {
		if l.spec.String.EscapeSequences && l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case 'r':
				sb.WriteByte('\r')
			case '\\':
				sb.WriteByte('\\')
			case delim:
				sb.WriteByte(delim)
			case 'x':
				l.readChar()
				c1 := l.ch
				l.readChar()
				c2 := l.ch
				hi, hiOK := hexVal(c1)
				lo, loOK := hexVal(c2)
				if hiOK && loOK {
					sb.WriteByte(hi*16 + lo)
				} else {
					sb.WriteString(`\x`)
					if c1 != 0 {
						sb.WriteByte(c1)
					}
					if c2 != 0 {
						sb.WriteByte(c2)
					}
				}
			case 0:
				sb.WriteByte('\\')
			default:
				sb.WriteByte('\\')
				sb.WriteByte(l.ch)
			}
			l.readChar()
			continue
		}
		sb.WriteByte(l.ch)
		l.readChar()
	}

	return sb.String()
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func hexVal(ch byte) (byte, bool) {
	switch {
	case '0' <= ch && ch <= '9':
		return ch - '0', true
	case 'a' <= ch && ch <= 'f':
		return ch - 'a' + 10, true
	case 'A' <= ch && ch <= 'F':
		return ch - 'A' + 10, true
	}
	return 0, false
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
