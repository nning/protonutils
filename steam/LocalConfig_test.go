package steam

import (
	"bytes"
	"encoding/binary"
	"strconv"
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

		// Position of the first match of the search used by older versions
		needle := []byte("appid\x00\x00\x00\x00\x00")
		appID, _ := strconv.ParseUint(id, 10, 64)
		binary.LittleEndian.PutUint32(needle[6:], uint32(appID))
		legacyPos := bytes.Index(s.AppInfo.data, needle)

		legacyName := ""
		if legacyPos >= 10 {
			if ln, err := ParseBinaryVdf(s.AppInfo.data[legacyPos-10:]); err == nil {
				legacyName = ln.FirstByName("common").FirstByName("name").String() + " (appid " + ln.FirstByName("appid").String() + ")"
			}
		}

		t.Logf("invalid: id=%v subtree=%v inAppInfo=%v getErr=%v apiName=%q apiErr=%v legacyPos=%v legacyName=%q",
			id, x.FirstChild() != nil || x.String() == "", n != nil, getErr, apiName, apiErr, legacyPos, legacyName)
	}
}
