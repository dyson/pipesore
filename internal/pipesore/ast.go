package pipesore

type ast struct {
	filters []filter
}

func newAST() *ast {
	return &ast{
		filters: []filter{},
	}
}

type filter struct {
	name      string
	arguments []any
	position
}

func (f filter) isNot() bool {
	return f.name[0:1] == "!"
}
