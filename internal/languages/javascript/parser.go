package javascript

import "calltree/internal/core"

type Parser struct{}

// compile-time check
var _ core.LanguageParser = (*Parser)(nil)

func NewParser() core.LanguageParser {
	return &Parser{}
}

func (p *Parser) Parse(source []byte) (*core.FileAnalysis, error) {
	return &core.FileAnalysis{
		Language:  "javascript",
		Functions: map[string]*core.Function{}, // !Note: Why not default set it to Empty?
	}, nil
}
