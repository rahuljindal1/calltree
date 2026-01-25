package core

type LanguageParser interface {
	Parse(source []byte) (*FileAnalysis, error)
}
