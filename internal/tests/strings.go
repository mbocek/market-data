package tests

// Ptr returns a pointer to the given value for any type T.
func Ptr[T any](v T) *T {
	return &v
}
