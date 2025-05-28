package repo

type Conn interface {
	GetEngine() any
}

type QuerySpec interface {
	Apply(conn any) (any, error)
	And(other QuerySpec) QuerySpec
	Or(other QuerySpec) QuerySpec
}
