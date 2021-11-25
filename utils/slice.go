package utils

// Slice build upon Go slice and adds methods
type Slice[T comparable] []T

// DeleteIndex deletes index i in Slice s
func (s Slice[T]) DeleteIndex(i int) Slice[T] {
	if len(s) == 0 || i < 0 || i > len(s)-1 {
		return s
	}

	if len(s) >= i+1 {
		return append(s[:i], s[i+1:]...)
	}

	return s[:i-1]
}

// DeleteValue deletes the first entry with value str in Slice s
func (s Slice[T]) DeleteValue(value T) Slice[T] {
	if len(s) == 0 {
		return s
	}

	for i, v := range s {
		if value == v {
			return s.DeleteIndex(i)
		}
	}

	return s
}
