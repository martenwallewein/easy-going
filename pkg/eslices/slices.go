package eslices

func AppendToSliceIfMissing[T comparable](slice []T, i T) []T {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func RemoveByIndex[T comparable](index int, s []T) []T {
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

func RemoveFromSlice[T comparable](s []T, i T) []T {
	index := IndexOf(i, s)
	if index == -1 {
		return s
	}
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

func IndexOf[T comparable](element T, data []T) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
