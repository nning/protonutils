package set

import (
	"sort"
)

type Set map[string]bool

func (set Set) Add(s string) {
	set[s] = true
}

func Init(set Set) Set {
	if set == nil {
		return make(Set)
	}

	return set
}

func (set Set) Includes(s string) bool {
	return set[s]
}

func (set Set) Sort() []string {
	var keys []string

	for key := range set {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}
