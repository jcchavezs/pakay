package values

func GetFromMap[T any](m map[string]interface{}, key string) (T, bool) {
	var value T
	rVal, ok := m[key]
	if ok {
		value, ok = rVal.(T)
	}

	return value, ok
}
