package params

func GetMap[T any](t any) (map[string]T, bool) {
	ma, ok := t.(map[string]any)
	if !ok {
		return nil, false
	}
	m := make(map[string]T)
	for k, v := range ma {
		v, ok := v.(T)
		if !ok {
			return nil, false
		}
		m[k] = v
	}
	return m, true
}
