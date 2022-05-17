package steam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CompatTools_Add(t *testing.T) {
	compatTools := make(CompatTools)
	version := "proton_63"
	id := "620"

	s, err := New("", testSteamRoot, true)
	assert.Empty(t, err)

	game, _, err := s.GetGameData(id)
	assert.Empty(t, err)

	versionName, err := s.GetCompatToolName(version)
	assert.Empty(t, err)
	assert.Equal(t, "Proton 6.3-8", versionName)

	compatTools.Add(version, versionName)
	compatTools.AddGame(version, game)

	v := compatTools[version]
	assert.NotEmpty(t, v)
	assert.Equal(t, version, v.ID)
	assert.Equal(t, "Proton 6.3-8", v.Name)
	assert.Equal(t, 1, len(v.Games))

	g := v.Games["Portal 2"]
	assert.NotEmpty(t, g)
	assert.Equal(t, id, g.ID)
	assert.Equal(t, "Portal 2", g.Name)
	assert.Equal(t, false, g.IsInstalled)
	assert.Equal(t, false, g.IsShortcut)
}

func Test_CompatTools_Merge(t *testing.T) {
	data1 := map[string][]string{
		"proton_63": {
			"620", // Portal 2
		},
	}

	data2 := map[string][]string{
		"proton_63": {
			"70", // Half-Life
		},
		"proton_experimental": {
			"400", // Portal
		},
	}

	s, err := New("", testSteamRoot, true)
	assert.Empty(t, err)

	compatTools1, _ := s.NewCompatTools(data1)
	compatTools2, _ := s.NewCompatTools(data2)

	compatTools1.Merge(compatTools2)
	x := *compatTools1

	assert.NotEmpty(t, x["proton_63"])
	assert.NotEmpty(t, x["proton_experimental"])
	assert.Equal(t, 2, len(x["proton_63"].Games))
	assert.Equal(t, 1, len(x["proton_experimental"].Games))

	assert.Equal(t, "620", x["proton_63"].Games["Portal 2"].ID)
	assert.Equal(t, "Portal 2", x["proton_63"].Games["Portal 2"].Name)

	assert.Equal(t, "70", x["proton_63"].Games["Half-Life"].ID)
	assert.Equal(t, "Half-Life", x["proton_63"].Games["Half-Life"].Name)

	assert.Equal(t, "400", x["proton_experimental"].Games["Portal"].ID)
	assert.Equal(t, "Portal", x["proton_experimental"].Games["Portal"].Name)
}

func Test_CompatTools_Read(t *testing.T) {
	s, err := New("", testSteamRoot, true)
	assert.Empty(t, err)

	err = s.ReadCompatTools()
	assert.Empty(t, err)

	compatTools := s.CompatTools

	assert.Equal(t, 7, len(compatTools))
	assert.Equal(t, 45, len(compatTools[""].Games))
	assert.Equal(t, 2, len(compatTools["proton_63"].Games))
	assert.Equal(t, 1, len(compatTools["proton_experimental"].Games))
	assert.Equal(t, 2, len(compatTools["Proton-7.0rc6-GE-1"].Games))
	assert.Equal(t, 1, len(compatTools["Proton-6.20-GE-1"].Games))
}
