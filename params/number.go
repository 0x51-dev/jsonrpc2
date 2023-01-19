package params

import "encoding/json"

func GetFloat64(t any) (float64, bool) {
	n, ok := t.(json.Number)
	if !ok {
		return 0, false
	}
	f, err := n.Float64()
	if err != nil {
		return 0, false
	}
	return f, true
}

func GetInt64(t any) (int64, bool) {
	n, ok := t.(json.Number)
	if !ok {
		return 0, false
	}
	i, err := n.Int64()
	if err != nil {
		return 0, false
	}
	return i, true
}
