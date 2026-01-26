package core

type TreeNode struct {
	Name     string
	File     string
	Children []*TreeNode
}

func BuildCallTree(functions map[string]*Function) map[string]*TreeNode {

	result := make(map[string]*TreeNode)

	var build func(name string, visited map[string]bool) *TreeNode

	build = func(name string, visited map[string]bool) *TreeNode {

		if visited[name] {
			return &TreeNode{
				Name:     name,
				File:     functions[name].File,
				Children: []*TreeNode{},
			}
		}

		visited[name] = true

		fn := functions[name]

		file := ""
		calls := []string{}

		if fn != nil {
			file = fn.File
			calls = fn.Calls
		}

		node := &TreeNode{
			Name:     name,
			File:     file,
			Children: []*TreeNode{},
		}

		for _, called := range calls {

			node.Children = append(
				node.Children,
				build(called, copyVisited(visited)),
			)
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
