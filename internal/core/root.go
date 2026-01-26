package core

// FindRoots returns functions that are never called by others.
func FindRoots(functions map[string]*Function) []string {

	called := make(map[string]bool)

	for _, fn := range functions {
		for _, c := range fn.Calls {
			called[c] = true
		}
	}

	roots := []string{}

	for name := range functions {
		if !called[name] {
			roots = append(roots, name)
		}
	}

	return roots
}
