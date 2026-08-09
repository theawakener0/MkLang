# MkLang

MkLang generates a complete, interpreted programming language from a single
JSON config file. Describe your syntax, features, and keyword spellings, and it
writes a self-contained Go project with a lexer, parser, evaluator, and REPL
tailored to your language — ready to `go build` and run.

Bundled example languages:

- [Zod](langs/zod.json) — the built-in language, a small dynamically typed
  language with functions, closures, arrays, hashes, and matrices.
- [Pythonish](examples/configs/pythonish.json) — a Python-flavored dialect
  (`def`, `elif`, `and`/`or`/`not`, `#` comments, `True`/`False`/`None`).

## Install

Prebuilt binaries for Linux, macOS, and Windows are attached to every
[GitHub Release](https://github.com/theawakener0/MkLang/releases).

### Download script

```sh
curl -fsSL https://raw.githubusercontent.com/theawakener0/MkLang/main/install.sh | sh
```

### go install

```sh
go install github.com/theawakener0/MkLang@latest
```

Requires [Go](https://go.dev/dl/) 1.22 or newer. (You only need Go to run the
generator; you also need it to build the language it generates.)

## Quickstart

```sh
# 1. Write a template config
mklang init

# 2. Edit lang.json: name, keywords, operators, features, output style
vim lang.json

# 3. Generate the language project
mklang generate -config lang.json -out ./mylang

# 4. Build and run your new language
cd mylang
go build -o mylang .
./mylang                     # REPL
./mylang program.txt         # run a script
```

Or generate one of the bundled languages directly:

```sh
mklang generate -config langs/zod.json
mklang generate -config examples/configs/pythonish.json
```

## How it works

MkLang embeds a generic interpreter engine (token → lexer → parser →
evaluator → object). The engine is *spec-driven*: every grammar decision is
read from a `Spec` built from your config. Generating a language renders three
config-specific files from that spec and copies the engine unchanged:

| Generated file            | From your config                        |
| ------------------------- | --------------------------------------- |
| `token/token.go`          | `keywords` and `symbols` tables         |
| `token/spec.go`           | `defaultSpec()` = the whole config      |
| `main.go`, `go.mod`, `README.md`, `lang.json` | name, version, module       |

The engine reads its config at runtime through `token.Current`, so the same
interpreter implements every generated language. Each generated project is
self-contained (no MkLang dependency) and uses only the Go standard library.

## The config file

A config is JSON. Any field you omit falls back to the bundled defaults, so a
config only needs to say what you change. Providing `keywords`, `symbols`, or
`builtins` **replaces** the default table entirely — deleting a line from an
edited template removes that token from your language.

`mklang init` writes the full template (the Zod language). Key fields:

```jsonc
{
  "name": "MyLang",                 // language name (also the default dir/binary)
  "module": "github.com/you/mylang",// Go module path for the generated project
  "version": "0.1.0",

  "features": {
    "floats": true,                 // float literals
    "strings": true,                // string literals
    "arrays": true,                 // [a, b, c]
    "hashes": true,                 // {k: v, ...}
    "for_loops": true,              // for (init; cond; update) { }
    "infinite_loops": true,         // loop { }
    "try": true,                    // try() builtin / recoverable errors
    "prefix_postfix_inc_dec": true, // ++ and --
    "compound_assignment": true,    // += -= *= /=
    "elseif": true,                 // elseif / else-if chains
    "matrices": true,               // matrix() builtin and arithmetic
    "auto_semicolons": true,        // newlines end statements
    "comments": true,               // line comments
    "block_comments": true          // /* */ comments
  },

  "comments":   { "line": "//", "block_start": "/*", "block_end": "*/" },
  "identifier": { "start_letters": true, "extra_start": "_", "part_letters_digits": true, "extra_part": "_" },
  "string":     { "delimiter": "\"", "escape_sequences": true },

  "keywords": { "fn": "FUNCTION", "let": "LET", "if": "IF", "elseif": "ELSEIF", "...": "..." },
  "symbols":  { "=": "ASSIGN", "==": "EQ", "&&": "LAND", "+": "PLUS", "...": "..." },
  "builtins": { "len": "len", "println": "println", "print": "println", "...": "..." },

  "statement_end_tokens": ["IDENT", "INT", "FLOAT", "STRING", "TRUE", "FALSE", "NULL", "RPAREN", "RBRACKET", "RBRACE", "RETURN", "BREAK", "CONTINUE"],

  "output": { "true": "true", "false": "false", "null": "null", "function": "fn", "...": "..." }
}
```

### Grammar notes

- `keywords` maps each reserved word to a canonical token type
  (`FUNCTION`, `LET`, `TRUE`, `FALSE`, `IF`, `ELSEIF`, `ELSE`, `RETURN`,
  `FOR`, `LOOP`, `BREAK`, `CONTINUE`, `NULL`). You choose the spellings.
- `symbols` maps operator/delimiter spellings to token types (`ASSIGN`,
  `ASSIGNCHAR` (`:=`), `EQ`, `NOTEQ`, `LAND`, `LOR`, `BANG`, `PLUS`, `MINUS`,
  `ASTERISK`, `SLASH`, `LT`, `LTEQ`, `GT`, `GTEQ`, `INC`, `DEC`, `INCASSIGN`,
  `DECDASSIGN`, `MLTASSIGN`, `DIVASSIGN`, `LPAREN`, `RPAREN`, `LBRACE`,
  `RBRACE`, `LBRACKET`, `RBRACKET`, `COMMA`, `COLOMN`, `SEMICOLON`, `DOT`).
  Longer spellings match first.
- `builtins` maps the spelling of each built-in function to its canonical
  name, so you can rename them (`"print": "println"`, `"str": "string"`).
- `output` controls how values print: `true`/`false`/`null`/`break`/
  `continue`, the `function` keyword, and array/hash delimiters.
- A language needs a way to declare variables: a `LET` keyword, or the
  `:=` (`ASSIGNCHAR`) symbol. `x = 5` alone only reassigns existing variables
  — the engine's `x = 5` creates nothing, `let x = 5` and `x := 5` do.
- `for_loops` requires a `;` symbol: the `for` clause is always C-style
  `for (init; cond; update) { }`.
- Conditions take parentheses: `if (x > 5) { }`.
- Blocks always use braces, and functions are anonymous
  (`f = fn (x) { x * 2 }`).

`mklang generate` validates your config before writing anything and explains
what it found wrong.

## Examples

The [examples](examples/) directory holds runnable programs for the bundled
Zod language:

```sh
mklang generate -config langs/zod.json -out ./zod
cd zod
go build -o zod .
./zod ../examples/hello.zd
./zod ../examples/closures.zd
./zod ../examples/hashes.zd
```

For a taste of a different dialect, generate Pythonish and run its demo:

```sh
mklang generate -config examples/configs/pythonish.json -out ./pysh
cd pysh
go build -o pysh .
./pysh ../examples/demo.pysh
```

## CLI

```
mklang init [-out file.json]                  write a template config
mklang generate -config file.json [-out dir]  generate a language project
mklang -version                               print the version
```

## Development

```sh
go test ./...
```

The engine tests run against the bundled Zod config; the generator tests
generate, build, and execute Pythonish end-to-end.

## License

MIT License. See [`LICENSE`](LICENSE).
