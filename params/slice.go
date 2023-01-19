package params

func GetSlice[T any](t any) ([]T, bool) {
	s, ok := t.([]any)
	if !ok {
		return nil, false
	}
	var slice []T
	for _, v := range s {
		v, ok := v.(T)
		if !ok {
			return nil, false
		}
		slice = append(slice, v)
	}
	return slice, true
}
