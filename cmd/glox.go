package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/fadyZohdy/gLox/pkg/ast"
	"github.com/fadyZohdy/gLox/pkg/scanner"
)

var hadError bool

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
}

func runPrompt() {
	for {
		fmt.Printf("> ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Could not read from stdin")
			return
		}
		if len(line) == 0 {
			log.Println("break")
			break
		}
		run(string(line))
		hadError = false
	}
}

func run(source string) {
	// TODO: check why slow
	t1 := time.Now().UnixNano() / int64(time.Millisecond)
	scanner := scanner.NewScanner(source, error)
	tokens := scanner.ScanTokens()
	t2 := time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Println("scanner: ", t2-t1)
	parser := ast.NewParser(tokens, report)
	parser.Parse()
	t3 := time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Println("Parser: ", t3-t2)
}

func error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	log.Println("[line ", line, "] Error", where, ": ", message)
	hadError = true
}
