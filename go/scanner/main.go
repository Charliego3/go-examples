package main

import (
	"github.com/kataras/golog"
	"go/scanner"
	"go/token"
	"os"
)

func main() {
	src, err := os.ReadFile("/Users/charlie/dev/go/temp/go/scanner/main.go")
	if err != nil {
		golog.Fatal(err)
	}

	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	test()

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		golog.Errorf("Pos: %-*v Tok: %-*q Lit: %q", 10, fset.Position(pos), 13, tok, lit)
	}
}

func test() {
	golog.Println("this is test func")
}
