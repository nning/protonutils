package steam

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"
	"testing"

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
	if !assert.Equal(t, 245, len(games)) {
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

		legacyPos, legacyName := legacyLookup(s, id)

		t.Logf("invalid: id=%v subtree=%v inAppInfo=%v getErr=%v apiName=%q apiErr=%v legacyPos=%v legacyName=%q",
			id, x.FirstChild() != nil || x.String() == "", n != nil, getErr, apiName, apiErr, legacyPos, legacyName)
	}
}

// legacyLookup returns position and name (with app ID of the entry) found by
// the byte search used by older versions
func legacyLookup(s *Steam, id string) (int, string) {
	needle := []byte("appid\x00\x00\x00\x00\x00")
	appID, _ := strconv.ParseUint(id, 10, 64)
	binary.LittleEndian.PutUint32(needle[6:], uint32(appID))
	pos := bytes.Index(s.AppInfo.data, needle)

	name := ""
	if pos >= 10 {
		if n, err := ParseBinaryVdf(s.AppInfo.data[pos-10:]); err == nil {
			name = n.FirstByName("common").FirstByName("name").String() + " (appid " + n.FirstByName("appid").String() + ")"
		}
	}

	return pos, name
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

		pos, legacyName := legacyLookup(s, g.ID)
		if pos >= 0 && !strings.HasPrefix(legacyName, " (") {
			continue
		}

		n, _ := s.AppInfo.Get(g.ID)
		t.Logf("named now: id=%v name=%q appinfoName=%q legacyPos=%v legacyName=%q",
			g.ID, g.Name, n.FirstByName("common").FirstByName("name").String(), pos, legacyName)
	}
}
