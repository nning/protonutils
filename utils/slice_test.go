package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Slice_instantiation(t *testing.T) {
	s := Slice[string]{"foo"}
	assert.Equal(t, 1, len(s))
}

func Test_Slice_DeleteIndex_empty(t *testing.T) {
	s := Slice[string]{}
	s = s.DeleteIndex(0)

	assert.Equal(t, s, Slice[string]{})
}

func Test_Slice_DeleteIndex_first(t *testing.T) {
	s := Slice[int]{0, 1}
	s = s.DeleteIndex(0)

	assert.Equal(t, s, Slice[int]{1})
}

func Test_Slice_DeleteIndex_last(t *testing.T) {
	type ty = float64

	s := Slice[ty]{0.0, 1.0}
	s = s.DeleteIndex(1)

	assert.Equal(t, s, Slice[ty]{0.0})

	s = Slice[ty]{0.0, 1.0, 2.0}
	s = s.DeleteIndex(2)

	assert.Equal(t, s, Slice[ty]{0.0, 1.0})
}

func Test_Slice_DeleteIndex_middle(t *testing.T) {
	s := Slice[string]{"0", "1", "2"}
	s = s.DeleteIndex(1)

	assert.Equal(t, s, Slice[string]{"0", "2"})

	s = Slice[string]{"0", "1", "2", "3"}
	s = s.DeleteIndex(1)

	assert.Equal(t, s, Slice[string]{"0", "2", "3"})

	s = Slice[string]{"0", "1", "2", "3"}
	s = s.DeleteIndex(2)

	assert.Equal(t, s, Slice[string]{"0", "1", "3"})
}

func Test_Slice_DeleteIndex_missing(t *testing.T) {
	type ty = *int
	i := 0
	p := &i

	s := Slice[ty]{p}
	s = s.DeleteIndex(1)

	assert.Equal(t, s, Slice[ty]{p})

	s = Slice[ty]{p}
	s = s.DeleteIndex(2)

	assert.Equal(t, s, Slice[ty]{p})

	s = Slice[ty]{p}
	s = s.DeleteIndex(-1)

	assert.Equal(t, s, Slice[ty]{p})

	s = Slice[ty]{p}
	s = s.DeleteIndex(0)

	assert.Equal(t, s, Slice[ty]{})
	assert.Equal(t, 0, len(s))
}

func Test_Slice_DeleteValue(t *testing.T) {
	s := Slice[string]{"0"}
	s = s.DeleteValue("0")

	assert.Equal(t, s, Slice[string]{})

	s = Slice[string]{"0", "1"}
	s = s.DeleteValue("0")

	assert.Equal(t, s, Slice[string]{"1"})

	s = Slice[string]{"0", "1"}
	s = s.DeleteValue("1")

	assert.Equal(t, s, Slice[string]{"0"})

	s = Slice[string]{"0", "1", "2"}
	s = s.DeleteValue("1")

	assert.Equal(t, s, Slice[string]{"0", "2"})

	s = Slice[string]{"0", "1", "2", "3"}
	s = s.DeleteValue("2")

	assert.Equal(t, s, Slice[string]{"0", "1", "3"})

	s = Slice[string]{}
	s = s.DeleteValue("0")

	assert.Equal(t, s, Slice[string]{})
}