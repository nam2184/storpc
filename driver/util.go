package driver

func Truncate[T any](s *[]T, index int) {
	if index < 0 || index > len(*s) {
		panic("truncate: index out of range")
	}
	var toClear []T
	*s, toClear = (*s)[:index], (*s)[index:]
	var zero T
	for i := range toClear {
		toClear[i] = zero
	}
}
