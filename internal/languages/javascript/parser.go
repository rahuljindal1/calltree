package javascript

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	javascript "github.com/smacker/go-tree-sitter/javascript"

	"calltree/internal/core"
)

type Parser struct{}

var _ core.LanguageParser = (*Parser)(nil)

func NewParser() core.LanguageParser {
	return &Parser{}
}

func (p *Parser) Parse(source []byte, fileName string) (*core.FileAnalysis, error) {

	parser := sitter.NewParser()
	parser.SetLanguage(javascript.GetLanguage())

	tree, err := parser.ParseCtx(
		context.Background(),
		nil,
		source,
	)
	if err != nil {
		return nil, err
	}

	analysis := &core.FileAnalysis{
		Language:  "javascript",
		Functions: make(map[string]*core.Function),
	}

	visitor := NewVisitor(source, analysis)

	Walk(
		tree.RootNode(),
		visitor.Enter,
		visitor.Exit,
	)

	return analysis, nil
}
