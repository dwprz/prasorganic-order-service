package helper

func MapValues[T any](m map[string]T) []T {
	values := make([]T, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}

	return values
}
