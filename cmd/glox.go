package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	scanner "github.com/fadyZohdy/gLox/pkg/scanner"
)

var hadError bool

func main() {
	if len(os.Args) > 2 {
		log.Println("Usage: jlox [script]")
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
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
	scanner := scanner.NewScanner(source, error)
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
}

func error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	log.Println("[line ", line, "] Error", where, ": ", message)
	hadError = true
}
