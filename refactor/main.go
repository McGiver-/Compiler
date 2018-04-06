package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/McGiver-/Compiler/refactor/Lex"
	"github.com/McGiver-/Compiler/refactor/Syn"
	graph "github.com/awalterschulze/gographviz"
)

var suffix = 1

func main() {
	src := flag.String("source", "tester", "source file to be compiled")
	outputGraph := flag.String("graph", "graph.dot", "output of the graph made with graphiz")
	flag.Parse()
	r, err := os.Open(*src)
	if err != nil {
		log.Fatal("could not open file")
	}
	defer r.Close()
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
	nodeList := getNodes(rootNode)
	strip(nodeList)
	// fmt.Printf("%s", spew.Sdump(rootNode))
	errs = append(errs, ec...)
	for _, err := range errs {
		fmt.Printf("%v\n", err)
	}
	dotFile, err := os.Create(*outputGraph)
	if err != nil {
		log.Fatal("could not open file")
	}
	g := graph.NewGraph()
	if err := g.SetName("ast"); err != nil {
		panic(err)
	}
	if err := g.SetDir(true); err != nil {
		panic(err)
	}

	rootNodeName := fmt.Sprintf("%s_%d", rootNode.Type, suffix)
	g.AddNode("ast", rootNodeName, nil)
	suffix++
	makeGraph(rootNode, g, rootNodeName)
	fmt.Fprint(dotFile, g.String())
	dotFile.Close()
	// err = exec.Command("dot", "Tpng graph.dot > astGraph.png").Run()
	// if err != nil {
	// 	log.Fatalf("could not exec dot commant : %v", err)
	// }
}

func makeGraph(n *Syn.Node, g *graph.Graph, pName string) {
	if n.LeftMostChild == nil {
		return
	}
	child := n.LeftMostChild
	for child != nil {
		name := fmt.Sprintf("%s_%d", child.Type, suffix)
		if name == fmt.Sprintf("+_%d", suffix) {
			name = fmt.Sprintf("plus_%d", suffix)
		}
		if name == fmt.Sprintf("-_%d", suffix) {
			name = fmt.Sprintf("minus%d", suffix)
		}
		if name == fmt.Sprintf("==_%d", suffix) {
			name = fmt.Sprintf("eq%d", suffix)
		}
		g.AddNode("ast", name, nil)
		suffix++
		g.AddEdge(pName, name, true, nil)
		makeGraph(child, g, name)
		child = child.RightSibling
	}
}

func getNodes(n *Syn.Node) []*Syn.Node {
	if n == nil {
		return []*Syn.Node{}
	}
	nodeList := []*Syn.Node{}
	child := n.LeftMostChild
	for child != nil {
		nodeList = append(nodeList, child)
		nodeList = append(nodeList, getNodes(child)...)
		child = child.RightSibling
	}
	return nodeList
}

func strip(nodeList []*Syn.Node) {
	for i := 0; i < len(nodeList); i++ {
		if nodeList[i] != nil && nodeList[i].RightSibling != nil && nodeList[i].RightSibling.Type == "EPSILON" {
			nodeList[i].RightSibling = nil
		}
		if nodeList[i] != nil && nodeList[i].LeftMostChild != nil && nodeList[i].LeftMostChild.Type == "EPSILON" {
			nodeList[i].LeftMostChild = nil
		}
	}
}
