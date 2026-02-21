package typescript

import (
	"strings"

	"github.com/rahuljindal1/calltree/internal/core"
	sitter "github.com/smacker/go-tree-sitter"
)

type Visitor struct {
	source          []byte
	fileName        string
	analysis        *core.FileAnalysis
	fnStack         []*core.Function
	classStack      []string
	includeBuiltins bool
}

func NewVisitor(
	source []byte,
	fileName string,
	analysis *core.FileAnalysis,
	opts core.ParseOptions,
) *Visitor {
	return &Visitor{
		source:          source,
		fileName:        fileName,
		analysis:        analysis,
		fnStack:         []*core.Function{},
		classStack:      []string{},
		includeBuiltins: opts.IncludeBuiltins,
	}
}

func (v *Visitor) current() *core.Function {
	if len(v.fnStack) == 0 {
		return nil
	}
	return v.fnStack[len(v.fnStack)-1]
}

func (v *Visitor) Enter(node *sitter.Node) {

	switch node.Type() {

	// ----------------------------------------
	// class Foo { ... }
	// ----------------------------------------
	case "class_declaration":
		nameNode := node.ChildByFieldName("name")
		if nameNode != nil {
			v.classStack = append(
				v.classStack,
				nameNode.Content(v.source),
			)
		}

	// ----------------------------------------
	// function foo() {}
	// ----------------------------------------
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
		v.fnStack = append(v.fnStack, fn)

	// ----------------------------------------
	// class Foo { bar() {} }
	// ----------------------------------------
	case "method_definition":

		nameNode := node.ChildByFieldName("name")
		if nameNode == nil {
			return
		}

		method := nameNode.Content(v.source)

		qualified := method
		if len(v.classStack) > 0 {
			qualified = v.classStack[len(v.classStack)-1] + "." + method
		}

		fn := &core.Function{
			Name:  qualified,
			File:  v.fileName,
			Calls: []string{},
		}

		v.analysis.Functions[qualified] = fn
		v.fnStack = append(v.fnStack, fn)

	// ----------------------------------------
	// foo(), this.foo(), obj.foo(), a.b.c()
	// ----------------------------------------
	case "call_expression":

		current := v.current()
		if current == nil {
			return
		}

		callee := node.ChildByFieldName("function")
		if callee == nil {
			return
		}

		var callName string

		switch callee.Type() {
		case "identifier":
			callName = callee.Content(v.source)

		case "member_expression":
			callName = extractMemberName(callee, v.source)
		}

		if callName == "" {
			return
		}

		if !v.includeBuiltins && isIgnoredCall(stripMember(callName)) {
			return
		}

		current.Calls = append(current.Calls, callName)
	}
}

func (v *Visitor) Exit(node *sitter.Node) {

	switch node.Type() {

	case "class_declaration":
		if len(v.classStack) > 0 {
			v.classStack = v.classStack[:len(v.classStack)-1]
		}

	case "function_declaration", "method_definition":
		if len(v.fnStack) > 0 {
			v.fnStack = v.fnStack[:len(v.fnStack)-1]
		}
	}
}

// ----------------------------------------
// Helpers
// ----------------------------------------

func extractMemberName(node *sitter.Node, source []byte) string {

	object := node.ChildByFieldName("object")
	property := node.ChildByFieldName("property")

	if property == nil {
		return ""
	}

	prop := property.Content(source)

	if object == nil {
		return prop
	}

	switch object.Type() {

	case "identifier":
		return object.Content(source) + "." + prop

	case "member_expression":
		parent := extractMemberName(object, source)
		if parent != "" {
			return parent + "." + prop
		}
	}

	return prop
}

func stripMember(name string) string {
	if i := strings.LastIndex(name, "."); i >= 0 {
		return name[i+1:]
	}
	return name
}
