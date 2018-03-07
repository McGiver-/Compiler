package Lex

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

type Token struct {
	Type     string
	Lexeme   string
	Location string
}

type LScanner struct {
	reader   *bufio.Reader
	line     int
	col      int
	reserved map[string]string
}

// Reserved words with lexemes as their keys and types as their values
// The types must be atocc compatible
// Order matters. Put the longer one first i.e '==' '='
var reserved = map[string]string{
	"==": "eq", "=": "=", "+": "+",
	"-": "-", "*/": "*/", "*": "*",
	"<>": "neq", "<=": "leq", "<": "lt",
	">=": "geq", ">": "gt", "(": "(",
	")": ")", "{": "{", "}": "}",
	"[": "[", "]": "]", "::": "sr",
	":": ":", "else": "else", "then": "then",
	"//": "//", "/*": "/*", "/": "/",
	"for": "for", "if": "if", "class": "class",
	"and": "and", "int": "int", ";": ";",
	"not": "not", "float": "float", ",": ",",
	"or": "or", "get": "get", ".": ".",
	"put": "put", "return": "return", "program": "program",
}

// Create Scanner by passing the file to be scanned and return a pointer to LScanner
func CreateScanner(fileName string) (*LScanner, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	return &LScanner{reader, 1, 1, reserved}, nil
}

// Check if identifier and return length if so.
func (ls *LScanner) id() int {
	i := 1
	for {
		b, err := ls.reader.Peek(i)
		if err == io.EOF {
			return i - 1
		}
		if !isAlphanum(b[len(b)-1]) {
			return i - 1
		}
		i++
	}
}

// Check if reserved
func (ls *LScanner) isReserved() bool {
	for i := 7; i > 0; i-- {
		chars, _ := ls.reader.Peek(i)
		_, ok := ls.reserved[string(chars)]
		if ok {
			return true
		}
	}
	return false
}

// Returns the reserved words Type
func (ls *LScanner) getReserved() (string, string) {
	chars, _ := ls.reader.Peek(7)
	for i := 7; i > 0; i-- {
		s, ok := ls.reserved[string(chars[:i])]
		if ok {
			return s, string(chars[:i])
		}
	}
	return "", ""
}

func (ls *LScanner) integer() int {
	i := 1
	for {
		b, err := ls.reader.Peek(i)
		if err == io.EOF {
			return i - 1
		}
		if !isDigit(b[len(b)-1]) {
			return i - 1
		}
		i++
	}
}

func (ls *LScanner) isletter() bool {
	b, _ := ls.reader.Peek(1)
	return isLetter(b[0])
}

func isLetter(b byte) bool {
	return (b >= 97 && b <= 122) || (b >= 65 && b <= 90)
}

func (ls *LScanner) isAlphanum() bool {
	b, _ := ls.reader.Peek(1)
	return isAlphanum(b[0])
}

func isAlphanum(b byte) bool {
	return isLetter(b) || isDigit(b) || b == 95
}

func (ls *LScanner) isDigit() bool {
	b, _ := ls.reader.Peek(1)
	return isDigit(b[0])
}

func isDigit(b byte) bool {
	return b >= 48 && b <= 57
}

func (ls *LScanner) isFraction(offset int) (int, bool) {
	minimum, err := ls.reader.Peek(offset + 2)
	if err == io.EOF || string(minimum[offset]) != "." || !isDigit(minimum[offset+1]) {
		return 0, false
	}
	existE := false
	i := 1
	for {
		b, err := ls.reader.Peek(i + offset)
		if err == io.EOF {
			return i - 1 + offset, true
		}
		if string(b[len(b)-1]) == "e" && existE {
			return i - 1 + offset, true
		}
		if string(b[len(b)-1]) == "e" {
			existE = true
			b, err = ls.reader.Peek(i + offset + 1)
			if err == io.EOF {
				return i - 1 + offset, true
			}
			if string(b[len(b)-1]) == "-" {
				b, err = ls.reader.Peek(i + offset + 2)
				if err == io.EOF {
					return i - 1 + offset, true
				}
				if !isNonzero(b[len(b)-1]) {
					return i - 1 + offset, true
				}
			}
			if !isNonzero(b[len(b)-1]) {
				return i - 1 + offset, true
			}
		}
		if !isFraction(string(b[offset:])) {
			return i - 1 + offset, true
		}
		i++
	}
}

