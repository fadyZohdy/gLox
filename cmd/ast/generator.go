package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("Usage: generate_ast <output directory>")
		os.Exit(65)
	}
	outputDir := os.Args[1]

	defineAst(outputDir, "Expr", []string{
		"Binary   : left Expr, operator scanner.Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value interface{}",
		"Unary    : operator scanner.Token, right Expr",
		"Ternary  : condition Expr, trueBranch Expr, falseBranch Expr",
	})

	defineAst(outputDir, "Stmt", []string{
		"Expression : expression Expr",
		"Print      : expression Expr",
	})
}

func defineAst(outputDir, basename string, types []string) {
	path := outputDir + "/" + strings.ToLower(basename) + ".go"

	//cleanup first
	os.Remove(path)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writeToFile(file, "package ast")
	writeToFile(file, "")
	writeToFile(file, "import \"github.com/fadyZohdy/gLox/pkg/scanner\"")
	writeToFile(file, "")
	writeToFile(file, "type "+basename+" interface {")
	writeToFile(file, "\t"+"accept(visitor Visitor) any")
	writeToFile(file, "}")
	writeToFile(file, "")

	for _, nodeType := range types {
		defineType(file, nodeType, basename)
	}
}

func writeToFile(file *os.File, content string) {
	_, err := fmt.Fprintln(file, content)
	if err != nil {
		log.Fatal(err)
	}
}

func defineType(file *os.File, nodeType string, basename string) {
	structName := strings.TrimSpace(strings.Split(nodeType, ":")[0])
	fieldsStr := strings.TrimSpace(strings.Split(nodeType, ":")[1])
	fields := strings.Split(fieldsStr, ",")
	fmt.Fprintln(file, "type "+structName+" struct {")
	for _, field := range fields {
		field = strings.TrimSpace(field)
		writeToFile(file, "\t"+field)
	}
	writeToFile(file, "}")
	writeToFile(file, "")

	writeToFile(file, "func (expr *"+structName+") accept(visitor Visitor) any {")
	writeToFile(file, "\t"+"return visitor.Visit"+structName+basename+"(expr)")
	writeToFile(file, "}")
	writeToFile(file, "")
}
