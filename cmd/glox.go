package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/fadyZohdy/gLox/pkg/ast"
	"github.com/fadyZohdy/gLox/pkg/scanner"
)

var hadError bool
var hadRuntimeError bool

func main() {
	// f, err := os.Create("cpu.prof")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()
	if len(os.Args) == 1 {
		runPrompt()
	} else {
		runFile(os.Args[1])
	}
}

func runFile(fileName string) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("Could not open file: " + fileName)
		return
	}
	run(string(content))
	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	// TODO: support arrow keys and history
	interpreter := ast.NewInterpreter(runtimeError)
	for {
		fmt.Printf("> ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal()
			return
		}
		if len(line) == 1 {
			continue
		}
		source := string(line)
		scanner := scanner.NewScanner(source, func(line int, message string) { report(line, "", message) })
		tokens := scanner.ScanTokens()
		parser := ast.NewParser(tokens, report)
		stmts, err := parser.Parse()
		if err != nil {
			log.Println(err)
			return
		}
		// if expr, ok := stmts[0].(*ast.Expression); ok {
		for _, stmt := range stmts {
			if expr, ok := stmt.(*ast.Expression); ok {
				fmt.Println(interpreter.Evaluate(expr))
			} else {
				interpreter.Interpret([]ast.Stmt{stmt})
			}
			hadError = false
		}
	}
}

func run(source string) {
	scanner := scanner.NewScanner(source, func(line int, message string) { report(line, "", message) })
	tokens := scanner.ScanTokens()
	parser := ast.NewParser(tokens, report)
	stmts, err := parser.Parse()
	if err != nil {
		log.Println(err)
		return
	}
	interpreter := ast.NewInterpreter(runtimeError)
	interpreter.Interpret(stmts)
}

func runtimeError(err *ast.RuntimeError) {
	log.Println(err.Message, "[line ", err.Token.Line, "]")
	hadRuntimeError = true
}

func report(line int, where string, message string) {
	log.Println("[line ", line, "] Error", where, ": ", message)
	hadError = true
}
