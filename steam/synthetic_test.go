package steam

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nning/protonutils/config"
	"github.com/nning/protonutils/internal/steamtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupOffline writes a synthetic Steam root and makes the store API
// unreachable for the duration of the test (tests using the encrypted real data
// still query the real API)
func setupOffline(t *testing.T, r steamtest.Root) string {
	t.Helper()

	old := StoreAPIURL
	StoreAPIURL = "http://127.0.0.1:1/"
	t.Cleanup(func() { StoreAPIURL = old })

	return steamtest.Setup(t, r)
}

func newSynthetic(t *testing.T, r steamtest.Root, user string) *Steam {
	t.Helper()

	root := setupOffline(t, r)
	s, err := New(user, &config.Config{SteamRoot: root}, false)
	require.NoError(t, err)

	return s
}

func Test_Synthetic_AppInfoVersions(t *testing.T) {
	for _, version := range []int{27, 28, 29} {
		r := steamtest.Default()
		r.AppInfoVersion = version
		s := newSynthetic(t, r, "")

		assert.Equal(t, uint32(0x07564400+version/10*16+version%10), s.AppInfo.Version)
		assert.Equal(t, 6, s.AppInfo.Len())

		name, err := s.GetName("620")
		assert.NoError(t, err)
		assert.Equal(t, "Portal 2", name, version)

		dir, err := s.GetInstalldir("1091500")
		assert.NoError(t, err)
		assert.Equal(t, "Cyberpunk 2077", dir)

		dir, err = s.GetInstalldir("11111111")
		assert.NoError(t, err)
		assert.Equal(t, "", dir)

		game, err := s.GetGame("1091500")
		assert.NoError(t, err)
		assert.Equal(t, DeckCompatibilityPlayable, game.DeckCompatibility.Category)

		name, err = s.GetCompatToolName("proton_63")
		assert.NoError(t, err)
		assert.Equal(t, "Proton 6.3-8", name)
	}
}

func Test_Synthetic_ParseAppInfoErrors(t *testing.T) {
	_, err := ParseAppInfo([]byte{1, 2, 3, 4, 5, 6, 7, 8}, "x")
	assert.ErrorContains(t, err, "unsupported version")

	data := steamtest.AppInfo(28, steamtest.App{ID: 10, Name: "x"})
	_, err = ParseAppInfo(data[:len(data)-10], "x")
	assert.Error(t, err)

	a, err := ParseAppInfo(nil, "x")
	assert.NoError(t, err)
	assert.Equal(t, 0, a.Len())
}

func Test_Synthetic_ParseBinaryVdf(t *testing.T) {
	data := steamtest.BinaryVdf(steamtest.Tree("a",
		steamtest.Str("s", "value"),
		steamtest.Int("i", -5),
		steamtest.KV{Key: "u", Value: uint64(1) << 40},
		steamtest.Tree("empty"),
	))

	n, err := ParseBinaryVdf(data)
	require.NoError(t, err)
	assert.Equal(t, "a", n.Name())
	assert.Equal(t, "value", n.FirstByName("s").String())
	assert.Equal(t, int32(-5), n.FirstByName("i").Int())
	assert.Equal(t, uint64(1)<<40, n.FirstByName("u").Uint64())
	assert.NotNil(t, n.FirstByName("empty"))

	for i := 1; i < len(data)-1; i++ {
		_, err = ParseBinaryVdf(data[:i])
		assert.Error(t, err, i)
	}

	_, err = ParseBinaryVdf([]byte{0x42, 'x', 0})
	assert.ErrorContains(t, err, "unknown type")
}

func Test_Synthetic_Users(t *testing.T) {
	r := steamtest.Default()

	// Most recent user is chosen by default
	s := newSynthetic(t, r, "")
	assert.Equal(t, r.Users[0].ID32(), s.UID)

	for _, u := range []string{"bob", "BOB", r.Users[1].ID32(), "76561198050517828"} {
		root := setupOffline(t, r)
		s, err := New(u, &config.Config{SteamRoot: root}, false)
		require.NoError(t, err, u)
		assert.Equal(t, r.Users[1].ID32(), s.UID, u)
	}

	root := setupOffline(t, r)
	_, err := New("mallory", &config.Config{SteamRoot: root}, false)
	assert.ErrorContains(t, err, "Steam user not found: mallory")

	_, err = New("12345", &config.Config{SteamRoot: root}, false)
	assert.ErrorContains(t, err, "No userdata found")

	// Empty userdata directory must not panic
	r.Users = nil
	root = setupOffline(t, r)
	_, err = New("", &config.Config{SteamRoot: root}, false)
	assert.ErrorContains(t, err, "No Steam user found")
}

