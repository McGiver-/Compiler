package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/McGiver-/Compiler/refactor/Lex"
	"github.com/McGiver-/Compiler/refactor/Syn"
	"github.com/McGiver-/Compiler/refactor/Syn/Sem"
	"github.com/McGiver-/Compiler/refactor/Syn/ast"
	graph "github.com/awalterschulze/gographviz"
	"github.com/olekukonko/tablewriter"
)

var suffix = 1
var errs = []error{}

func main() {

	outputTable := flag.String("table", "", "output of the symbol table")
	outputGraph := flag.String("graph", "", "output of the graph made with graphiz")
	outputTokens := flag.String("tokens", "", "output of tokens from the Lexer")
	var (
		tknsOut  *os.File
		dotFile  *os.File
		tableOut *os.File
	)
	flag.Parse()
	src := flag.Args()[0]
	r, err := os.Open(src)
	if err != nil {
		log.Fatal("could not open source file")
	}
	defer r.Close()
	defer func() {
		if r := recover(); r != nil {
			for _, err := range errs {
				fmt.Printf("%v\n", err)
			}
		}
	}()
	lexer, err := Lex.CreateLexer(r)
	if err != nil {
		log.Fatal(err)
	}

	tc, errs := lexer.GetTokensNoChan()

	if *outputTokens == "" {
		tknsOut, err = ioutil.TempFile("", "tokensOut")
		if err != nil {
			log.Fatal("could not open tokens temp file")
		}
		defer os.Remove(tknsOut.Name())
	} else {
		tknsOut, err = os.Create(*outputTokens)
		if err != nil {
			log.Fatal("could not open tokens file")
		}
		defer tknsOut.Close()
	}

	for _, tk := range tc {
		fmt.Fprintf(tknsOut, "type:<%s> lexeme:<%s> position:<%s>\n", tk.Token.String(), tk.Lit, tk.Position)
	}
	analyzer, err := Syn.CreateAnalyzer(tc)
	if err != nil {
		log.Fatal(err)
	}

	ec, rootNode := analyzer.Parse()
	nodeList := getNodes(rootNode)
	replaceLeading(nodeList)
	strip(nodeList)
	errs = append(errs, ec...)

	if *outputGraph == "" {
		dotFile, err = ioutil.TempFile("", "dotfile")
		if err != nil {
			log.Fatal("could not open graph temp file")
		}
		defer os.Remove(dotFile.Name())
	} else {
		dotFile, err = os.Create(*outputGraph)
		if err != nil {
			log.Fatal("could not open graph file")
		}
		defer dotFile.Close()
	}

	g := graph.NewGraph()
	if err := g.SetName("ast"); err != nil {
		log.Fatal("could not set name of graph")
	}
	if err := g.SetDir(true); err != nil {
		log.Fatal("could not set dir of graph")
	}

	rootNodeName := fmt.Sprintf("%s_%d", rootNode.Type, suffix)
	g.AddNode("ast", rootNodeName, nil)
	suffix++
	makeGraph(rootNode, g, rootNodeName)
	fmt.Fprint(dotFile, g.String())
	dotFile.Close()

	makeTableVisitor := &ast.TableCreationVisitor{}
	errs = append(errs, rootNode.Accept(makeTableVisitor)...)
	errs = append(errs, rootNode.Table.GetShadows()...)
	errs = append(errs, rootNode.CheckUndeclaredDataMemeber()...)
	rootNode.SetFcallType(rootNode.Table)
	errs = append(errs, rootNode.CheckReturnParams()...)
	errs = append(errs, rootNode.CheckReturnType()...)
	errs = append(errs, rootNode.CheckCalledFuncDeclared()...)
	errs = append(errs, rootNode.Table.CheckMemberFuncNoDef()...)
	errs = sortErrors(errs, tc)
	for _, err := range errs {
		fmt.Printf("%v\n", err)
	}
	if *outputTable == "" {
		tableOut, err = ioutil.TempFile("", "tableOut")
		if err != nil {
			log.Fatal("could not open table temp file")
		}
		defer os.Remove(tableOut.Name())
	} else {
		tableOut, err = os.Create(*outputTable)
		if err != nil {
			log.Fatal("could not open table file")
		}
		defer tableOut.Close()
	}

	preorder(rootNode.Table, printTable, tableOut)
	// fmt.Printf("%s", spew.Sdump(rootNode.Table))
}

