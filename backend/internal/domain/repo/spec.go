package repo

type QuerySpec interface {
	Apply(conn any) (any, error)
}

type BaseQuerySpec struct{}

func (q *BaseQuerySpec) Apply(conn any) (any, error) {
	return conn, nil
}

type AndQuerySpec struct {
	BaseQuerySpec
	Left  QuerySpec
	Right QuerySpec
}

func (q AndQuerySpec) Apply(conn any) (any, error) {
	qs, err := q.Left.Apply(conn)
	if err != nil {
		return nil, err
	}
	qs, err = q.Right.Apply(qs)
	if err != nil {
		return nil, err
	}
	return qs, nil
}

func And(left, right QuerySpec) QuerySpec {
	return &AndQuerySpec{Left: left, Right: right}
}
