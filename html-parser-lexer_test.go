package htmlparserlexer

import (
	"bytes"
	"fmt"
	"github.com/mh-cbon/state-lexer"
	"os"
	"testing"
	"text/tabwriter"
)

func Test_Fixture1(t *testing.T) {

	expected := []lexer.Token{
		lexer.Token{Type: 2, Value: "<"},
		lexer.Token{Type: 0, Value: "html"},
		lexer.Token{Type: 3, Value: ">"},
		lexer.Token{Type: 4, Value: "</"},
		lexer.Token{Type: 1, Value: "html"},
		lexer.Token{Type: 5, Value: ">"},
	}

	input := `<html></html>`

	runTest(t, expected, input, false)
}

func Test_Fixture2(t *testing.T) {

	expected := []lexer.Token{
		lexer.Token{Type: 10, Value: "\n"},
		lexer.Token{Type: 2, Value: "<"},
		lexer.Token{Type: 0, Value: "html"},
		lexer.Token{Type: 3, Value: ">"},
		lexer.Token{Type: 10, Value: "\n"},
		lexer.Token{Type: 2, Value: "<"},
		lexer.Token{Type: 0, Value: "input"},
		lexer.Token{Type: 10, Value: " "},
		lexer.Token{Type: 6, Value: "type"},
		lexer.Token{Type: 7, Value: "="},
		lexer.Token{Type: 8, Value: "\""},
		lexer.Token{Type: 9, Value: "button"},
		lexer.Token{Type: 8, Value: "\""},
		lexer.Token{Type: 10, Value: " "},
		lexer.Token{Type: 6, Value: "disabled"},
		lexer.Token{Type: 10, Value: " "},
		lexer.Token{Type: 6, Value: "value"},
		lexer.Token{Type: 7, Value: "="},
		lexer.Token{Type: 9, Value: "42"},
		lexer.Token{Type: 10, Value: " "},
		lexer.Token{Type: 3, Value: "/>"},
		lexer.Token{Type: 10, Value: "\n"},
		lexer.Token{Type: 4, Value: "</"},
		lexer.Token{Type: 1, Value: "html"},
		lexer.Token{Type: 5, Value: ">"},
		lexer.Token{Type: 10, Value: "\n"},
	}

	input := `
<html>
<input type="button" disabled value=42 />
</html>
`

	runTest(t, expected, input, false)
}

func Test_Fixture3(t *testing.T) {

	expected := []lexer.Token{
		lexer.Token{Type: 2, Value: "<"},
		lexer.Token{Type: 0, Value: "input"},
		lexer.Token{Type: 10, Value: " "},
		lexer.Token{Type: 6, Value: "type"},
		lexer.Token{Type: 7, Value: "="},
		lexer.Token{Type: 8, Value: "\""},
		lexer.Token{Type: 9, Value: "butt\\\"on"},
		lexer.Token{Type: 8, Value: "\""},
		lexer.Token{Type: 10, Value: " "},
		lexer.Token{Type: 3, Value: "/>"},
	}

	input := `<input type="butt\"on" />`

	runTest(t, expected, input, false)
}

func Test_Fixture4(t *testing.T) {

	expected := []lexer.Token{
		lexer.Token{Type: 12, Value: "<!--"},
		lexer.Token{Type: 13, Value: "comments works too"},
		lexer.Token{Type: 14, Value: "-->"},
	}

	input := `<!--comments works too-->`

	runTest(t, expected, input, false)
}

func Test_Fixture5(t *testing.T) {

	expected := []lexer.Token{
		lexer.Token{Type: 11, Value: "whatever "},
		lexer.Token{Type: 2, Value: "<"},
		lexer.Token{Type: 0, Value: "html"},
		lexer.Token{Type: 3, Value: ">"},
		lexer.Token{Type: 10, Value: " "},
		lexer.Token{Type: 11, Value: "it works ?"},
	}

	input := `whatever <html> it works ?`

	runTest(t, expected, input, false)
}

func Test_Fixture6(t *testing.T) {

	expected := []lexer.Token{
		lexer.Token{Type: 2, Value: "<"},
		lexer.Token{Type: 0, Value: "html"},
		lexer.Token{Type: 10, Value: " "},
		lexer.Token{Type: 3, Value: ""},
		lexer.Token{Type: 2, Value: "<"},
		lexer.Token{Type: 0, Value: "div"},
		lexer.Token{Type: 3, Value: ">"},
		lexer.Token{Type: 10, Value: " "},
		lexer.Token{Type: 11, Value: "shame this..."},
	}

	input := `<html <div> shame this...`

	runTest(t, expected, input, false)
}

