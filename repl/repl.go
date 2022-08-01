package repl

import (
	"bufio"
	"fmt"
	"github.com/nanjingblue/go-monkey/lexer"
	"github.com/nanjingblue/go-monkey/token"
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

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
