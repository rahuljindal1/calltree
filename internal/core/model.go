package core

type Position struct {
	Line   uint32
	Column uint32
}

type Function struct {
	Name  string
	Start Position
	End   Position
	Calls []string
}

type FileAnalysis struct {
	Language  string
	Functions map[string]*Function
}
