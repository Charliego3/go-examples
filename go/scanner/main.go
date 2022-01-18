package main

import (
	"github.com/kataras/golog"
	"go/scanner"
	"go/token"
)

func main() {
	src := []byte("cos(x) + 2i*sin(x) // Euler")

	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		golog.Errorf("Pos: %+v, Tok: %q, Lit: %q", fset.Position(pos), tok, lit)
	}
}
