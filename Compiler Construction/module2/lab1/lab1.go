package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	inPath := flag.String("i", "", "input .go file")
	outPath := flag.String("o", "", "output .go file (optional; default stdout)")
	flag.Parse()

	if *inPath == "" {
		fmt.Fprintln(os.Stderr, "usage: go run lab2.go -i demo.go [-o demo_out.go]")
		os.Exit(2)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, *inPath, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	rewritten := rewriteSwitches(file)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		fmt.Fprintf(os.Stderr, "format error: %v\n", err)
		os.Exit(1)
	}

	if *outPath == "" {
		_, _ = os.Stdout.Write(buf.Bytes())
	} else {
		if err := os.WriteFile(*outPath, buf.Bytes(), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "write error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Fprintf(os.Stderr, "done: rewritten %d switch statement(s)\n", rewritten)
}

func rewriteSwitches(file *ast.File) int {
	count := 0

	ast.Inspect(file, func(n ast.Node) bool {
		sw, ok := n.(*ast.SwitchStmt)
		if !ok {
			return true
		}

		if sw.Tag == nil {
			return true
		}

		tag := sw.Tag
		sw.Tag = nil

		for _, stmt := range sw.Body.List {
			cl, ok := stmt.(*ast.CaseClause)
			if !ok {
				continue
			}

			if len(cl.List) == 0 {
				continue
			}

			cond := buildOrEquals(tag, cl.List)
			cl.List = []ast.Expr{cond}
		}

		count++
		return true
	})

	return count
}

func buildOrEquals(tag ast.Expr, values []ast.Expr) ast.Expr {
	var result ast.Expr

	for i, v := range values {
		eq := &ast.BinaryExpr{
			X:  tag,
			Op: token.EQL,
			Y:  v,
		}

		if i == 0 {
			result = eq
		} else {
			result = &ast.BinaryExpr{
				X:  result,
				Op: token.LOR,
				Y:  eq,
			}
		}
	}

	return result
}
