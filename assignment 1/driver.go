package main

import(
	"github.com/McGiver-/Compiler/Lex"
	"fmt"
	"log"
	"os"
	"io/ioutil"
)



func main(){
	tokens := []*Lex.Token{}
	errs := []error{}
	inputFile := os.Args[1]
	FileA2CC := os.Args[2]
	ErrorFile := os.Args[3]
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
			errs = append(errs,err)
			continue		
		}
		if tkn == nil {
			continue
		}
		tokens = append(tokens,tkn)
	}
	types := "" 
	for _,v := range tokens {
		types += v.Type +" "
		fmt.Printf("%v\n",v.Type)
	}
	ioutil.WriteFile(FileA2CC,[]byte(types),os.ModePerm)

	errors := ""
	for _,v := range errs {
		errors += v.Error() +"\n"
		fmt.Printf("%v \n",v)
	}
	ioutil.WriteFile(ErrorFile,[]byte(errors),os.ModePerm)
}