func Test_Synthetic_UsersWithoutMostRecent(t *testing.T) {
	r := steamtest.Default()
	r.Users[0].MostRecent = false

	root := setupOffline(t, r)
	s, err := New("", &config.Config{SteamRoot: root}, false)
	require.NoError(t, err)
	assert.Equal(t, r.Users[0].ID32(), s.UID)
}

func Test_Synthetic_FindRoot(t *testing.T) {
	root := setupOffline(t, steamtest.Default())
	home := os.Getenv("HOME")

	_, err := FindRoot("")
	assert.ErrorContains(t, err, "Steam root not found")

	flatpakRoot := filepath.Join(home, ".var/app/com.valvesoftware.Steam/data/Steam")
	require.NoError(t, os.MkdirAll(filepath.Dir(flatpakRoot), 0700))
	require.NoError(t, os.Symlink(root, flatpakRoot))

	for _, configured := range []string{"", LegacyDefaultRoot} {
		found, err := FindRoot(configured)
		assert.NoError(t, err)
		assert.Equal(t, flatpakRoot, found)
	}

	found, err := FindRoot("~/.var/app/com.valvesoftware.Steam/data/Steam")
	assert.NoError(t, err)
	assert.Equal(t, flatpakRoot, found)

	_, err = FindRoot("/nonexistent")
	assert.ErrorContains(t, err, "not a directory")
}

func Test_Synthetic_ReadCompatTools(t *testing.T) {
	s := newSynthetic(t, steamtest.Default(), "")
	require.NoError(t, s.ReadCompatTools())

	tools := s.CompatTools

	def := tools[""]
	require.NotNil(t, def)
	assert.True(t, def.IsDefault)
	assert.Equal(t, "GE-Proton9-1 (Default)", def.Name)
	// Proton 6.3 (1580130) is a tool and not listed
	assert.Equal(t, []string{"1091500"}, def.Games.Sort())
	assert.True(t, s.AppInfo.IsTool("1580130"))
	assert.False(t, s.AppInfo.IsTool("1091500"))
	assert.False(t, s.AppInfo.IsTool("70"))

	assert.Equal(t, "Proton 6.3-8", tools["proton_63"].Name)
	assert.False(t, tools["proton_63"].IsCustom)
	assert.Equal(t, []string{"3228583970", "620"}, tools["proton_63"].Games.Sort()) // sorted by name
	assert.Equal(t, "My Game", tools["proton_63"].Games["3228583970"].Name)
	assert.True(t, tools["proton_63"].Games["3228583970"].IsShortcut)

	for _, id := range []string{"GE-Proton6-1", "GE-Proton7-1", "GE-Proton8-1", "GE-Proton9-1"} {
		require.NotNil(t, tools[id], id)
		assert.True(t, tools[id].IsCustom, id)
		assert.True(t, tools[id].IsInstalled(s), id)
	}

	assert.False(t, tools["GE-Proton7-1"].Games["400"].IsInstalled)
	assert.Equal(t, "GE-Proton8-1", s.GetGameVersion("70").ID)
	assert.Nil(t, s.GetGameVersion("11111111"))
}

func Test_Synthetic_UnusedCustomToolDirs(t *testing.T) {
	r := steamtest.Default()
	// Second tool in directory of used tool
	r.Tools = append(r.Tools, steamtest.Tool{Dir: "GE-Proton8-1", ID: "GE-Proton8-1-extra", DisplayName: "Extra"})
	// Directory name differs from internal name
	r.Tools = append(r.Tools, steamtest.Tool{Dir: "renamed", ID: "proton-custom", DisplayName: "Custom"})
	r.Mapping["1091500"] = "proton-custom"

	s := newSynthetic(t, r, "")
	require.NoError(t, s.ReadCompatTools())

	assert.Equal(t, "Custom", s.CompatTools["proton-custom"].Name)

	var names []string
	for _, dir := range s.UnusedCustomToolDirs() {
		names = append(names, filepath.Base(dir))
	}

	// GE-Proton9-1 is the default, GE-Proton8-1 and renamed are used by
	// installed games, GE-Proton7-1 is only used by a game that is not
	// installed
	assert.Equal(t, []string{"GE-Proton6-1", "GE-Proton7-1"}, names)
}

