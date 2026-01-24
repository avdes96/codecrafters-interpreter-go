package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/app/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/errs"
	"github.com/codecrafters-io/interpreter-starter-go/app/interpreter"
	"github.com/codecrafters-io/interpreter-starter-go/app/parser"
	"github.com/codecrafters-io/interpreter-starter-go/app/scanner"
	"github.com/codecrafters-io/interpreter-starter-go/app/utils"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh (tokenize | parse) <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	filename := os.Args[2]
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "tokenize":
		tokenize(string(source))
		if errs.HadError {
			os.Exit(65)
		}
	case "parse":
		parse(string(source))
	case "evaluate":
		evaluate(string(source))
	case "run":
		run(string(source))
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func tokenize(source string) {
	s := scanner.NewScanner(source)
	for _, t := range s.ScanTokens() {
		fmt.Println(t)
	}
	fmt.Println()
}

func parse(source string) {
	s := scanner.NewScanner(source)
	tokens := s.ScanTokens()
	p := parser.NewParser(tokens)
	expr := p.ParseExpression()
	if errs.HadError {
		os.Exit(65)
	}
	ast.Print(expr)
}

func evaluate(source string) {
	s := scanner.NewScanner(source)
	tokens := s.ScanTokens()
	p := parser.NewParser(tokens)
	expr := p.ParseExpression()
	if errs.HadError {
		os.Exit(65)
	}
	i := interpreter.NewInterpreter()
	eval := i.Evaluate(expr)
	if errs.HadRuntimeError {
		os.Exit(70)
	}
	fmt.Println(utils.Stringify(eval))
}

func run(source string) {
	s := scanner.NewScanner(source)
	tokens := s.ScanTokens()
	p := parser.NewParser(tokens)
	statements := p.ParseProgram()
	if errs.HadError {
		os.Exit(65)
	}
	i := interpreter.NewInterpreter()
	i.Interpret(statements)
	if errs.HadRuntimeError {
		os.Exit(70)
	}
}
