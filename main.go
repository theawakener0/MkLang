package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/theawakener0/MkLang/config"
	"github.com/theawakener0/MkLang/generator"
)

// The engine packages are embedded and copied into every generated project.
// The generator renders config-specific replacements for token/token.go,
// token/spec.go, main.go, go.mod, and README.md on top of them.
//
//go:embed token lexer parser ast evaluator object repl LICENSE .gitignore
var engineFS embed.FS

var version = "dev"

const usageText = `mklang - generate a programming language from a JSON config

Usage:
  mklang init [-out file.json]                  write a template config file
  mklang generate -config file.json [-out dir]  generate a language project
  mklang -version                               print the version

Examples:
  mklang init
  mklang generate -config lang.json -out ./mylang
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usageText)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		runInit(os.Args[2:])
	case "generate":
		runGenerate(os.Args[2:])
	case "-version", "--version", "version":
		fmt.Printf("mklang %s\n", version)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s", os.Args[1], usageText)
		os.Exit(1)
	}
}

func runInit(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	out := fs.String("out", "lang.json", "output file for the template config")
	_ = fs.Parse(args)

	spec := config.Default()
	if err := spec.WriteTo(*out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("wrote template config to %s\n", *out)
	fmt.Printf("edit it, then run: mklang generate -config %s\n", *out)
}

func runGenerate(args []string) {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	configPath := fs.String("config", "", "path to the language config JSON")
	outDir := fs.String("out", "", "output directory (default: ./<name>)")
	force := fs.Bool("force", false, "overwrite an existing output directory")
	_ = fs.Parse(args)

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "generate requires -config <file.json>")
		os.Exit(1)
	}

	spec, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := spec.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	dir := *outDir
	if dir == "" {
		dir = "./" + strings.ToLower(spec.Name)
	}

	if err := generator.Generate(spec, engineFS, dir, *force); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin := strings.ToLower(spec.Name)
	fmt.Printf("generated language %q in %s\n", spec.Name, dir)
	fmt.Printf("  cd %s\n  go build -o %s .\n  ./%s\n", dir, bin, bin)
}
