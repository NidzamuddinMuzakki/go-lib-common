package cast

func NewPointer[T any](t T) *T {
	return &t
}

// K = source model
// V = target model
// V as a pointer since pointer of interface doesnt give us as reference
// we dont want any memory allocation bigger than this so we must append v as a pointer instead
func Convert[K any, V any](k *[]K, v *[]V, f func(*K) V) {
	for idx := range *k {
		*v = append(*v, f(&(*k)[idx]))
	}
}

// with v allocation memory
// K = source model
// V = target model
// V as a pointer since pointer of interface doesnt give us as reference
// we dont want any memory allocation bigger than this so we must append v as a pointer instead
func ConvertAndAllocate[K any, V any](k *[]K, v *[]V, f func(*K) V) {
	*v = make([]V, 0, len(*k))
	for idx := range *k {
		*v = append(*v, f(&(*k)[idx]))
	}
}