func Test_Synthetic_GetAppIDAndName(t *testing.T) {
	s := newSynthetic(t, steamtest.Default(), "")
	require.NoError(t, s.ReadCompatTools())

	id, name, err := s.GetAppIDAndName("portal")
	assert.NoError(t, err)
	assert.Equal(t, "400", id)
	assert.Equal(t, "Portal", name)

	id, _, err = s.GetAppIDAndName("Portal 2")
	assert.NoError(t, err)
	assert.Equal(t, "620", id)

	id, _, err = s.GetAppIDAndName("my game")
	assert.NoError(t, err)
	assert.Equal(t, "3228583970", id)

	_, _, err = s.GetAppIDAndName("p")
	var ambiguous *AmbiguousNameError
	assert.ErrorAs(t, err, &ambiguous)
	assert.Len(t, ambiguous.Matches, 2)

	_, _, err = s.GetAppIDAndName("nothing")
	assert.ErrorContains(t, err, "not found")

	p, name, err := s.GetInstalldirPath("portal 2")
	assert.NoError(t, err)
	assert.Equal(t, "Portal 2", name)
	assert.True(t, strings.HasSuffix(p, "library/steamapps/common/Portal 2"), p)

	_, _, err = s.GetCompatdataPath("portal")
	assert.ErrorContains(t, err, "not installed")
}

func Test_Synthetic_UpdateAndSave(t *testing.T) {
	s := newSynthetic(t, steamtest.Default(), "")
	ctm := s.CompatToolMapping

	before, err := os.ReadFile(ctm.Path)
	require.NoError(t, err)

	require.NoError(t, ctm.Update("620", "GE-Proton8-1"))
	require.NoError(t, ctm.Update("1091500", "proton_63"))
	require.NoError(t, ctm.Update("70", ""))
	assert.Error(t, ctm.Update("0", ""))
	require.NoError(t, ctm.Save())

	backup, err := os.ReadFile(ctm.Path + BackupSuffix)
	require.NoError(t, err)
	assert.Equal(t, before, backup)

	entries, err := os.ReadDir(filepath.Dir(ctm.Path))
	require.NoError(t, err)
	for _, e := range entries {
		assert.NotContains(t, e.Name(), ".tmp", "temporary file left behind")
	}

	s2, err := New("", &config.Config{SteamRoot: s.Root}, false)
	require.NoError(t, err)
	require.NoError(t, s2.ReadCompatTools())

	assert.Equal(t, "GE-Proton8-1", s2.GetGameVersion("620").ID)
	assert.Equal(t, "proton_63", s2.GetGameVersion("1091500").ID)
	assert.True(t, s2.GetGameVersion("70").IsDefault)
	assert.Equal(t, "GE-Proton9-1", s2.CompatToolMapping.GetDefaultID())
}

func Test_Synthetic_ParseTextConfigError(t *testing.T) {
	p := filepath.Join(t.TempDir(), "broken.vdf")
	require.NoError(t, os.WriteFile(p, []byte("\"a\"\n{\n\t\"b\"\t\"c\"\n"), 0600))

	_, err := ParseTextConfig(p)
	assert.ErrorContains(t, err, "broken.vdf")

	r := steamtest.Default()
	root := setupOffline(t, r)
	require.NoError(t, os.WriteFile(filepath.Join(root, "config", "config.vdf"), []byte("\"x\"\n{\n"), 0600))

	_, err = New("", &config.Config{SteamRoot: root}, false)
	assert.ErrorContains(t, err, "config.vdf")
}

func Test_Synthetic_MissingCompatToolMapping(t *testing.T) {
	root := setupOffline(t, steamtest.Default())
	require.NoError(t, os.WriteFile(filepath.Join(root, "config", "config.vdf"), []byte("\"InstallConfigStore\"\n{\n}\n"), 0600))

	s, err := New("", &config.Config{SteamRoot: root}, false)
	require.NoError(t, err)
	require.NoError(t, s.ReadCompatTools())
	assert.Equal(t, "Default", s.CompatTools[""].Name)

	require.NoError(t, s.CompatToolMapping.Update("620", "proton_63"))
	require.NoError(t, s.CompatToolMapping.Save())

	s2, err := New("", &config.Config{SteamRoot: root}, false)
	require.NoError(t, err)
	assert.Equal(t, "proton_63", s2.CompatToolMapping.Node.FirstByName("620").FirstByName("name").String())
}

