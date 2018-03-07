package Syn

import (
	"fmt"
	"strings"

	"github.com/McGiver-/Compiler/Lex"
	"github.com/alediaferia/stackgo"
)

type SynAnalyzer struct {
	tokens           []*Lex.Token
	tokenIndex       int
	attributeGrammar [96]string
	predictionTable  map[string]map[string]int
	parsingStack     *stackgo.Stack
	semanticStack    *stackgo.Stack
}

// CreateAnalyzer creates the Analyzer by passing the token array
func CreateAnalyzer(tokens []*Lex.Token) (*SynAnalyzer, error) {

	return &SynAnalyzer{
		tokens,
		-1,
		attributeGrammar,
		predictionTable,
		stackgo.NewStack(),
		stackgo.NewStack(),
	}, nil
}

func (syn *SynAnalyzer) Parse() (errorList []error) {
	skipping := false
	pStack := syn.parsingStack // This is the stack that holds the parsing symbols that are pushed.
	var symbol string          // Symbol variable

	pStack.Push("$")
	printStack(pStack)
	pStack.Push("Prog")         // Nonterminal that starts the progra
	token, _ := syn.nextToken() // First token before loop
	printStack(pStack)
	for pStack.Top() != "$" { // Loop while not at the end of the program yet
		printStack(pStack)
		if !skipping {
			symbol = pStack.Top().(string) // Peek at the top symbol
		}
		if symbol == "EPSILON" {
			pStack.Pop()
			continue
		}
		fmt.Printf("symbol %s\n", symbol)
		if _, ok := terminals[symbol]; ok { // Enter if the symbol is a terminal
			if symbol == token.Type { // The symbol at the top of the stack matches the read token
				fmt.Printf("matched %s\n", token.Type)
				pStack.Pop()
				printStack(pStack)
				token, _ = syn.nextToken()
			} else {
				skipping = true
				errorList = append(errorList, fmt.Errorf("Expected %s at %s", symbol, token.Location))
				token, err := syn.nextToken() // take the next token
				if err != nil {
					return
				}
				for symbol != token.Type { // if the terminal does not match the token than keep looping until it does
					token, err = syn.nextToken()
					if err != nil {
						return
					}
				}
			}
		} else { // If the symbol was not a terminal i.e nonterminal
			rhsList, err := syn.getProduction(symbol, token.Type)
			fmt.Printf("symbol %s token %s lexeme %s rhsList %v position %s\n", symbol, token.Type, token.Lexeme, rhsList, token.Location)
			if err != nil {
				err = fmt.Errorf("error %v at %s token %s lexeme %s symbol %s", err, token.Location, token.Type, token.Lexeme, symbol)
				errorList = append(errorList, err)
				skipping = true
				token, err := syn.nextToken()
				if err != nil {
					return
				}
				for {
					for _, v := range rhsList {
						if token.Type == v {
							goto endLoop
						}
						token, err = syn.nextToken()
						if err != nil {
							return
						}
					}

				}
			endLoop:
			} else {
				pStack.Pop()
				inverseRhsMultiplePush(pStack, rhsList)
			}
		}
	}
	return
}

func printStack(stack *stackgo.Stack) {
	newPrint := fmt.Sprintf("%s", stack.String())
	fmt.Printf("STACK : %s\n", newPrint)
}

func inverseRhsMultiplePush(stack *stackgo.Stack, rhsList []string) {
	for _, v := range reverse(rhsList) {
		stack.Push(v)
	}
}

func (syn *SynAnalyzer) nextToken() (*Lex.Token, error) {
	tokens := syn.tokens
	if syn.tokenIndex >= len(tokens)-1 {
		return nil, fmt.Errorf("Reached the end of the program")
	}
	syn.tokenIndex++
	return tokens[syn.tokenIndex], nil
}

func (syn *SynAnalyzer) getProduction(nonterminal, terminal string) ([]string, error) {
	productionNumber, expected, err := getProductionNumber(syn.predictionTable, nonterminal, terminal)
	if err != nil {
		return expected, err
	}
	return strings.Split(syn.attributeGrammar[productionNumber], " "), err
}

func getProductionNumber(predictionTable map[string]map[string]int, nonterminal, terminal string) (int, []string, error) {
	nonTerminalMap := predictionTable[nonterminal]
	productionNum, ok := nonTerminalMap[terminal]
	if ok {
		return productionNum, []string{}, nil
	}
	expected := ""
	expectedList := make([]string, 1)
	for v := range nonTerminalMap {
		expected += " " + v
		expectedList = append(expectedList, v)
	}
	return 0, expectedList, fmt.Errorf("Expecting one of the following: %s", expected)

}

// ---------------------------------------------------------------------- HELPERS --------------------------------------------------------------

func reverse(list []string) []string {
	for i := 0; i < len(list)/2; i++ {
		j := len(list) - i - 1
		list[i], list[j] = list[j], list[i]
	}
	return list
}
