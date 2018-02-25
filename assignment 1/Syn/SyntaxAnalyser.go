package Syn

import (
	"fmt"
	"strings"

	"github.com/McGiver-/Compiler/Lex"
	"github.com/alediaferia/stackgo"
)

type SynAnalyzer struct {
	tokens           *[]Lex.Token
	parseTable       map[string]string
	parsingStack     *stackgo.Stack
	semanticStack    *stackgo.Stack
	current          *Lex.Token
	attributeGrammar [96]string
	predictionTable  map[string]map[string]int
}

// This is the production list. There are 96 productions. Number corresponds to what should be pushed onto the stack.
// This push should be done in reverse.
var attributeGrammar = [96]string{
	"ProgBody",
	"ClassDecl Prog",
	"FuncDef Prog",
	"program FuncBody ; FuncDefList",
	"class id Inherit { ClassVarDecl } ;",
	": id InheritList",
	"EPSILON",
	", id InheritList",
	"EPSILON",
	"Decl ClassVarDeclTail",
	"EPSILON",
	"Type id",
	"ArraySizeTail ; ClassVarDecl",
	"FuncDefListTail",
	"FuncHead FuncBody ;",
	"Decl FuncDefListTail",
	"EPSILON",
	"Type id FuncHeadTail",
	"sr id ( FParams )",
	"( FParams )",
	"EPSILON",
	"FuncDef FuncDefList",
	"{ FuncVarDecl }",
	"StatementNoAssign StatementTail",
	"Type FuncVarDeclTail",
	"EPSILON",
	"id ArraySizeTail ; FuncVarDecl",
	"VariableTail AssignStatTail ; StatementTail",
	"Variable AssignStatTail",
	"AssignOp expr",
	"AssignStat ;",
	"StatementNoAssign",
	"Statement StatementTail",
	"EPSILON",
	"for ( Decl AssignOp Expr ; RelExpr ; AssignStat ) StatBlock ;",
	"{ StatementTail }",
	"Statement",
	"EPSILON",
	"ArithExpr RelExprTail",
	"ArithExpr RelOp ArithExpr",
	"RelOp ArithExpr",
	"EPSILON",
	"Term ArithExprTail",
	"AddOp Term ArithExprTail",
	"EPSILON",
	"Factor TermTail",
	"MultOp Factor TermTail",
	"EPSILON",
	"EvalExprHead",
	"intNum",
	"floatNum",
	"+",
	"-",
	"id EvalExprTail",
	"EvalIndiceTail",
	"FunctionCallExpr",
	"Indice EvalIndiceTail",
	"EvalNestHead",
	"EvalIndiceHead",
	"EPSILON",
	". EvalExprHead",
	"id VariableTail",
	"VarIndiceTail",
	"VarNest",
	"Indice VariableTail",
	"EPSILON",
	". Variable",
	"( AParams )",
	"[ ArithExpr ]",
	"[ intNum ]",
	"ArraySize ArraySizeTail",
	"EPSILON",
	"float",
	"id",
	"int",
	"Type id ArraySizeTail FParamsTail",
	"EPSILON",
	", FParams",
	"EPSILON",
	"Expr AParamsTail",
	"EPSILON",
	", AParams",
	"EPSILON",
	"=",
	"eq",
	"neq",
	"lt",
	"gt",
	"leq",
	"geq",
	"+",
	"-",
	"or",
	"*",
	"/",
	"and",
}

