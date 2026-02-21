package typescript

import (
	"context"

	"github.com/rahuljindal1/calltree/internal/core"
	sitter "github.com/smacker/go-tree-sitter"
	typescript "github.com/smacker/go-tree-sitter/typescript/tsx"
)

type Parser struct {
	opts core.ParseOptions
}

var _ core.LanguageParser = (*Parser)(nil)

func NewParser(opts core.ParseOptions) core.LanguageParser {
	return &Parser{opts: opts}
}

func (p *Parser) Parse(source []byte, fileName string) (*core.FileAnalysis, error) {

	parser := sitter.NewParser()
	parser.SetLanguage(typescript.GetLanguage())

	tree, err := parser.ParseCtx(
		context.Background(),
		nil,
		source,
	)
	if err != nil {
		return nil, err
	}

	analysis := &core.FileAnalysis{
		Language:  "typescript",
		Functions: make(map[string]*core.Function),
	}

	visitor := NewVisitor(source, fileName, analysis, core.ParseOptions{IncludeBuiltins: p.opts.IncludeBuiltins})

	Walk(
		tree.RootNode(),
		visitor.Enter,
		visitor.Exit,
	)

	return analysis, nil
}
