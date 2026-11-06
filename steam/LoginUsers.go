package steam

import (
	"log/slog"
	"path/filepath"
	"strings"
)

// LoginUsersVdf represents parsed VDF config for app config from
// loginusers.vdf
type LoginUsersVdf struct {
	Vdf
}

// GetID64 return steamID64 for given username
func (lu *LoginUsersVdf) GetID64(username string) string {
	slog.Debug("LoginUsersVdf.GetID64", "username", username)

	if username == "" {
		return ""
	}

	for x := lu.Root.FirstSubTree(); x != nil; x = x.NextSubTree() {
		if strings.EqualFold(x.FirstByName("AccountName").String(), username) {
			return x.Name()
		}
	}

	return ""
}

// GetMostRecentID64 returns steamID64 of the user that logged in most recently
// (empty string if unknown)
func (lu *LoginUsersVdf) GetMostRecentID64() string {
	for x := lu.Root.FirstSubTree(); x != nil; x = x.NextSubTree() {
		if x.FirstByName("MostRecent").String() == "1" {
			return x.Name()
		}
	}

	return ""
}

func (s *Steam) initLoginUsers() error {
	p := filepath.Join(s.Root, "config", "loginusers.vdf")

	n, err := ParseTextConfig(p)
	if err != nil {
		return err
	}

	s.LoginUsers = &LoginUsersVdf{Vdf{n, nil, p, s}}

	return nil
}
