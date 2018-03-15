package Syn

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/McGiver-/Compiler/Lex"
	"github.com/alediaferia/stackgo"
)

type SynAnalyzer struct {
	tokens           []*Lex.Token
	tokenIndex       int
	attributeGrammar [93]string
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
	sStack := syn.semanticStack
	var symbol string // Symbol variable
	var previousToken *Lex.Token

	pStack.Push("$")
	pStack.Push("Prog")         // Nonterminal that starts the progra
	token, _ := syn.nextToken() // First token before loop
	for pStack.Top() != "$" {   // Loop while not at the end of the program yet
		if !skipping {
			symbol = pStack.Top().(string) // Peek at the top symbol
		}
		if symbol == "EPSILON" {
			pStack.Pop()
			continue
		}
		if string(symbol[0]) == "@" {
			pStack.Pop()
			handleSemanticAction(string(symbol[1:]), previousToken, sStack)
			continue
		}
		if _, ok := terminals[symbol]; ok { // Enter if the symbol is a terminal
			if symbol == token.Type { // The symbol at the top of the stack matches the read token
				pStack.Pop()
				previousToken = token
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
	classListNode := sStack.Pop().(*Node)
	fmt.Println(spew.Sdump(classListNode))
	// fmt.Printf("classMemberInheritId %v\n", classListNode.LeftMostChild.LeftMostChild.RightSibling.LeftMostChild.Token)
	return
}

func handleSemanticAction(action string, token *Lex.Token, stack *stackgo.Stack) {
	if action == "id" {
		id := makeNode("id", "id", token)
		stack.Push(id)
	}
	if action == "InheritListMember" {
		id := stack.Pop().(*Node)
		inheritListMember := makeNode("InheritListMember", "InheritListMember", id.Token)
		inheritListMember.adoptChildren(id)
		top := stack.Top().(*Node)
		if top.Type == "InheritListMember" {
			top.makeSiblings(inheritListMember)
		} else {
			stack.Push(inheritListMember)
		}
	}
	if action == "InheritList" {
		inheritList := makeNode("InheritList", "InheritList", token)
		top := stack.Top().(*Node)
		if top.Type == "InheritListMember" {
			stack.Pop()
			inheritList.adoptChildren(top)
			inheritList.Token = inheritList.LeftMostChild.Token
			stack.Push(inheritList)
		} else {
			stack.Push(makeNode("InheritList", "EPSILON", nil))
		}
	}
	if action == "ClassMember" {
		classMember := makeNode("ClassMember", "ClassMember", token)
		inheritList := stack.Pop().(*Node)
		id := stack.Pop().(*Node)
		classMember.adoptChildren(id.makeSiblings(inheritList))
		classMember.Token = classMember.LeftMostChild.Token
		top := stack.Top()
		if top == nil {
			stack.Push(classMember)
			return
		}
		node := top.(*Node)
		if node.Type == "ClassMember" {
			node.makeSiblings(classMember)
		}
	}
	if action == "ClassList" {
		classList := makeNode("ClassList", "ClassList", token)
		classList.adoptChildren(stack.Pop().(*Node))
		classList.Token = classList.LeftMostChild.Token
		stack.Push(classList)
	}
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
