package steam

import (
	"path"

	log "github.com/sirupsen/logrus"
)

// LoginUsersVdf represents parsed VDF config for app config from
// loginusers.vdf
type LoginUsersVdf struct {
	Vdf
}

// GetID64 return steamID64 for given username
func (lu *LoginUsersVdf) GetID64(username string) string {
	log.Debug("LoginUsersVdf.GetID64(", username, ")")

	x := lu.Root.FirstSubTree()

	for ; x != nil; x = x.NextChild() {
		if x.FirstByName("AccountName").String() == username {
			return x.Name()
		}
	}

	return ""
}

func (s *Steam) initLoginUsers() error {
	p := path.Join(s.Root, "config", "loginusers.vdf")

	n, err := ParseTextConfig(p)
	if err != nil {
		return err
	}

	s.LoginUsers = &LoginUsersVdf{Vdf{n, nil, p, s}}

	return nil
}
