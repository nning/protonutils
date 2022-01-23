package steam2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InitCompatToolMapping(t *testing.T) {
	t.Parallel()

	s, err := New("", testSteamRoot, false)
	assert.Empty(t, err)

	assert.Empty(t, err)
	assert.NotEmpty(t, s.CompatToolMapping)
	assert.NotEmpty(t, s.CompatToolMapping.Root)
	assert.NotEmpty(t, s.CompatToolMapping.Node)
	assert.Equal(t, "../test/root/config/config.vdf", s.CompatToolMapping.Path)
}

func Test_ReadCompatTools(t *testing.T) {
	t.Parallel()

	s, err := New("", testSteamRoot, false)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.CompatToolMapping)

	compatTools, err := s.CompatToolMapping.ReadCompatTools()
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

		assert.NotEqual(t, 0, len(compatTool.Games))
	}

	compatTools, err = s.CompatToolMapping.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 7, len(compatTools))
}

func Test_Add(t *testing.T) {
	t.Parallel()

	s, err := New("", testSteamRoot, false)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.CompatToolMapping)

	ctm := s.CompatToolMapping
	assert.NotEmpty(t, ctm)

	compatTools, err := ctm.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 7, len(compatTools))

	s.CompatToolMapping.Add("1", "foo")

	compatTools, err = ctm.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 8, len(compatTools))
	assert.NotEmpty(t, compatTools["foo"])
	if compatTools["foo"] != nil {
		assert.Equal(t, 1, len(compatTools["foo"].Games))
	}
}

func Test_Update(t *testing.T) {
	t.Parallel()

	s, err := New("", testSteamRoot, false)
	assert.Empty(t, err)

	// ctm := s.CompatToolMapping
	// assert.NotEmpty(t, ctm)

	compatTools, err := s.CompatToolMapping.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 7, len(compatTools))

	s.CompatToolMapping.Add("1", "foo")

	compatTools, err = s.CompatToolMapping.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 8, len(compatTools))
	assert.NotEmpty(t, compatTools["foo"])
	if compatTools["foo"] != nil {
		assert.Equal(t, 1, len(compatTools["foo"].Games))
	}

	s.CompatToolMapping.Update("1", "bar")

	compatTools, err = s.CompatToolMapping.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 8, len(compatTools))
	assert.NotEmpty(t, compatTools["bar"])
	if compatTools["bar"] != nil {
		assert.Equal(t, 1, len(compatTools["bar"].Games))
	}
}
