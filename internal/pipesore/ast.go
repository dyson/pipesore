package pipesore

type ast struct {
	functions []function
}

func newAST() *ast {
	return &ast{
		functions: []function{},
	}
}

type function struct {
	name      string
	arguments []any
}

func (f function) isNot() bool {
	return f.name[0:1] == "!"
}
