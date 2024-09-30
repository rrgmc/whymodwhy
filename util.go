package main

// dedupeSlice remove duplicates from a slice.
func dedupeSlice[T comparable](sliceList []T) []T {
	dedupeMap := make(map[T]struct{})
	var list []T

	for _, slice := range sliceList {
		if _, exists := dedupeMap[slice]; !exists {
			dedupeMap[slice] = struct{}{}
			list = append(list, slice)
		}
	}

	return list
}
