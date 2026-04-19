package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: astprint <path-to-file.go>")
		return
	}

	inputPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Println("invalid path:", err)
		return
	}

	fset := token.NewFileSet()

	file, err := parser.ParseFile(
		fset,
		inputPath,
		nil,
		parser.ParseComments,
	)
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}

	fmt.Println("=== AST for:", inputPath, "===\n")
	ast.Fprint(os.Stdout, fset, file, nil)
}
