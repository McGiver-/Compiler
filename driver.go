package main

import (
	"fmt"
	"log"
	"os"

	"github.com/McGiver-/Compiler/Lex"
	"github.com/McGiver-/Compiler/Sem"
	"github.com/McGiver-/Compiler/Syn"
	graph "github.com/awalterschulze/gographviz"
)

func main() {
	tokens := []*Lex.Token{}
	errs := []error{}
	inputFile := os.Args[1]
	// FileA2CC := os.Args[2]
	// ErrorFile := os.Args[3]
	scanner, err := Lex.CreateScanner(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	for {
		tkn, err := scanner.NextToken()
		if err != nil && err.Error() == "EOF" {
			break
		}
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if tkn == nil {
			continue
		}
		tokens = append(tokens, tkn)
	}
	types := ""
	for _, v := range tokens {
		types += v.Type + " "
	}
	// ioutil.WriteFile(FileA2CC, []byte(types), os.ModePerm)

	errors := ""
	for _, v := range errs {
		errors += v.Error() + "\n"
		fmt.Printf("%v \n", v)
	}
	// ioutil.WriteFile(ErrorFile, []byte(errors), os.ModePerm)

	analyzer, err := Syn.CreateAnalyzer(tokens)

	if err != nil {
		fmt.Print(err)
	}

	errorList, rootNode := analyzer.Parse()
	for _, v := range errorList {
		fmt.Println(v)
	}

	tableAnalyzer := Sem.CreateAnalyzer(rootNode)
	errs = tableAnalyzer.CreateTables()
	for _, v := range errs {
		fmt.Println(v)
	}
	g := graph.NewGraph()
	if err := g.SetName("ast"); err != nil {
		panic(err)
	}
	if err := g.SetDir(true); err != nil {
		panic(err)
	}

	g.AddNode("ast", rootNode.Value, nil)
	makeGraph(rootNode, g)
	fmt.Print(g.String())

}

func makeGraph(n *Syn.Node, g *graph.Graph) {
	if n.LeftMostChild == nil {
		return
	}
	child := n.LeftMostChild
	for child != nil {
		g.AddNode("ast", child.Value, nil)
		g.AddEdge(n.Value, child.Value, true, nil)
		makeGraph(child, g)
		child = child.RightSibling
	}
}
