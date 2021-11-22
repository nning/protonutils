package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringSlice_Delete_empty(t *testing.T) {
	slice := StringSlice{}
	slice = slice.Delete(0)

	assert.Equal(t, slice, StringSlice{})
}

func Test_StringSlice_Delete_first(t *testing.T) {
	slice := StringSlice{"0", "1"}
	slice = slice.Delete(0)

	assert.Equal(t, slice, StringSlice{"1"})
}

func Test_StringSlice_Delete_last(t *testing.T) {
	slice := StringSlice{"0", "1"}
	slice = slice.Delete(1)

	assert.Equal(t, slice, StringSlice{"0"})

	slice = StringSlice{"0", "1", "2"}
	slice = slice.Delete(2)

	assert.Equal(t, slice, StringSlice{"0", "1"})
}

func Test_StringSlice_Delete_middle(t *testing.T) {
	slice := StringSlice{"0", "1", "2"}
	slice = slice.Delete(1)

	assert.Equal(t, slice, StringSlice{"0", "2"})

	slice = StringSlice{"0", "1", "2", "3"}
	slice = slice.Delete(1)

	assert.Equal(t, slice, StringSlice{"0", "2", "3"})

	slice = StringSlice{"0", "1", "2", "3"}
	slice = slice.Delete(2)

	assert.Equal(t, slice, StringSlice{"0", "1", "3"})
}

func Test_StringSlice_Delete_missing(t *testing.T) {
	slice := StringSlice{"0"}
	slice = slice.Delete(1)

	assert.Equal(t, slice, StringSlice{"0"})

	slice = StringSlice{"0"}
	slice = slice.Delete(2)

	assert.Equal(t, slice, StringSlice{"0"})

	slice = StringSlice{"0"}
	slice = slice.Delete(-1)

	assert.Equal(t, slice, StringSlice{"0"})
}
