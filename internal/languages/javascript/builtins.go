package javascript

var ignoredCalls = map[string]bool{
	// primitives / globals
	"Number":  true,
	"String":  true,
	"Boolean": true,
	"Object":  true,
	"Array":   true,
	"Promise": true,
	"Set":     true,
	"Map":     true,
	"Date":    true,

	// common prototype methods
	"toLowerCase": true,
	"trim":        true,
	"map":         true,
	"includes":    true,
	"filter":      true,
	"reduce":      true,
	"forEach":     true,
	"flatMap":     true,
	"push":        true,

	// framework / DI noise
	"Inject": true,
}

func isIgnoredCall(name string) bool {
	return ignoredCalls[name]
}
