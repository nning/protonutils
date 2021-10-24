package steam

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/MrWaggel/gosteamconv"
)

func (s *Steam) GetId64(u string) (string, error) {
	cfg, err := s.GetLoginUsers()
	if err != nil {
		return "", err
	}

	for id, c := range cfg {
		if c.(MapLevel)["AccountName"] == u {
			return id, nil
		}
	}

	return "", errors.New("User not found: " + u)
}

func (s *Steam) UserToId32(u string) (string, error) {
	id_str64, err := s.GetId64(u)
	if err != nil {
		x, e := strconv.ParseInt(u, 10, 32)
		if e != nil {
			return "", err
		}

		_, e = gosteamconv.SteamInt32ToString(int32(x))
		if e != nil {
			return "", err
		} else {
			return u, nil
		}
	}

	id_int64, err := strconv.ParseInt(id_str64, 10, 64)
	if err != nil {
		return "", err
	}

	str, err := gosteamconv.SteamInt64ToString(id_int64)
	if err != nil {
		return "", err
	}

	id_int32, err := gosteamconv.SteamStringToInt32(str)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(id_int32), nil
}
