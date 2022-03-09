package steam

import (
	"path"
)

// LoginUsersVdf represents parsed VDF config for app config from
// loginusers.vdf
type LoginUsersVdf struct {
	Vdf
}

// GetID64 return steamID64 for given username
func (lu *LoginUsersVdf) GetID64(username string) string {
	x := lu.Root.FirstSubTree()

	for {
		if x.FirstByName("AccountName").String() == username {
			return x.Name()
		}

		x := x.NextSubTree()
		if x == nil {
			break
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
