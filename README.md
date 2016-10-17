# html-parser-lexer

Parse html content with [state-lexer](https://github.com/mh-cbon/state-lexer).

## Usage

```go
package main

import (
	"bytes"
	"fmt"
	"github.com/mh-cbon/html-parser-lexer"
	"github.com/mh-cbon/state-lexer"
	"os"
	"text/tabwriter"
)

func main() {
	input := `<html>c o n t e n t</html>`

	b := bytes.NewBufferString(input)
	l := htmlparserlexer.NewHtmlParserLexer(b)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	fmt.Fprintf(w, "%v\t %v\t %v\n", "name", "type", "value")
	l.Scan(func(tok lexer.Token) {
		fmt.Fprintf(w, "%v\t %v\t %q\n", htmlparserlexer.TokenName(tok), tok.Type, tok.Value)
	})

	w.Flush()
	//Output:
	// name             | type | value
	// TagOpenToken     | 0    | "<html"
	// TagOpenEndToken  | 3    | ">"
	// TextToken        | 11   | "c o n t e n t"
	// TagCloseToken    | 1    | "</html"
	// TagCloseEndToken | 5    | ">"
}
```
