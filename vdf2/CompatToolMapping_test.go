package vdf2

import (
	"testing"

	"github.com/nning/protonutils/steam"
	"github.com/stretchr/testify/assert"
)

func Test_GetCompatToolMapping(t *testing.T) {
	t.Parallel()

	s, err := steam.New("", testSteamRoot, true)
	assert.Empty(t, err)

	cmt, err := GetCompatToolMapping(s)
	assert.Empty(t, err)
	assert.NotEmpty(t, cmt)
	assert.NotEmpty(t, cmt.Root)
	assert.NotEmpty(t, cmt.Node)
	assert.Equal(t, "../test/root/config/config.vdf", cmt.Path)
	assert.Equal(t, s, cmt.Steam)
}

func Test_ReadCompatTools(t *testing.T) {
	t.Parallel()

	s, err := steam.New("", testSteamRoot, true)
	assert.Empty(t, err)

	cmt, err := GetCompatToolMapping(s)
	assert.Empty(t, err)
	assert.NotEmpty(t, cmt)

	compatTools, err := cmt.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 7, len(compatTools))

	for compatToolID, compatTool := range compatTools {
		assert.Equal(t, compatToolID, compatTool.ID)

		if compatTool.IsDefault {
			assert.Equal(t, "", compatTool.ID)
			assert.Equal(t, "Proton 6.3-8 (Default)", compatTool.Name)
		} else {
			assert.NotEqual(t, "", compatTool.ID)
			assert.NotEqual(t, "", compatTool.Name)
		}

		if compatTool.IsCustom {
			assert.Equal(t, compatTool.ID, compatTool.Name)
		}

		// fmt.Println(compatTool)

		assert.NotEqual(t, 0, len(compatTool.Games))
	}

	compatTools, err = cmt.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 7, len(compatTools))
}

func Test_Add(t *testing.T) {
	t.Parallel()

	s, err := steam.New("", testSteamRoot, true)
	assert.Empty(t, err)

	cmt, err := GetCompatToolMapping(s)
	assert.Empty(t, err)
	assert.NotEmpty(t, cmt)

	compatTools, err := cmt.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 7, len(compatTools))

	cmt.Add("1", "foo")

	compatTools, err = cmt.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 8, len(compatTools))
	assert.Equal(t, 1, len(compatTools["foo"].Games))
}

func Test_Update(t *testing.T) {
	t.Parallel()

	s, err := steam.New("", testSteamRoot, true)
	assert.Empty(t, err)

	cmt, err := GetCompatToolMapping(s)
	assert.Empty(t, err)
	assert.NotEmpty(t, cmt)

	compatTools, err := cmt.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 7, len(compatTools))

	cmt.Add("1", "foo")

	compatTools, err = cmt.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 8, len(compatTools))
	assert.Equal(t, 1, len(compatTools["foo"].Games))

	cmt.Update("1", "bar")

	compatTools, err = cmt.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 8, len(compatTools))
	assert.Equal(t, 1, len(compatTools["bar"].Games))
}