func printTable(table *Sem.Table, w io.Writer) {
	if table == nil {
		return
	}
	tw := tablewriter.NewWriter(w)
	tw.SetHeader([]string{"name", "kind", "type", "link"})
	tw.SetCaption(true, table.Name)
	for _, v := range table.Entries {
		if v == nil {
			continue
		}
		link := "false"
		if v.Child != nil {
			link = "true"
		}
		tw.Append([]string{v.Name, v.Kind, v.Typ, link})
	}
	tw.Render()
}

func preorder(table *Sem.Table, fn func(*Sem.Table, io.Writer), w io.Writer) {
	fn(table, w)
	if table == nil {
		return
	}
	for _, v := range table.Entries {
		if v == nil {
			continue
		}
		preorder(v.Child, fn, w)
	}
}

func makeGraph(n *ast.Node, g *graph.Graph, pName string) {
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

func getNodes(n *ast.Node) []*ast.Node {
	if n == nil {
		return []*ast.Node{}
	}
	nodeList := []*ast.Node{}
	child := n.LeftMostChild
	for child != nil {
		nodeList = append(nodeList, child)
		nodeList = append(nodeList, getNodes(child)...)
		child = child.RightSibling
	}
	return nodeList
}

func strip(nodeList []*ast.Node) {
	for i := 0; i < len(nodeList); i++ {
		if nodeList[i] != nil && nodeList[i].RightSibling != nil && nodeList[i].RightSibling.Type == "EPSILON" {
			nodeList[i].RightSibling = nil
		}
		if nodeList[i] != nil && nodeList[i].LeftMostChild != nil && nodeList[i].LeftMostChild.Type == "EPSILON" {
			nodeList[i].LeftMostChild = nil
		}
	}
}

func replaceLeading(nodeList []*ast.Node) {
	for i := 0; i < len(nodeList); i++ {
		if nodeList[i] != nil &&
			nodeList[i].LeftMostChild != nil &&
			nodeList[i].LeftMostChild.Type == "EPSILON" &&
			nodeList[i].LeftMostChild.RightSibling != nil {
			nodeList[i].LeftMostChild = nodeList[i].LeftMostChild.RightSibling
		}
	}
}

func sortErrors(errors []error, tks []*Lex.Token) []error {

	for i := 0; i < len(errors); i++ {
		if string(errors[i].Error()[0]) == "<" {
			lit := strings.Split(strings.Split(errors[i].Error(), "<")[1], ">")[0]
			for _, k := range tks {
				if lit == k.Lit {
					errors[i] = fmt.Errorf("%s %s", k.Position, errors[i].Error())
				}
			}
		}
	}

	for i := 0; i < len(errors)-1; i++ {
		for j := i; j < len(errors)-2; j++ {
			p1 := strings.Split(strings.Split(errors[j].Error(), " ")[0], ":")
			l1, _ := strconv.Atoi(p1[0])
			p2 := strings.Split(strings.Split(errors[j+1].Error(), " ")[0], ":")
			l2, _ := strconv.Atoi(p2[0])
			if l1 > l2 {
				errors[j], errors[j+1] = errors[j+1], errors[j]
			}
		}
	}

	for i := 0; i < len(errors)-1; i++ {
		for j := i; j < len(errors)-2; j++ {
			p1 := strings.Split(strings.Split(errors[j].Error(), " ")[0], ":")
			l1, _ := strconv.Atoi(p1[0])
			c1, _ := strconv.Atoi(p1[1])
			p2 := strings.Split(strings.Split(errors[j+1].Error(), " ")[0], ":")
			l2, _ := strconv.Atoi(p2[0])
			c2, _ := strconv.Atoi(p2[1])
			if l1 == l2 && c1 > c2 {
				errors[j], errors[j+1] = errors[j+1], errors[j]
			}
		}
	}

	return removeDuplicates(errors)
}

func removeDuplicates(errors []error) []error {
	newErrs := []error{}
	found := map[string]bool{}
	for i := 0; i < len(errors); i++ {
		f := found[errors[i].Error()]
		if !f {
			newErrs = append(newErrs, errors[i])
			found[errors[i].Error()] = true
		}
	}
	return newErrs
}
