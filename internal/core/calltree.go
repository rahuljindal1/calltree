package core

type TreeNode struct {
	Name     string
	Children []*TreeNode
}

func BuildCallTree(functions map[string]*Function) map[string]*TreeNode {

	result := make(map[string]*TreeNode)

	var build func(name string, visited map[string]bool) *TreeNode

	build = func(name string, visited map[string]bool) *TreeNode {

		if visited[name] {
			return &TreeNode{
				Name:     name,
				Children: []*TreeNode{},
			}
		}

		visited[name] = true

		node := &TreeNode{
			Name:     name,
			Children: []*TreeNode{},
		}

		fn := functions[name]
		if fn == nil {
			return node
		}

		for _, called := range fn.Calls {
			if functions[called] != nil {
				node.Children = append(
					node.Children,
					build(called, copyVisited(visited)),
				)
			}
		}

		return node
	}

	for name := range functions {
		result[name] = build(name, map[string]bool{})
	}

	return result
}

func copyVisited(src map[string]bool) map[string]bool {
	dst := make(map[string]bool)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
