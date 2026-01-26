package core

type LanguageParser interface {
	Parse(source []byte, fileName string) (*FileAnalysis, error)
}
