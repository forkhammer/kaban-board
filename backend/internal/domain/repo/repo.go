package repo

type RetrieveRepo[T any, TID comparable] interface {
	Get(Id TID) (*T, error)
}

type ListRepo[T any] interface {
	List(spec QuerySpec) ([]T, error)
	Count(spec QuerySpec) (int64, error)
}

type CreateRepo[T any] interface {
	Create(item *T) (*T, error)
}

type UpdateRepo[T any] interface {
	Update(item *T) (*T, error)
}

type DeleteRepo[TID comparable] interface {
	Delete(id TID) error
}

type ReadRepo[T any, TID comparable] interface {
	RetrieveRepo[T, TID]
	ListRepo[T]
}

type RWRepo[T any, TID comparable] interface {
	ReadRepo[T, TID]
	CreateRepo[T]
	UpdateRepo[T]
	DeleteRepo[TID]
}
