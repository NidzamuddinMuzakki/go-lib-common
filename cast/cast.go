package cast

func NewPointer[T any](t T) *T {
	return &t
}
