package params

func Get[T any](t any) (T, bool) {
	m, ok := t.(T)
	return m, ok
}
