package javascript

import (
	sitter "github.com/smacker/go-tree-sitter"

	"calltree/internal/core"
)

type Visitor struct {
	source   []byte
	fileName string
	analysis *core.FileAnalysis
	stack    []*core.Function
}

func NewVisitor(
	source []byte,
	fileName string,
	analysis *core.FileAnalysis,
) *Visitor {
	return &Visitor{
		source:   source,
		analysis: analysis,
		fileName: fileName,
		stack:    []*core.Function{},
	}
}

func (v *Visitor) current() *core.Function {
	if len(v.stack) == 0 {
		return nil
	}
	return v.stack[len(v.stack)-1]
}

func (v *Visitor) Enter(node *sitter.Node) {

	switch node.Type() {

	// ------------------------------------------------
	// function foo() {}
	// ------------------------------------------------
	case "function_declaration":

		nameNode := node.ChildByFieldName("name")
		if nameNode == nil {
			return
		}

		name := nameNode.Content(v.source)

		fn := &core.Function{
			Name:  name,
			File:  v.fileName,
			Calls: []string{},
		}

		v.analysis.Functions[name] = fn
		v.stack = append(v.stack, fn)

	// ------------------------------------------------
	// foo()
	// ------------------------------------------------
	case "call_expression":

		current := v.current()
		if current == nil {
			return
		}

		callee := node.ChildByFieldName("function")
		if callee == nil {
			return
		}

		if callee.Type() == "identifier" {
			current.Calls = append(
				current.Calls,
				callee.Content(v.source),
			)
		}
	}
}

func (v *Visitor) Exit(node *sitter.Node) {

	if node.Type() == "function_declaration" {
		if len(v.stack) > 0 {
			v.stack = v.stack[:len(v.stack)-1]
		}
	}
}
