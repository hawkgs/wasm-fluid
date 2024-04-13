package utils

func FilterSlice[T any](target []T, test func(T) bool) []T {
	filtered := []T{}

	for _, item := range target {
		if test(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}
