package steam

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"testing"

	"github.com/BenLubar/vdf"
	"github.com/stretchr/testify/assert"
)

func Test_GetGames_EnabledViewedSteamPlay(t *testing.T) {
	requireTestRoot(t)
	s, _ := New("", testConfig, false)

	games, err := s.LocalConfig.GetGames(true)
	assert.Empty(t, err)
	assert.Equal(t, 36, len(games))
}

func Test_GetGames_DisabledViewedSteamPlay(t *testing.T) {
	requireTestRoot(t)
	s, _ := New("", testConfig, false)

	games, err := s.LocalConfig.GetGames(false)
	assert.Empty(t, err)
	// Compared to older versions: +1 Half-Life 2 (dropped because the byte
	// search hit a wrong entry), -4 config apps (Steam Client, Steam
	// Screenshots, Steam Input Configs, EALink)
	if !assert.Equal(t, 242, len(games)) {
		logInvalidGames(t, s, games)
		logGamesWithoutLegacyName(t, s, games)
	}
}

// logInvalidGames logs diagnostics for app IDs in localconfig.vdf that are not
// returned by GetGames
func logInvalidGames(t *testing.T, s *Steam, games []*Game) {
	t.Helper()

	valid := make(map[string]bool)
	for _, g := range games {
		valid[g.ID] = true
	}

	for x := s.LocalConfig.Node.FirstChild(); x != nil; x = x.NextChild() {
		id := x.Name()
		if valid[id] {
			continue
		}

		n, getErr := s.AppInfo.Get(id)
		apiName, apiErr := s.GetNameFromAPI(id)

		legacyPos, legacyName, _ := legacyLookup(s, id)

		t.Logf("not listed: id=%v type=%q inAppInfo=%v getErr=%v apiName=%q apiErr=%v legacyPos=%v legacyName=%q",
			id, s.AppInfo.GetType(id), n != nil, getErr, apiName, apiErr, legacyPos, legacyName)
	}
}

// legacyLookup returns position, name and parse error of the lookup used by
// older versions (byte search, parsing with BenLubar/vdf until the end of
// appinfo.vdf, errors ignored)
func legacyLookup(s *Steam, id string) (int, string, error) {
	needle := []byte("appid\x00\x00\x00\x00\x00")
	appID, _ := strconv.ParseUint(id, 10, 64)
	binary.LittleEndian.PutUint32(needle[6:], uint32(appID))
	pos := bytes.Index(s.AppInfo.data, needle)
	if pos < 10 {
		return pos, "", nil
	}

	var n vdf.Node
	err := n.UnmarshalBinary(s.AppInfo.data[pos-10:])
	return pos, n.FirstByName("common").FirstByName("name").String(), err
}

// logGamesWithoutLegacyName logs games that have a name now, but had none with
// the lookup used by older versions (which did not query the store API if an
// appinfo.vdf entry was found)
func logGamesWithoutLegacyName(t *testing.T, s *Steam, games []*Game) {
	t.Helper()

	for _, g := range games {
		if IsShortcut(g.ID) {
			continue
		}

		pos, legacyName, legacyErr := legacyLookup(s, g.ID)
		if pos < 0 || legacyName != "" {
			continue
		}

		n, _ := s.AppInfo.Get(g.ID)
		t.Logf("named now: id=%v name=%q appinfoName=%q legacyPos=%v legacyErr=%v",
			g.ID, g.Name, n.FirstByName("common").FirstByName("name").String(), pos, legacyErr)
	}
}
