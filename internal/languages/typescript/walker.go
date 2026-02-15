package typescript

import sitter "github.com/smacker/go-tree-sitter"

// Walk performs a depth-first traversal of the syntax tree.
func Walk(
	node *sitter.Node,
	enter func(*sitter.Node),
	exit func(*sitter.Node),
) {
	if node == nil {
		return
	}

	enter(node)

	for i := 0; i < int(node.ChildCount()); i++ {
		Walk(node.Child(i), enter, exit)
	}

	exit(node)
}
