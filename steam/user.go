package steam

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/MrWaggel/gosteamconv"
)

func (s *Steam) getID64(u string) (string, error) {
	cfg, err := s.getLoginUsers()
	if err != nil {
		return "", err
	}

	for id, c := range cfg {
		if c.(mapLevel)["AccountName"] == u {
			return id, nil
		}
	}

	return "", errors.New("User not found: " + u)
}

func (s *Steam) userToID32(u string) (string, error) {
	idStr64, err := s.getID64(u)
	if err != nil {
		x, e := strconv.ParseInt(u, 10, 32)
		if e != nil {
			return "", err
		}

		_, e = gosteamconv.SteamInt32ToString(int32(x))
		if e != nil {
			return "", err
		}

		return u, nil
	}

	idInt64, err := strconv.ParseInt(idStr64, 10, 64)
	if err != nil {
		return "", err
	}

	str, err := gosteamconv.SteamInt64ToString(idInt64)
	if err != nil {
		return "", err
	}

	idInt32, err := gosteamconv.SteamStringToInt32(str)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(idInt32), nil
}
