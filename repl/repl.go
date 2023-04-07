package repl

import (
	"bufio"
	"fmt"
	"github.com/nanjingblue/go-monkey/compiler"
	"github.com/nanjingblue/go-monkey/lexer"
	"github.com/nanjingblue/go-monkey/parser"
	"github.com/nanjingblue/go-monkey/vm"
	"io"
)

const PROMPT = ">>"
const GOMONKET = `
  __  __     __ __  __  __  _ _  _______   __
 / _]/__\ __|  V  |/__\|  \| | |/ / __\ 'v' /
| [/\ \/ |__| \_/ | \/ | | ' |   <| _| '. .'
\__/\__/   |_| |_|\__/|_|\__|_|\_\___| !_!
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	io.WriteString(out, GOMONKET)
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}

		lastPopped := machine.LastPoppedStackElem()
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, "parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
