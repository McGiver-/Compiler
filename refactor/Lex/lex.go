package Lex

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/McGiver-/Compiler/refactor/Lex/scanner"
	"github.com/McGiver-/Compiler/refactor/Lex/token"
)

type Token struct {
	token.Token
	Lit string
	token.Position
}

type Lexer struct {
	scanner.Scanner
	fset *token.FileSet
}

func CreateLexer(r io.Reader) (*Lexer, error) {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not get src")
	}
	fset := token.NewFileSet()                      // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(src)) // register input "file"
	lexer := &Lexer{scanner.Scanner{}, fset}
	lexer.Init(file, src, nil /* no error handler */, 2)

	return lexer, nil
}

func (l *Lexer) GetTokens() (<-chan *Token, <-chan error) {
	tc, ec := make(chan *Token), make(chan error)

	// Repeated calls to Scan yield the token sequence found in the input.
	go func() {
		for {
			pos, tok, lit := l.Scan()
			switch tok {
			case token.ILLEGAL:
				ec <- fmt.Errorf("%s\t%s\tis not an accepted character", l.fset.Position(pos), lit)
			case token.EOF:
				break
			default:
				tc <- &Token{tok, lit, l.fset.Position(pos)}
			}
			if tok == token.EOF {
				break
			}
		}
	}()
	return tc, ec
}

func (l *Lexer) GetTokensNoChan() (tks []*Token, ers []error) {
	for {
		pos, tok, lit := l.Scan()
		switch tok {
		case token.ILLEGAL:
			ers = append(ers, fmt.Errorf("%s <%s> is not an accepted character", l.fset.Position(pos), lit))
		case token.EOF:
			break
		default:
			tks = append(tks, &Token{tok, lit, l.fset.Position(pos)})
		}
		if tok == token.EOF {
			break
		}
	}
	return
}
