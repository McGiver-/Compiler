package Lex

import(
	"bufio"
	"os"
	"fmt"
	"errors"
	"io"
)

type Token struct {
	Type string
	Lexeme string
	location string
}

type LScanner struct{
	reader *bufio.Reader
	col int
	line int
}

func CreateScanner(fileName string) (*LScanner, error){
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	return &LScanner{reader,0,0}, nil
}

func (ls *LScanner) id() int {
	i := 1
	for{
		b, err := ls.reader.Peek(i)
		if err == io.EOF {
			return i-1
		}
		if !isAlphanum(b[len(b)-1]) {
			return i-1
		}
		i++
	}
}

func (ls *LScanner) integer() int {
	i := 1
	for{
		b, err := ls.reader.Peek(i)
		if err == io.EOF {
			return i-1
		}
		if !isDigit(b[len(b)-1]) {
			return i-1
		}
		i++
	}
}

func (ls *LScanner) isletter() bool {
	b, _ := ls.reader.Peek(1)
	return (b[0] >= 97 && b[0] <= 122) || (b[0] >= 65 && b[0] <= 90)
}

func isLetter(b byte) bool {
	return (b >= 97 && b <= 122) || (b >= 65 && b <= 90)
}

func (ls *LScanner) isAlphanum() bool {
	b, _ := ls.reader.Peek(1)
	return ls.isletter() || ls.isDigit() || b[0] == 95
}

func isAlphanum(b byte) bool {
	return isLetter(b) || isDigit(b) || b == 95
}

func (ls *LScanner) isDigit() bool {
	b, _ := ls.reader.Peek(1)
	return b[0] >= 48 && b[0] <= 57
}

func isDigit(b byte) bool {
	return b >= 48 && b <= 57
}

func (ls *LScanner) isFraction(offset int) (int,bool) {
	minimum,err := ls.reader.Peek(offset+2)
	if err == io.EOF || string(minimum[offset]) != "." || !isDigit(minimum[offset+1]) {
		return 0, false;
	}
	existE := false
	i := 1
	for{
		b, err := ls.reader.Peek(i+offset)
		if err == io.EOF {
			return i-1+offset , true
		}
		if string(b[len(b)-1]) == "e" && existE {
			return i-1+offset , true
		}
		if  string(b[len(b)-1]) == "e"{
			existE = true
			b, err = ls.reader.Peek(i+offset+1)
			if err == io.EOF {
				return i-1+offset , true
			}
			if string(b[len(b)-1]) == "-"{
				b, err = ls.reader.Peek(i+offset+2)
				if err == io.EOF {
					return i-1+offset , true
				}
				if !isNonzero(b[len(b)-1]){
					return i-1+offset, true
				}
			}
			if !isNonzero(b[len(b)-1]) {
					return i-1+offset, true
			}
		}
		if !isFraction(string(b[offset:])) {
			return i-1+offset, true
		}
		i++
	}
}

