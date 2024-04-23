package utils

// FilterSlice filters a slice by a provided test function.
// The output is a new slice of type T.
func FilterSlice[T any](target []T, test func(T) bool) []T {
	filtered := []T{}

	for _, item := range target {
		if test(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}