func isFraction(b string) bool {
	numE := 0
	eJustHappened := false
	if string(b[0]) != "." {
		return false
	}
	if len(b) == 2 && string(b[0]) == "0" {
		return false
	}
	if len(b) == 2 && string(b[1]) == "0" {
		return true
	}
	if b[len(b)-1] == 0 && b[len(b)-2] == 0 {
		return false
	}
	for i := 1; i < len(b); i++ {
		if !isDigit(b[i]) {
			if numE == 0 && string(b[i]) == "e" {
				numE++
				eJustHappened = true
				continue
			}
			if eJustHappened {
				if string(b[i]) == "-" {
					eJustHappened = false
					continue
				}
			}
			eJustHappened = false
			return false
		}
	}
	return true
}

func isNonzero(b byte) bool {
	return b >= 49 && b <= 57
}

func (ls *LScanner) isNonzero() bool {
	b, _ := ls.reader.Peek(1)
	return isNonzero(b[0])
}

func (ls *LScanner) isEOF() bool {
	_, err := ls.reader.Peek(1)
	return err == io.EOF
}

func (ls *LScanner) isIdent(s string) bool {
	chars, err := ls.reader.Peek(len(s))
	if err != nil {
		return false
	}
	return string(chars) == s
}

func (ls *LScanner) isNext(i int) bool {
	chars, err := ls.reader.Peek(1)
	if err != nil {
		return false
	}
	return chars[0] == byte(i)
}

func (ls *LScanner) read(s string) string {
	chars := make([]byte, len(s))
	ls.reader.Read(chars)
	return string(chars)
}

func (ls *LScanner) token(t, l string, line, col int) *Token {
	ls.col += len(l)
	lexeme := " "
	if l != " " {
		lexeme = ls.read(l)
		fmt.Printf("lexeme: %s\t\tline:%d column:%d\n", lexeme, ls.line, col)
	} else {
		lexeme = ls.read(l)
	}
	return &Token{t, lexeme, fmt.Sprintf("%d %d", line, col)}
}

func (ls *LScanner) NextToken() (*Token, error) {
	switch {
	case ls.isEOF(): // is eof
		return nil, errors.New("EOF")
	case ls.isNext(32): // space
		ls.token("space", " ", ls.line, ls.col)
		return nil, nil
	case ls.isNext(10): // new line
		ls.line++
		ls.col = 0
		ls.token("new line", " ", ls.line, ls.col)
		return nil, nil
	case ls.isNonzero():
		n := ls.integer()
		f, t := ls.isFraction(n)
		if !t {
			str, _ := ls.reader.Peek(n)
			return ls.token("intNum", string(str), ls.line, ls.col), nil
		}
		str, _ := ls.reader.Peek(f)
		return ls.token("floatNum", string(str), ls.line, ls.col), nil
	case ls.isIdent("0"):
		n, t := ls.isFraction(1)
		if !t {
			return ls.token("intNum", "0", ls.line, ls.col), nil
		}
		str, _ := ls.reader.Peek(n)
		return ls.token("floatNum", string(str), ls.line, ls.col), nil
		// ----------------------------- Reserved Words ------------------------------------------------
	case ls.isIdent("_"):
		ls.read("_")
		return ls.token("_", "_", ls.line, ls.col), fmt.Errorf("%d:%d    invalid identifier", ls.line, ls.col)
	case ls.isReserved():
		t, l := ls.getReserved()
		return ls.token(t, l, ls.line, ls.col), nil
	case ls.isletter():
		n := ls.id()
		str, _ := ls.reader.Peek(n)
		return ls.token("id", string(str), ls.line, ls.col), nil
	default:
		badChar, _ := ls.reader.Peek(1)
		return ls.token(" ", string(badChar), ls.line, ls.col), fmt.Errorf("%d:%d %s is not an accepted character", ls.line, ls.col, badChar)
		// ----------------------------- Reserved Words ------------------------------------------------
	}
	return nil, nil
}