func  isFraction(b string) bool {
	numE := 0;
	eJustHappened := false
	if string(b[0]) != "."{
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
	for i := 1; i< len(b); i++ {
		if !isDigit(b[i]){
			if numE == 0 && string(b[i]) == "e"{
				numE++
				eJustHappened = true
				continue
			}
			if eJustHappened {
				if string(b[i]) == "-"{
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
	return b[0] >= 49 && b[0] <= 57
}

func (ls *LScanner) isEOF() bool{
	_, err := ls.reader.Peek(1)
	return err == io.EOF
}

func (ls *LScanner) isIdent(s string) bool{
	chars, err := ls.reader.Peek(len(s))
	if err != nil {
		return false
	}
	return string(chars) == s
}

func (ls *LScanner)  isNext(i int) bool {
	chars, err := ls.reader.Peek(1)
	if err != nil {
		return false
	}
	return chars[0] == byte(i)
}

func (ls *LScanner) read(s string) string{
	chars := make([]byte,len(s))
	ls.reader.Read(chars)
	return string(chars)
}

func (ls *LScanner) readN(i int) string{
	chars := make([]byte,i)
	ls.reader.Read(chars)
	return string(chars)
}

func (ls *LScanner )token(t,l string, line,col int) *Token{
	ls.col += len(l)
	return &Token{t,l,fmt.Sprintf("%d %d",line,col)} 
}

func (ls *LScanner) NextToken() (*Token,error){
	switch {
	case ls.isEOF(): // is eof
		return nil, errors.New("EOF")
	case ls.isNext(32): // space
		ls.read(" ")
		ls.col++
		return nil,nil
	case ls.isNext(10): // space
		ls.read(" ")
		ls.line++
		ls.col = 0
		return nil,nil
	case ls.isNonzero():
		n := ls.integer()
		l,t := ls.isFraction(n)
		if !t{
			return ls.token("integer",ls.readN(n),ls.line,ls.col),nil
		}
		return ls.token("float",ls.readN(l),ls.line,ls.col),nil
	case ls.isIdent("0"):
		n, t := ls.isFraction(1)
		if !t{
			return ls.token("integer",ls.read("0"),ls.line,ls.col), nil
		}
		return ls.token("float",ls.readN(n),ls.line,ls.col),nil
		// ----------------------------- Reserved Words ------------------------------------------------
	case ls.isIdent("_"):
		ls.read("_")
		return ls.token("_","_",ls.line,ls.col), fmt.Errorf("%d:%d    invalid identifier",ls.line,ls.col)
	case ls.isIdent("program"):
		ls.read("program")
		return ls.token("program","program",ls.line,ls.col), nil
	case ls.isIdent("return"):
		ls.read("return")
		return ls.token("return","return",ls.line,ls.col), nil
	case ls.isIdent("put"):
		ls.read("put")
		return ls.token("put","put",ls.line,ls.col), nil
	case ls.isIdent("get"):
		ls.read("get")
		return ls.token("get","get",ls.line,ls.col), nil
	case ls.isIdent("float"):
		ls.read("float")
		return ls.token("float","float",ls.line,ls.col), nil
	case ls.isIdent("int"):
		ls.read("int")
		return ls.token("int","int",ls.line,ls.col), nil
	case ls.isIdent("class"):
		ls.read("class")
		return ls.token("class","class",ls.line,ls.col), nil
	case ls.isIdent("for"):
		ls.read("for")
		return ls.token("for","for",ls.line,ls.col), nil
	case ls.isIdent("else"):
		ls.read("else")
		return ls.token("else","else",ls.line,ls.col), nil
	case ls.isIdent("then"):
		ls.read("then")
		return ls.token("then","then",ls.line,ls.col), nil
	case ls.isIdent("if"):
		ls.read("if")
		return ls.token("if","if",ls.line,ls.col), nil
	case ls.isIdent("("):
		ls.read("(")
		return ls.token("openParen","(",ls.line,ls.col), nil
	case ls.isIdent(")"):
		ls.read(")")
		return ls.token("closeParen",")",ls.line,ls.col), nil
	case ls.isIdent("{"):
		ls.read("{")
		return ls.token("openCurly{","{",ls.line,ls.col), nil
	case ls.isIdent("}"):
		ls.read("}")
		return ls.token("closeCurly","}",ls.line,ls.col), nil
	case ls.isIdent("["):
		ls.read("[")
		return ls.token("openSquare","[",ls.line,ls.col), nil
	case ls.isIdent("]"):
		ls.read("]")
		return ls.token("closeSquare","]",ls.line,ls.col), nil
	case ls.isIdent("not"):
		ls.read("not")
		return ls.token("not","not",ls.line,ls.col), nil
	case ls.isIdent("and"):
		ls.read("and")
		return ls.token("and","and",ls.line,ls.col), nil
	case ls.isIdent("or"):
		ls.read("or")
		return ls.token("or","or",ls.line,ls.col), nil
	case ls.isIdent(";"):
		ls.read(";")
		return ls.token(";",";",ls.line,ls.col), nil
	case ls.isIdent(","):
		ls.read(",")
		return ls.token(",",",",ls.line,ls.col), nil
	case ls.isIdent("."):
		ls.read(".")
		return ls.token(".",".",ls.line,ls.col), nil
	case ls.isIdent("::"):
		ls.read("::")
		return ls.token("::","::",ls.line,ls.col), nil
	case ls.isIdent(":"):
		ls.read(":")
		return ls.token(":",":",ls.line,ls.col), nil
	case ls.isIdent("//"):
		ls.read("//")
		return ls.token("//","//",ls.line,ls.col), nil
	case ls.isIdent("*/"):
		ls.read("*/")
		return ls.token("star/","*/",ls.line,ls.col), nil
	case ls.isIdent("*"):
		ls.read("*")
		return ls.token("star","*",ls.line,ls.col), nil
	case ls.isIdent("/*"):
		ls.read("/*")
		return ls.token("/star","/*",ls.line,ls.col), nil
	case ls.isIdent("/"):
		ls.read("/")
		return ls.token("/","/",ls.line,ls.col), nil
	case ls.isIdent("<>"):
		ls.read("<>")
		return ls.token("<>","<>",ls.line,ls.col), nil
	case ls.isIdent("<="):
		ls.read("<=")
		return ls.token("<=","<=",ls.line,ls.col), nil
	case ls.isIdent("<"):
		ls.read("<")
		return ls.token("<","<",ls.line,ls.col), nil
	case ls.isIdent(">="):
		ls.read(">=")
		return ls.token(">=",">=",ls.line,ls.col), nil
	case ls.isIdent(">"):
		ls.read(">")
		return ls.token(">",">",ls.line,ls.col), nil
	case ls.isIdent("=="):
		ls.read("==")
		return ls.token("==","==",ls.line,ls.col), nil
	case ls.isIdent("="):
		ls.read("=")
		return ls.token("=","=",ls.line,ls.col), nil
	case ls.isIdent("+"):
		ls.read("+")
		return ls.token("plus","+",ls.line,ls.col), nil
	case ls.isIdent("-"):
		ls.read("-")
		return ls.token("minus","-",ls.line,ls.col), nil
	case ls.isletter():
		n := ls.id()
		id := ls.readN(n)
		return ls.token("id",id,ls.line,ls.col),nil
	default:
		return ls.token(" ",ls.read(" "),ls.line,ls.col), fmt.Errorf("%d:%d    not an accepted character",ls.line,ls.col)
		// ----------------------------- Reserved Words ------------------------------------------------
	}
	return nil,nil
}


