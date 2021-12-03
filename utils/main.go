package utils

// StringSlice is a slice of strings
type StringSlice []string

// DeleteIndex deletes index i in StringSlice s
func (s StringSlice) DeleteIndex(i int) StringSlice {
	if len(s) == 0 || i < 0 || i > len(s)-1 {
		return s
	}

	if len(s) >= i+1 {
		return append(s[:i], s[i+1:]...)
	}

	return s[:i-1]
}

// DeleteValue deletes the first entry with value str in StringSlice s
func (s StringSlice) DeleteValue(str string) StringSlice {
	if len(s) == 0 {
		return s
	}

	for i, v := range s {
		if str == v {
			return s.DeleteIndex(i)
		}
	}

	return s
}

// Clone returns clone of StringSlice s
func (s StringSlice) Clone() StringSlice {
	x := StringSlice{}
	x = append(x, s...)
	return x
}

// DeleteValues removes each value in toDelete from s
func (s StringSlice) DeleteValues(toDelete StringSlice) StringSlice {
	newSlice := s.Clone()

	for _, v := range toDelete {
		newSlice = newSlice.DeleteValue(v)
	}

	return newSlice
}
