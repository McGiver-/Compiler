package main

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/McGiver-/Compiler/refactor/Lex"
	"github.com/McGiver-/Compiler/refactor/Syn"
	graph "github.com/awalterschulze/gographviz"
)

func main() {
	inputFile := os.Args[1]
	r, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("could not open file")
	}
	lexer, err := Lex.CreateLexer(r)
	if err != nil {
		log.Fatal(err)
	}
	tc, errs := lexer.GetTokensNoChan()
	for _, err := range errs {
		fmt.Printf("%v\n", err)
	}
	analyzer, err := Syn.CreateAnalyzer(tc)
	if err != nil {
		fmt.Print(err)
	}

	ec, rootNode := analyzer.Parse()
	fmt.Printf("%s", spew.Sdump(rootNode))
	errs = append(errs, ec...)
	for _, err := range errs {
		fmt.Printf("%v\n", err)
	}
	g := graph.NewGraph()
	if err := g.SetName("ast"); err != nil {
		panic(err)
	}
	if err := g.SetDir(true); err != nil {
		panic(err)
	}

	g.AddNode("ast", rootNode.Token.Lit, nil)
	makeGraph(rootNode, g)
	fmt.Print(g.String())
}

func makeGraph(n *Syn.Node, g *graph.Graph) {
	if n.LeftMostChild == nil {
		return
	}
	child := n.LeftMostChild
	for child != nil {
		g.AddNode("ast", child.Token.Lit, nil)
		g.AddEdge(n.Token.Lit, child.Token.Lit, true, nil)
		makeGraph(child, g)
	}
}