// This table has as key the Nonterminal. This is a key to another map which has as a key
// the expected terminals and as value the number of the production that should be pushed to the stack.
// These numbers correspond to the productions is the attributeGrammar array.
var predictionTable = map[string]map[string]int{
	"Prog":              map[string]int{"program": 0, "class": 1, "float": 2, "id": 2, "int": 2},
	"ProgBody":          map[string]int{"program": 3},
	"ClassDecl":         map[string]int{"class": 4},
	"Inherit":           map[string]int{":": 5, "{": 6},
	"InheritList":       map[string]int{",": 7, "{": 8},
	"ClassVarDecl":      map[string]int{"float": 9, "id": 9, "int": 9, "}": 10},
	"Decl":              map[string]int{"float": 11, "id": 11, "int": 11},
	"ClassVarDeclTail":  map[string]int{"[": 12, ";": 12, "float": 13, "id": 13, "int": 13},
	"FuncDef":           map[string]int{"float": 14, "id": 14, "int": 14},
	"FuncDefList":       map[string]int{"float": 15, "id": 15, "int": 15, "$": 16, "}": 16},
	"FuncHead":          map[string]int{"float": 17, "id": 17, "int": 17},
	"FuncHeadTail":      map[string]int{"sr": 18, "{": 20, "(": 19},
	"FuncDefListTail":   map[string]int{"float": 21, "id": 21, "int": 21},
	"FuncBody":          map[string]int{"{": 22},
	"FuncVarDecl":       map[string]int{"for": 23, "float": 24, "id": 24, "int": 24, "}": 25},
	"FuncVarDeclTail":   map[string]int{"id": 26, "[": 27, ".": 27, "=": 27},
	"AssignStat":        map[string]int{"id": 28},
	"AssignStatTail":    map[string]int{"=": 29},
	"Statement":         map[string]int{"id": 30, "for": 31},
	"StatementTail":     map[string]int{"id": 32, "for": 32, "}": 33},
	"StatementNoAssign": map[string]int{"for": 34},
	"StatBlock":         map[string]int{"{": 35, "id": 36, "for": 36, ";": 37},
	"Expr":              map[string]int{"id": 38},
	"RelExpr":           map[string]int{"id": 39},
	"RelExprTail":       map[string]int{"eq": 40, "neq": 40, "lt": 40, "gt": 40, "leq": 40, "geq": 40, ",": 41, ";": 41, ")": 41},
	"ArithExpr":         map[string]int{"id": 42},
	"ArithExprTail":     map[string]int{"+": 43, "-": 43, "or": 43, "]": 44, "eq": 44, "neq": 44, "lt": 44, "gt": 44, "leq": 44, "geq": 44, ",": 44, ";": 44, ")": 44},
	"Term":              map[string]int{"id": 45},
	"TermTail":          map[string]int{"*": 46, "/": 46, "and": 46, "+": 47, "-": 47, "or": 47, "[": 47, "eq": 47, "neq": 47, "lt": 47, "gt": 47, "leq": 47, "geq": 47, ",": 47, ";": 47, ")": 47},
	"Factor":            map[string]int{"id": 48},
	"Num":               map[string]int{"intNum": 49, "floatNum": 50},
	"Sign":              map[string]int{"+": 51, "-": 52},
	"EvalExprHead":      map[string]int{"id": 53},
	"EvalExprTail":      map[string]int{".": 54, "]": 54, "(": 55},
	"EvalIndiceHead":    map[string]int{"[": 56},
	"EvalIndiceTail":    map[string]int{".": 57, "[": 58, "*": 59, "/": 59, "and": 59, "+": 59, "-": 59, "or": 59, "]": 59, "eq": 59, "neq": 59, "lt": 59, "gt": 59, "leq": 59, "geq": 59, ",": 59, ";": 59, ")": 59},
	"EvalNestHead":      map[string]int{".": 60},
	"Variable":          map[string]int{"id": 61},
	"VariableTail":      map[string]int{"[": 62, ".": 63},
	"VarIndiceTail":     map[string]int{"[": 64, "=": 65},
	"VarNest":           map[string]int{".": 66},
	"FunctionCallExpr":  map[string]int{"(": 67},
	"Indice":            map[string]int{"[": 68},
	"ArraySize":         map[string]int{"[": 69},
	"ArraySizeTail":     map[string]int{"[": 70, ",": 71, ";": 71, ")": 71},
	"Type":              map[string]int{"float": 72, "id": 73, "int": 74},
	"FParams":           map[string]int{"float": 75, "id": 75, "int": 75, ")": 76},
	"FParamsTail":       map[string]int{",": 77, ")": 78},
	"AParams":           map[string]int{"id": 79, ")": 80},
	"AParamsTail":       map[string]int{".": 81, ")": 82},
	"AssignOp":          map[string]int{"=": 83},
	"RelOp":             map[string]int{"eq": 84, "neq": 85, "lt": 86, "gt": 87, "leq": 88, "geq": 89},
	"AddOp":             map[string]int{"+": 90, "-": 91, "or": 92},
	"MultOp":            map[string]int{"*": 93, "/": 94, "and": 95},
}

func (syn *SynAnalyzer) consumeToken() error {
}

func getProduction(predictionTable map[string]map[string]int, attributeGrammar [96]string, nonterminal, terminal string) ([]string, error) {
	productionNumber, err := getProductionNumber(predictionTable, nonterminal, terminal)
	return strings.Split(attributeGrammar[productionNumber], " "), err
}

func getProductionNumber(predictionTable map[string]map[string]int, nonterminal, terminal string) (int, error) {
	nonTerminalMap := predictionTable[nonterminal]
	productionNum, ok := nonTerminalMap[terminal]
	if ok {
		return productionNum, nil
	}
	expected := ""
	for v := range nonTerminalMap {
		expected += " " + v
	}
	return 0, fmt.Errorf("Expecting one of the following: %s", expected)

}

// CreateAnalyzer creates the Analyzer by passing the token array
func CreateAnalyzer(tokens []Lex.Token) (*SynAnalyzer, error) {

	return &SynAnalyzer{
		&tokens,
		nil,
		stackgo.NewStack(),
		stackgo.NewStack(),
		&tokens[0],
		attributeGrammar,
		predictionTable,
	}, nil
}
