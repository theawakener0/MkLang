# MkLang

> Describe a programming language in one JSON file — MkLang builds the whole thing.

MkLang is a language *factory*. You write a JSON config that describes your
language — its name, keywords, operators, and features — and MkLang generates a
complete, self-contained Go project with a **lexer, parser, evaluator, and
REPL**, ready to `go build` and run. No compiler theory required.

> [!IMPORTANT]
> MkLang is in early development. It's not ready for production use yet.
> MkLang is based on an early version of [the Zod language](https://github.com/theawakener0/Zod).

Bundled example languages:

- [Zod](langs/zod.json) — the default language. A small dynamically typed
  language with functions, closures, arrays, hashes, matrices, and `try`.
- [Pythonish](examples/configs/pythonish.json) — a Python-flavored dialect
  (`def`, `elif`, `and`/`or`/`not`, `#` comments, `True`/`False`/`None`).

---

## What you get

Every generated project is a working interpreter with:

- A **tokenizer and parser** for exactly the syntax you configured
- A **tree-walking evaluator** with first-class functions and closures
- An interactive **REPL** and a **script runner** (`./mylang program.txt`)
- **Arrays**, **hashes**, **matrices**, **loops**, and **comments** (opt-in)
- **Zero external dependencies** — Go standard library only
- A `lang.json` copy of your config, so you can tweak and regenerate later

## Requirements

- [Go](https://go.dev/dl/) **1.22 or newer** — to run the generator, and to
  build the languages it generates. You don't need any other tools.

---

## Install

### Download script (Linux / macOS)

```sh
curl -fsSL https://raw.githubusercontent.com/theawakener0/MkLang/main/install.sh | sh
```

### go install

```sh
go install github.com/theawakener0/MkLang@latest
```

### Prebuilt binaries

Windows (and everyone else) can grab a binary from the
[GitHub Releases](https://github.com/theawakener0/MkLang/releases) page.

Check that it works:

```sh
mklang -version
```

---

## Quickstart

### Try a bundled language (30 seconds)

```sh
mklang generate -config langs/zod.json      # writes ./zod
cd zod
go build -o zod .
./zod                                       # the REPL
./zod ../examples/hello.zd                  # or run a script
```

### Create your own language (5 minutes)

```sh
# 1. Write a starter config
mklang init                                 # creates lang.json

# 2. Edit it: name, keywords, operators, features, output style
vim lang.json

# 3. Generate the language project
mklang generate -config lang.json -out ./mylang

# 4. Build and run your new language
cd mylang
go build -o mylang .
./mylang                                    # REPL
./mylang hello.txt                          # script
```

That's the whole loop. Everything after this page is details about step 2.

---

## How it works

MkLang embeds a **generic interpreter engine** (token → lexer → parser →
evaluator → object). The engine is *spec-driven*: it reads every grammar
decision from a `Spec` built from your config at runtime. Generating a language
renders a handful of config-specific files and copies the engine unchanged:

| Generated file        | Comes from your config                       |
| --------------------- | -------------------------------------------- |
| `token/token.go`      | `keywords` and `symbols` tables              |
| `token/spec.go`       | the whole config (as `defaultSpec()`)        |
| `main.go`, `go.mod`, `README.md`, `lang.json` | name, version, module |

Each generated project is **self-contained** (no MkLang dependency) and depends
only on the Go standard library. The same interpreter implements every
generated language, so your config *is* the language.

---

## The config file

A config is plain JSON. **Any field you omit falls back to the bundled
defaults**, so a config only needs to say what you change. One important
exception:

> Providing `keywords`, `symbols`, or `builtins` **replaces** the default table
> entirely — deleting a line from an edited template removes that token from
> your language.

`mklang init` writes the full template (the bundled Zod language). Here is the
exact output, unchanged:

```json
{
  "name": "Zod",
  "module": "github.com/yourname/yourlang",
  "version": "0.1.0",
  "features": {
    "floats": true,
    "strings": true,
    "arrays": true,
    "hashes": true,
    "for_loops": true,
    "infinite_loops": true,
    "try": true,
    "prefix_postfix_inc_dec": true,
    "compound_assignment": true,
    "elseif": true,
    "matrices": true,
    "auto_semicolons": true,
    "comments": true,
    "block_comments": true
  },
  "comments": {
    "line": "//",
    "block_start": "/*",
    "block_end": "*/"
  },
  "identifier": {
    "start_letters": true,
    "extra_start": "_",
    "part_letters_digits": true,
    "extra_part": "_"
  },
  "string": {
    "delimiter": "\"",
    "escape_sequences": true
  },
  "keywords": {
    "break": "BREAK",
    "continue": "CONTINUE",
    "else": "ELSE",
    "elseif": "ELSEIF",
    "false": "FALSE",
    "fn": "FUNCTION",
    "for": "FOR",
    "if": "IF",
    "let": "LET",
    "loop": "LOOP",
    "null": "NULL",
    "return": "RETURN",
    "true": "TRUE"
  },
  "symbols": {
    "!": "BANG",
    "!=": "NOTEQ",
    "\u0026\u0026": "LAND",
    "(": "LPAREN",
    ")": "RPAREN",
    "*": "ASTERISK",
    "*=": "MLTASSIGN",
    "+": "PLUS",
    "++": "INC",
    "+=": "INCASSIGN",
    ",": "COMMA",
    "-": "MINUS",
    "--": "DEC",
    "-=": "DECDASSIGN",
    ".": "DOT",
    "/": "SLASH",
    "/=": "DIVASSIGN",
    ":": "COLOMN",
    ":=": "ASSIGNCHAR",
    ";": "SEMICOLON",
    "\u003c": "LT",
    "\u003c=": "LTEQ",
    "=": "ASSIGN",
    "==": "EQ",
    "\u003e": "GT",
    "\u003e=": "GTEQ",
    "[": "LBRACKET",
    "]": "RBRACKET",
    "{": "LBRACE",
    "||": "LOR",
    "}": "RBRACE"
  },
  "builtins": {
    "color": "color",
    "contains": "contains",
    "exp": "exp",
    "first": "first",
    "float": "float",
    "input": "input",
    "insert": "insert",
    "int": "int",
    "keys": "keys",
    "last": "last",
    "len": "len",
    "make": "make",
    "matrix": "matrix",
    "pi": "pi",
    "pop": "pop",
    "printf": "printf",
    "println": "println",
    "push": "push",
    "random": "random",
    "remove": "remove",
    "sleep": "sleep",
    "string": "string",
    "try": "try",
    "type": "type",
    "vals": "vals"
  },
  "statement_end_tokens": [
    "IDENT",
    "INT",
    "FLOAT",
    "STRING",
    "TRUE",
    "FALSE",
    "NULL",
    "RPAREN",
    "RBRACKET",
    "RBRACE",
    "INC",
    "DEC",
    "RETURN",
    "BREAK",
    "CONTINUE"
  ],
  "output": {
    "true": "true",
    "false": "false",
    "null": "null",
    "break": "break",
    "continue": "continue",
    "function": "fn",
    "array_open": "[",
    "array_close": "]",
    "array_separator": ", ",
    "hash_open": "{",
    "hash_close": "}",
    "hash_separator": ", ",
    "hash_keyval": ": "
  }
}
```

To adapt it to your own language, change `name`, `module`, and `version`, then
tweak the spellings and features you care about. The
[field reference](#field-reference) table below explains what each part does.

### Field reference

| Field                     | What it does                                                            |
| ------------------------- | ----------------------------------------------------------------------- |
| `name`                    | Language name; also the default output directory and binary name (lowercased). |
| `module`                  | Go module path used in the generated `go.mod`.                          |
| `version`                 | Shown in the REPL banner.                                               |
| `features`                | Turn whole parts of the language on or off.                             |
| `comments`                | Comment markers: a `line` prefix and optional `block_start`/`block_end`. |
| `identifier`              | Which characters can form identifiers.                                  |
| `string`                  | String delimiter and escape sequences.                                  |
| `keywords`                | Reserved words → canonical token type. You pick the spellings.          |
| `symbols`                 | Operators/delimiters → canonical token type. Longer spellings match first. |
| `builtins`                | Spelling → built-in function, so you can rename them (`"print": "println"`). |
| `statement_end_tokens`    | Token types that end a statement at a newline (used with `auto_semicolons`). |
| `output`                  | How values print: booleans, `null`, functions, and array/hash delimiters. |

### Canonical token types

The right-hand side of each `keywords`/`symbols` entry must be one of these
canonical names (the *meaning* is fixed; only the *spelling* is up to you):

- **Keywords:** `FUNCTION`, `LET`, `TRUE`, `FALSE`, `NULL`, `IF`, `ELSEIF`,
  `ELSE`, `RETURN`, `FOR`, `LOOP`, `BREAK`, `CONTINUE`
- **Symbols:** `ASSIGN`, `ASSIGNCHAR`, `EQ`, `NOTEQ`, `LAND`, `LOR`, `BANG`,
  `PLUS`, `MINUS`, `ASTERISK`, `SLASH`, `LT`, `LTEQ`, `GT`, `GTEQ`, `INC`,
  `DEC`, `INCASSIGN`, `DECDASSIGN`, `MLTASSIGN`, `DIVASSIGN`, `LPAREN`,
  `RPAREN`, `LBRACE`, `RBRACE`, `LBRACKET`, `RBRACKET`, `COMMA`, `COLOMN`,
  `SEMICOLON`, `DOT`

For example, these three configs all mean the same thing — just spelled
differently:

| Language | FUNCTION keyword | not equals | boolean true | declares with |
| -------- | ---------------- | ---------- | ------------ | ------------- |
| Zod      | `fn`             | `!=`       | `true`       | `let` / `:=`  |
| Pythonish| `def`            | `!=`       | `True`       | `=` (`:=`)    |
| Rubish   | `fn`             | `!=`       | `true`       | `:=`          |

### Built-in functions

Every canonical builtin you map in `builtins` is available in the generated
language:

- **IO:** `println`, `printf`, `input`
- **Types:** `int`, `float`, `string`, `type`
- **Collections:** `len`, `first`, `last`, `push`, `pop`, `insert`, `remove`,
  `keys`, `vals`, `contains`, `make`
- **Math:** `exp`, `pi`, `random`
- **Other:** `color`, `sleep`, `matrix` (needs `features.matrices`), `try`
  (needs `features.try`)

---

## The fixed grammar

No matter what you configure, every generated language shares this syntax
(only the *spellings* change):

- **Functions are anonymous** — you bind them to a variable yourself:
  `add = fn (x, y) { x + y }`
- **Blocks always use braces** `{ }`
- **Conditionals take parentheses**: `if (x > 5) { ... }`
- **`for` is C-style**: `for (init; cond; update) { ... }` (requires a `;`
  symbol); `loop { ... }` runs forever
- **Variables:** a language needs a `LET` keyword *or* a symbol mapped to
  `ASSIGNCHAR`. The `LET` keyword and the `ASSIGNCHAR` spelling are what
  *declare* a new variable; a symbol mapped to `ASSIGN` (e.g. `=`) **only
  reassigns** an existing one. With Zod, `let x = 5` and `x := 5` declare;
  with Pythonish, `x = 5` declares (because `=` is the `ASSIGNCHAR`).
- **Operator precedence** is the usual one: `||` < `&&` < `==`/`!=` < `<`/`>`
  < `+`/`-` < `*`/`/` < unary prefix < call/index.

---

## Customizing: a worked example

Say you want a language with Python-ish spellings but a `&`/`|` bit-flavored
look, using `!` for not-equal and `~` for negation. Your config just overrides
the spellings:

```jsonc
{
  "name": "Punky",
  "module": "github.com/you/punky",
  "version": "0.1.0",
  "features": {
    "strings": true,
    "arrays": true,
    "hashes": true,
    "auto_semicolons": true
  },
  "comments": { "line": "#" },
  "keywords": {
    "func": "FUNCTION",
    "yeah": "TRUE",
    "nah": "FALSE",
    "nada": "NULL",
    "if": "IF", "elseif": "ELSEIF", "else": "ELSE",
    "give": "RETURN", "for": "FOR", "loop": "LOOP",
    "stop": "BREAK", "continue": "CONTINUE"
  },
  "symbols": {
    "=": "ASSIGNCHAR", "==": "EQ", "!": "NOTEQ", "~": "BANG",
    "+": "PLUS", "-": "MINUS", "*": "ASTERISK", "/": "SLASH",
    "+=": "INCASSIGN", "-=": "DECDASSIGN", "*=": "MLTASSIGN", "/=": "DIVASSIGN",
    "<": "LT", "<=": "LTEQ", ">": "GT", ">=": "GTEQ",
    "&&": "LAND", "||": "LOR",
    "(": "LPAREN", ")": "RPAREN", "{": "LBRACE", "}": "RBRACE",
    "[": "LBRACKET", "]": "RBRACKET",
    ",": "COMMA", ":": "COLOMN", ";": "SEMICOLON", ".": "DOT"
  },
  "builtins": {
    "print": "println", "len": "len", "int": "int", "string": "string",
    "push": "push", "pop": "pop"
  },
  "output": { "true": "yeah", "false": "nah", "null": "nada", "function": "func" }
}
```

Because `=` is mapped to `ASSIGNCHAR`, a bare `x = 5` declares the variable
(the LET keyword is only needed when you want `let x = 5` style). Now this is
valid Punky:

```
x = 5
double = func (n) { give n * 2 }
print(double(x))          # 10
if (x ! 3) { print("x is not 3") }
```

Run `mklang generate -config punky.json` and it just works. Because
`keywords`/`symbols` replace the defaults, every spelling you want must be in
the table — that's how you *remove* tokens too.

---

## Regenerating (the iteration loop)

Your config travels with the project as `lang.json`. After editing it, rebuild:

```sh
cd mylang
mklang generate -config lang.json -out . -force
go build -o mylang .
```

- `-out .` writes in place; `-force` is required because the directory already
  exists (generation refuses to overwrite a non-empty directory without it).
- Editing `lang.json` inside the project keeps everything in one place.

---

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
./zod ../examples/guess_the_number.zd
```

For a taste of a different dialect, generate Pythonish and run its demo:

```sh
mklang generate -config examples/configs/pythonish.json -out ./pysh
cd pysh
go build -o pysh .
./pysh ../examples/demo.pysh
```

Example programs to read: `examples/variables.zd`, `examples/control_flow.zd`,
`examples/arrays.zd`, `examples/closures.zd`, `examples/loops.zd`,
`examples/hashes.zd`.

---

## CLI

```
mklang init [-out file.json]                  write a template config
mklang generate -config file.json [-out dir]  generate a language project
mklang -version                               print the version
```

- `mklang generate` **validates your config** before writing anything and tells
  you what it found wrong.
- With no `-out`, the project is written to `./<lowercase-name>`.

---

## Troubleshooting

| Symptom | Cause & fix |
| ------- | ----------- |
| `no prefix parse function for LBRACE found` | You typed a word that isn't in your `keywords` table (e.g. `func` when your FUNCTION keyword is `fn`). Every reserved word must be listed under `keywords`. |
| `expected next token to be ...` | A spelling in your program doesn't match the config, or you used an operator you removed from `symbols`. |
| `output directory "..." is not empty` | Add `-force` to overwrite an existing directory. |
| `no way to declare variables` | Config validation: add a `LET` keyword or a symbol mapped to `ASSIGNCHAR` (`:=`, or `=` in Pythonish style). A symbol mapped to plain `ASSIGN` (`=`) only reassigns. |
| `for_loops feature requires a ";" symbol` | The C-style `for` clause needs `;` in `symbols`. |
| Keyword/operator still uses the old spelling after regenerating | `keywords`, `symbols`, and `builtins` replace the whole table. Keep every line you still want; deleting a line removes the token. |
| `{}` can't be used as a literal | `features.hashes` is off — braces are only for blocks. Turn `hashes` on. |
| Unexpected behavior with newlines | `features.auto_semicolons` inserts statement boundaries at line breaks. Disable it and use explicit `;`. |

---

## Development

```sh
go test ./...
```

The engine tests run against the bundled Zod config; the generator tests
generate, build, and execute Pythonish end-to-end.

## License

MIT License. See [`LICENSE`](LICENSE).
