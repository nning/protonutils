package steam

import (
	"errors"
	"strconv"
)

func (s *Steam) getName(id string) (string, bool, error) {
	name, valid := s.appidCache.Get(id)
	if name != "" && valid {
		return name, true, nil
	}

	name, err := s.findNameInAppInfo(id)
	if err != nil && !errors.Is(err, strconv.ErrRange) {
		return "", false, err
	}

	valid = true
	if name == "" {
		valid = false
	}

	s.appidCache.Add(id, name, valid)
	return name, valid, nil
}

func (s *Steam) getGameData(id string) (*gameData, error) {
	isInstalled, err := s.isInstalled(id)
	if err != nil {
		return nil, err
	}

	return &gameData{id, isInstalled}, nil
}

// SaveCache writes app ID cache to disk
func (s *Steam) SaveCache() error {
	return s.appidCache.Write()
}
