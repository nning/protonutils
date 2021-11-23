package utils

type StringSlice []string

func (s StringSlice) DeleteIndex(i int) StringSlice {
	if len(s) == 0 || i < 0 || i > len(s)-1 {
		return s
	}

	if len(s) >= i+1 {
		return append(s[:i], s[i+1:]...)
	}

	return s[:i-1]
}

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
