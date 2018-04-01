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

func (syn *SynAnalyzer) Parse() (errorList []error, rootNode *Node) {
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
	sStack.Pop()
	rootNode = sStack.Pop().(*Node)
	return
}

func handleSemanticAction(action string, token *Lex.Token, stack *stackgo.Stack) {
	if action == "EPSILON" {
		stack.Push(makeNode("EPSILON", "EPSILON", nil))
	}
	switch action {
	case "id", "float", "int":
		node := makeNode(action, action, token)
		stack.Push(node)
	}

	if action == "Type" {
		t := stack.Pop().(*Node)
		stack.Push(makeFamily("Type", "Type", t.Token, t))
	}

	if action == "ClassListAndFuncDefList" {
		if stack.Top().(*Node).Type == "FuncDef" {
			funcDefMember := stack.Pop().(*Node)
			if stack.Top() != nil && stack.Top().(*Node).Type == "ClassMember" {
				prog := makeFamily("Prog", "Prog", token, stack.Pop().(*Node), funcDefMember)
				prog.Token = prog.LeftMostChild.Token
				stack.Push(prog)
			} else {
				prog := makeFamily("Prog", "Prog", token, funcDefMember)
				prog.Token = prog.LeftMostChild.Token
				stack.Push(prog)
			}
		} else if stack.Top().(*Node).Type == "ClassMember" {
			prog := makeFamily("Prog", "Prog", token, stack.Pop().(*Node))
			prog.Token = prog.LeftMostChild.Token
			stack.Push(prog)
		}

	}

	if action == "InheritList" {
		top := stack.Top().(*Node)
		if top.Value == "EPSILON" {
			stack.Pop()
			if stack.Top().(*Node).Type == "InheritListMember" {
				inheritList := makeFamily("InheritList", "InheritList", token, stack.Pop().(*Node))
				inheritList.Token = inheritList.LeftMostChild.Token
				stack.Push(inheritList)
			} else {
				stack.Push(makeNode("InheritList", "EPSILON", nil))
			}
		}
	}

	if action == "InheritListMember" {
		id := stack.Pop().(*Node)
		inheritListMember := makeFamily("InheritListMember", "InheritListMember", id.Token, id)
		top := stack.Top().(*Node)
		if top.Value == "InheritListMember" {
			top.makeSiblings(inheritListMember)
		} else {
			stack.Push(inheritListMember)
		}
	}

	if action == "ClassMember" {
		memberList := stack.Pop().(*Node)
		inheritList := stack.Pop().(*Node)
		id := stack.Pop().(*Node)
		classMember := makeFamily("ClassMember", "ClassMember", token, id, inheritList, memberList)
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
	/// funcDef
	if action == "FuncDef" {
		if stack.Top().(*Node).Type == "EPSILON" {
			stack.Pop() // Poping EPSILON
		}
		statBlock := stack.Pop().(*Node)
		fParamList := stack.Pop().(*Node)
		id := stack.Pop().(*Node)
		scopeSpec := stack.Pop().(*Node)
		typeNode := stack.Pop().(*Node)
		funcDef := makeFamily("FuncDef", "FuncDef", token, typeNode, scopeSpec, id, fParamList, statBlock)
		funcDef.Token = funcDef.LeftMostChild.Token
		top := stack.Top()
		if top == nil || top.(*Node).Type == "ClassMember" {
			stack.Push(funcDef)
			return
		}
		node := top.(*Node)
		if node.Type == "FuncDef" {
			node.makeSiblings(funcDef)
		}
	}
	/// EmptyScopeSpec
	if action == "EmptyScope" {
		id := stack.Pop().(*Node)
		stack.Push(makeNode("Scope", "EPSILON", token))
		stack.Push(id)
	}
	///

	/// ScopeSpec
	if action == "ScopeSpec" {
		scope := makeFamily("Scope", "Scope", token, stack.Pop().(*Node))
		scope.Token = scope.LeftMostChild.Token
		stack.Push(scope)
	}
	/// ScopeSpec
	if action == "MemberList" {
		top := stack.Top().(*Node)
		if top.Value == "EPSILON" {
			stack.Pop()
			if stack.Top().(*Node).Type == "FuncDecl" || stack.Top().(*Node).Type == "VarDecl" {
				memberList := makeFamily("MemberList", "MemberList", token, stack.Pop().(*Node))
				memberList.Token = memberList.LeftMostChild.Token
				stack.Push(memberList)
			} else {
				stack.Push(makeNode("MemberList", "EPSILON", nil))
			}
		}
	}

	//StatBlock
	if action == "StatBlock" {
		top := stack.Top().(*Node)
		if top.Value == "EPSILON" {
			stack.Pop()
			if stack.Top().(*Node).Type == "VarDecl" || stack.Top().(*Node).Type == "Stat" {
				statBlock := makeFamily("StatBlock", "StatBlock", token, stack.Pop().(*Node))
				statBlock.Token = statBlock.LeftMostChild.Token
				stack.Push(statBlock)
			} else {
				stack.Push(makeNode("StatBlock", "EPSILON", nil))
			}
		} else {
			statBlock := makeFamily("StatBlock", "StatBlock", token, stack.Pop().(*Node))
			statBlock.Token = statBlock.LeftMostChild.Token
			stack.Push(statBlock)
		}
	}
	//

	if action == "VarDecl" {
		dimList := stack.Pop().(*Node)
		id := stack.Pop().(*Node)
		typeNode := stack.Pop().(*Node)
		varDecl := makeFamily("VarDecl", "VarDecl", typeNode.Token, typeNode, id, dimList)
		top := stack.Top().(*Node)
		if top.Type == "FuncDecl" || top.Type == "VarDecl" || top.Type == "Stat" {
			top.makeSiblings(varDecl)
		} else {
			stack.Push(varDecl)
		}
	}
	if action == "FuncDecl" {
		fparamList := stack.Pop().(*Node)
		id := stack.Pop().(*Node)
		typeNode := stack.Pop().(*Node)
		funcDecl := makeFamily("FuncDecl", "FuncDecl", typeNode.Token, typeNode, id, fparamList)
		top := stack.Top().(*Node)
		if top.Type == "FuncDecl" || top.Type == "VarDecl" {
			top.makeSiblings(funcDecl)
		} else {
			stack.Push(funcDecl)
		}
	}

	if action == "FparamList" {
		top := stack.Top().(*Node)
		if top.Value == "EPSILON" {
			stack.Pop()
			if stack.Top().(*Node).Type == "FparamMember" {
				fparamList := makeFamily("FparamList", "FparamList", token, stack.Pop().(*Node))
				fparamList.Token = fparamList.LeftMostChild.Token
				stack.Push(fparamList)
			} else {
				stack.Push(makeNode("FparamList", "EPSILON", nil))
			}
		}
	}
	if action == "FparamMember" {
		dimList := stack.Pop().(*Node)
		id := stack.Pop().(*Node)
		typeNode := stack.Pop().(*Node)
		fparamMember := makeFamily("FparamMember", "FparaMember", token, typeNode, id, dimList)
		fparamMember.Token = fparamMember.LeftMostChild.Token
		top := stack.Top().(*Node)
		if top.Type == "FparamMember" {
			top.makeSiblings(fparamMember)
		} else {
			stack.Push(fparamMember)
		}
	}

	if action == "DimList" {
		top := stack.Top().(*Node)
		if top.Value == "EPSILON" {
			stack.Pop()
			if stack.Top().(*Node).Type == "intNum" {
				dimList := makeFamily("DimList", "DimList", token, stack.Pop().(*Node))
				dimList.Token = dimList.LeftMostChild.Token
				stack.Push(dimList)
			} else {
				stack.Push(makeNode("DimList", "EPSILON", nil))
			}
		}
	}
	if action == "intNum" {
		intNum := makeNode("intNum", token.Lexeme, token)
		top := stack.Top().(*Node)
		if top.Type == "intNum" {
			top.makeSiblings(intNum)
		} else {
			stack.Push(intNum)
		}
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
