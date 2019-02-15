package commands

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

type CommandAST struct {
	Arguments []*ArgumentAST `{ @@ }`
}

type ArgumentAST struct {
	Command string `@String | @DoubleQuotedString`
}

func CreateParser() *participle.Parser {
	def, err := lexer.Regexp(`\"(?P<DoubleQuotedString>[^\"]+)\"|(?P<String>[^\s]+)|(\s+)`)
	if err != nil {
		panic(err)
	}

	def = lexer.Unquote(def, "DoubleQuotedString")

	parser, err := participle.Build(&CommandAST{}, def)
	if err != nil {
		panic(err)
	}

	return parser
}
