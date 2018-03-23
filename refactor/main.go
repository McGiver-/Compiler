package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/McGiver-/Compiler/refactor/Lex/scanner"
	"github.com/McGiver-/Compiler/refactor/Lex/token"
)

func main() {
	inputFile := os.Args[1]
	s := scanner.Scanner{}
	r, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("could not open file")
	}
	src, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal("could not get src")
	}
	fset := token.NewFileSet()                      // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(src)) // register input "file"
	s.Init(file, src, nil /* no error handler */, 2)

	// Repeated calls to Scan yield the token sequence found in the input.
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}

}