func Test_Fixture7(t *testing.T) {

	expected := []lexer.Token{
		lexer.Token{Type: 2, Value: "<"},
		lexer.Token{Type: 0, Value: "html"},
		lexer.Token{Type: 3, Value: ">"},
		lexer.Token{Type: 10, Value: "\n\n"},
		lexer.Token{Type: 2, Value: "<"},
		lexer.Token{Type: 0, Value: "esi:include"},
		lexer.Token{Type: 10, Value: " "},
		lexer.Token{Type: 6, Value: "src"},
		lexer.Token{Type: 7, Value: "="},
		lexer.Token{Type: 8, Value: "\""},
		lexer.Token{Type: 9, Value: "http://example.com/1.html"},
		lexer.Token{Type: 8, Value: "\""},
		lexer.Token{Type: 10, Value: " "},
		lexer.Token{Type: 6, Value: "alt"},
		lexer.Token{Type: 7, Value: "="},
		lexer.Token{Type: 8, Value: "\""},
		lexer.Token{Type: 9, Value: "http://bak.example.com/2.html"},
		lexer.Token{Type: 8, Value: "\""},
		lexer.Token{Type: 10, Value: " "},
		lexer.Token{Type: 6, Value: "onerror"},
		lexer.Token{Type: 7, Value: "="},
		lexer.Token{Type: 8, Value: "\""},
		lexer.Token{Type: 9, Value: "continue"},
		lexer.Token{Type: 8, Value: "\""},
		lexer.Token{Type: 3, Value: "/>"},
		lexer.Token{Type: 10, Value: "\n\n"},
		lexer.Token{Type: 4, Value: "</"},
		lexer.Token{Type: 1, Value: "html"},
		lexer.Token{Type: 5, Value: ">"},
		lexer.Token{Type: 10, Value: "\n"},
	}

	input := `<html>

<esi:include src="http://example.com/1.html" alt="http://bak.example.com/2.html" onerror="continue"/>

</html>
`

	runTest(t, expected, input, false)
}

func Test_Fixture8(t *testing.T) {

	expected := []lexer.Token{
		lexer.Token{Type: 2, Value: "<"},
		lexer.Token{Type: 0, Value: "html"},
		lexer.Token{Type: 3, Value: ">"},
		lexer.Token{Type: 10, Value: "\n  "},
		lexer.Token{Type: 11, Value: "div>\n"}, // -> TextToken
		lexer.Token{Type: 4, Value: "</"},
		lexer.Token{Type: 1, Value: "html"},
		lexer.Token{Type: 10, Value: "\n"},
		lexer.Token{Type: 5, Value: ""}, // -> empty TagCloseEndToken (auto closed)
		lexer.Token{Type: 11, Value: "Really broken..did you drink ? ah ah ah\n"},
	}

	input := `<html>
  div>
</html
Really broken..did you drink ? ah ah ah
`

	runTest(t, expected, input, false)
}

func runTest(t *testing.T, expected []lexer.Token, input string, debug bool) {

	b := bytes.NewBufferString(input)
	l := NewHtmlParserLexer(b)

	var tokens []lexer.Token
	l.Scan(func(tok lexer.Token) {
		tokens = append(tokens, tok)
	})

	if debug {
		fmt.Printf("%#v", tokens)
	}

	for i, e := range expected {
		if e.Type != tokens[i].Type {
			t.Errorf("Expected token type %v but got %v", e.Type, tokens[i].Type)
		}
		if e.Value != tokens[i].Value {
			t.Errorf("Expected token value %v but got %v", e.Value, tokens[i].Value)
		}
	}

	if len(expected) != len(tokens) {
		t.Errorf("Incoherent parsing expected %v tokens but got %v", len(expected), len(tokens))
	}
}

func Example_lexer() {
	input := `<html>c o n t e n t</html>`

	b := bytes.NewBufferString(input)
	l := NewHtmlParserLexer(b)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	fmt.Fprintf(w, "%v\t %v\t %v\n", "name", "type", "value")
	l.Scan(func(tok lexer.Token) {
		fmt.Fprintf(w, "%v\t %v\t %q\n", TokenName(tok), tok.Type, tok.Value)
	})

	w.Flush()
	//Output:
	// name               | type | value
	// TagOpenStartToken  | 2    | "<"
	// TagOpenToken       | 0    | "html"
	// TagOpenEndToken    | 3    | ">"
	// TextToken          | 11   | "c o n t e n t"
	// TagCloseStartToken | 4    | "</"
	// TagCloseToken      | 1    | "html"
	// TagCloseEndToken   | 5    | ">"

}
