package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringSlice_DeleteIndex_empty(t *testing.T) {
	slice := StringSlice{}
	slice = slice.DeleteIndex(0)

	assert.Equal(t, slice, StringSlice{})
}

func Test_StringSlice_DeleteIndex_first(t *testing.T) {
	slice := StringSlice{"0", "1"}
	slice = slice.DeleteIndex(0)

	assert.Equal(t, slice, StringSlice{"1"})
}

func Test_StringSlice_DeleteIndex_last(t *testing.T) {
	slice := StringSlice{"0", "1"}
	slice = slice.DeleteIndex(1)

	assert.Equal(t, slice, StringSlice{"0"})

	slice = StringSlice{"0", "1", "2"}
	slice = slice.DeleteIndex(2)

	assert.Equal(t, slice, StringSlice{"0", "1"})
}

func Test_StringSlice_DeleteIndex_middle(t *testing.T) {
	slice := StringSlice{"0", "1", "2"}
	slice = slice.DeleteIndex(1)

	assert.Equal(t, slice, StringSlice{"0", "2"})

	slice = StringSlice{"0", "1", "2", "3"}
	slice = slice.DeleteIndex(1)

	assert.Equal(t, slice, StringSlice{"0", "2", "3"})

	slice = StringSlice{"0", "1", "2", "3"}
	slice = slice.DeleteIndex(2)

	assert.Equal(t, slice, StringSlice{"0", "1", "3"})
}

func Test_StringSlice_DeleteIndex_missing(t *testing.T) {
	slice := StringSlice{"0"}
	slice = slice.DeleteIndex(1)

	assert.Equal(t, slice, StringSlice{"0"})

	slice = StringSlice{"0"}
	slice = slice.DeleteIndex(2)

	assert.Equal(t, slice, StringSlice{"0"})

	slice = StringSlice{"0"}
	slice = slice.DeleteIndex(-1)

	assert.Equal(t, slice, StringSlice{"0"})
}

func Test_StringSlice_DeleteValue(t *testing.T) {
	slice := StringSlice{"0"}
	slice = slice.DeleteValue("0")

	assert.Equal(t, slice, StringSlice{})

	slice = StringSlice{"0", "1"}
	slice = slice.DeleteValue("0")

	assert.Equal(t, slice, StringSlice{"1"})

	slice = StringSlice{"0", "1"}
	slice = slice.DeleteValue("1")

	assert.Equal(t, slice, StringSlice{"0"})

	slice = StringSlice{"0", "1", "2"}
	slice = slice.DeleteValue("1")

	assert.Equal(t, slice, StringSlice{"0", "2"})

	slice = StringSlice{"0", "1", "2", "3"}
	slice = slice.DeleteValue("2")

	assert.Equal(t, slice, StringSlice{"0", "1", "3"})

	slice = StringSlice{}
	slice = slice.DeleteValue("0")

	assert.Equal(t, slice, StringSlice{})
}
