package utils

type StringSlice []string

func (s StringSlice) Delete(i int) StringSlice {
	if len(s) == 0 || i < 0 || i > len(s)-1 {
		return s
	}

	if len(s) >= i+1 {
		return append(s[:i], s[i+1:]...)
	}

	return s[:i-1]
}
