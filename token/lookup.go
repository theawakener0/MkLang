package token

import "sort"

var (
	symbolSpellings []string
	symbolByType    = make(map[TokenType]string)
	keywordsByType  = make(map[TokenType]string)
	typeNames       = make(map[string]bool)
)

func init() {
	for s := range symbols {
		symbolSpellings = append(symbolSpellings, s)
		symbolByType[symbols[s]] = s
	}
	sort.Slice(symbolSpellings, func(i, j int) bool {
		if len(symbolSpellings[i]) != len(symbolSpellings[j]) {
			return len(symbolSpellings[i]) > len(symbolSpellings[j])
		}
		return symbolSpellings[i] < symbolSpellings[j]
	})
	for s := range keywords {
		keywordsByType[keywords[s]] = s
	}
	for _, tt := range allTypes {
		typeNames[string(tt)] = true
	}
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func LookupSymbol(spelling string) (TokenType, bool) {
	tok, ok := symbols[spelling]
	return tok, ok
}

// SymbolSpellings returns every configured symbol spelling, longest first,
// so the lexer can greedily match multi-character operators.
func SymbolSpellings() []string {
	return symbolSpellings
}

// SymbolLiteral returns the configured spelling for a symbol token type,
// or "" if the type has no configured symbol.
func SymbolLiteral(t TokenType) string {
	return symbolByType[t]
}

// KeywordLiteral returns the configured spelling for a keyword token type,
// or "" if the type has no configured keyword.
func KeywordLiteral(t TokenType) string {
	return keywordsByType[t]
}

// IsTypeName reports whether s is the canonical name of a token type.
func IsTypeName(s string) bool {
	return typeNames[s]
}