func Test_Synthetic_IsShortcut(t *testing.T) {
	assert.True(t, IsShortcut("3228583970"))
	assert.False(t, IsShortcut("620"))
	assert.False(t, IsShortcut("99999999999"))
	assert.False(t, IsShortcut("foo"))
}

func Test_Synthetic_UnsupportedAppInfoIsNotFatal(t *testing.T) {
	root := setupOffline(t, steamtest.Default())
	require.NoError(t, os.WriteFile(filepath.Join(root, "appcache", "appinfo.vdf"), []byte{0x30, 0x44, 0x56, 0x07, 1, 0, 0, 0}, 0600))

	s, err := New("", &config.Config{SteamRoot: root}, false)
	require.NoError(t, err)
	assert.Equal(t, 0, s.AppInfo.Len())
	require.NoError(t, s.ReadCompatTools())
}

func Test_Synthetic_CorruptShortcutsIsNotFatal(t *testing.T) {
	r := steamtest.Default()
	root := setupOffline(t, r)
	p := filepath.Join(root, "userdata", r.Users[0].ID32(), "config", "shortcuts.vdf")
	require.NoError(t, os.WriteFile(p, []byte{0x00, 's', 0, 0x42}, 0600))

	s, err := New("", &config.Config{SteamRoot: root}, false)
	require.NoError(t, err)
	assert.Equal(t, "", s.Shortcuts.GetName("3228583970"))
	require.NoError(t, s.ReadCompatTools())
}

func Test_Synthetic_BrokenAppInfoEntryIsNotFatal(t *testing.T) {
	r := steamtest.Default()
	// Unknown value type in one entry
	r.Apps[0].Extra = []steamtest.KV{{Key: "x", Value: uint64(1)}}
	root := setupOffline(t, r)

	p := filepath.Join(root, "appcache", "appinfo.vdf")
	data, err := os.ReadFile(p)
	require.NoError(t, err)
	// Replace the type byte of the uint64 node (0x07) of the first entry
	i := strings.Index(string(data), "\x07")
	for ; i >= 0 && i < len(data); i++ {
		if data[i] == 0x07 && i > 16 {
			data[i] = 0x42
			break
		}
	}
	require.NoError(t, os.WriteFile(p, data, 0600))

	s, err := New("", &config.Config{SteamRoot: root}, false)
	require.NoError(t, err)

	_, err = s.AppInfo.Get("70")
	assert.Error(t, err)

	require.NoError(t, s.ReadCompatTools())
	assert.Equal(t, "Portal 2", s.CompatTools["proton_63"].Games["620"].Name)

	n, err := s.AppInfo.Get("99999999999")
	assert.NoError(t, err)
	assert.Nil(t, n)
}

func Test_Synthetic_UsedToolIDsWithoutNames(t *testing.T) {
	r := steamtest.Default()
	// Shortcut of another user (not in alice's shortcuts.vdf)
	r.Mapping["3000000000"] = "GE-Proton6-1"
	// Installed game without appinfo.vdf entry (name lookup fails offline)
	r.Mapping["999999"] = "GE-Proton7-1"
	r.Installed = append(r.Installed, "999999")

	s := newSynthetic(t, r, "")
	require.NoError(t, s.ReadCompatTools())

	used := s.UsedToolIDs()
	assert.True(t, used["GE-Proton6-1"])
	assert.True(t, used["GE-Proton7-1"])
	assert.Empty(t, s.UnusedCustomToolDirs())
}

func Test_Synthetic_SaveThroughSymlink(t *testing.T) {
	s := newSynthetic(t, steamtest.Default(), "")
	ctm := s.CompatToolMapping

	real := filepath.Join(t.TempDir(), "config.vdf")
	data, err := os.ReadFile(ctm.Path)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(real, data, 0600))
	require.NoError(t, os.Remove(ctm.Path))
	require.NoError(t, os.Symlink(real, ctm.Path))

	require.NoError(t, ctm.Update("620", "GE-Proton6-1"))
	require.NoError(t, ctm.Save())

	info, err := os.Lstat(ctm.Path)
	require.NoError(t, err)
	assert.True(t, info.Mode()&os.ModeSymlink != 0, "symlink was replaced")

	data, err = os.ReadFile(real)
	require.NoError(t, err)
	assert.Contains(t, string(data), "GE-Proton6-1")
}
