package steam

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

	s, err := New("", testSteamRoot, true)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.CompatToolMapping)

	compatTools, err := s.CompatToolMapping.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 7, len(compatTools))

	for compatToolID, compatTool := range compatTools {
		assert.Equal(t, compatToolID, compatTool.ID)

		if compatTool.IsDefault {
			assert.Equal(t, "", compatTool.ID)
			// TODO
			// assert.Equal(t, "Proton 6.3-8 (Default)", compatTool.Name)
		} else {
			assert.NotEqual(t, "", compatTool.ID)
			assert.NotEqual(t, "", compatTool.Name)
		}

		// if compatTool.IsCustom {
		// 	assert.Equal(t, compatTool.ID, compatTool.Name)
		// }

		assert.NotEqual(t, 0, len(compatTool.Games))
	}

	compatTools, err = s.CompatToolMapping.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 7, len(compatTools))
}

func Test_Add(t *testing.T) {
	t.Parallel()

	s, err := New("", testSteamRoot, true)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.CompatToolMapping)

	ctm := s.CompatToolMapping
	assert.NotEmpty(t, ctm)

	id := "1593500"
	v := "Proton-8.3-GE-1"

	compatTools, err := ctm.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 7, len(compatTools))
	assert.Empty(t, compatTools[v])

	ctm.Add(id, v)

	compatTools, err = ctm.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 8, len(compatTools))
	assert.NotEmpty(t, compatTools[v])
	if compatTools[v] != nil {
		assert.Equal(t, 1, len(compatTools[v].Games))
	}
}

func Test_Update(t *testing.T) {
	t.Parallel()

	s, err := New("", testSteamRoot, true)
	assert.Empty(t, err)

	ctm := s.CompatToolMapping
	assert.NotEmpty(t, ctm)

	id := "1593500"
	v1 := "Proton-8.3-GE-1"
	v2 := "Proton-8.3-GE-2"

	compatTools, err := ctm.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 7, len(compatTools))
	assert.Empty(t, compatTools[v1])

	ctm.Add(id, v1)

	compatTools, err = ctm.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 8, len(compatTools))
	assert.NotEmpty(t, compatTools[v1])
	if compatTools[v1] != nil {
		assert.Equal(t, 1, len(compatTools[v1].Games))
	}

	ctm.Update(id, v2)

	compatTools, err = ctm.ReadCompatTools()
	assert.Empty(t, err)
	assert.Equal(t, 8, len(compatTools))
	assert.NotEmpty(t, compatTools[v2])
	if compatTools[v2] != nil {
		assert.Equal(t, 1, len(compatTools[v2].Games))
	}
}
