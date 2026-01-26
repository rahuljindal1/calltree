package core

type ParseOptions struct {
	IncludeBuiltins bool
}

type LanguageParser interface {
	Parse(source []byte, fileName string) (*FileAnalysis, error)
}
