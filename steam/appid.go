package steam

import (
	"strconv"
)

func isShortcut(id string) bool {
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		_, err := strconv.ParseInt(id, 10, 64)
		if err == nil {
			return true
		}
	}

	return false
}

func (s *Steam) getNameAndGameData(id string) (string, *gameData, bool, error) {
	var err error

	sc := isShortcut(id)
	name, valid := s.appidCache.Get(id)
	if name != "" && valid {
		data, err := s.getGameData(id, sc)
		if err != nil {
			return "", nil, false, err
		}

		return name, data, true, nil
	}

	if sc {
		name, err = s.findNameInShortcuts(id)
	} else {
		name, err = s.findNameInAppInfo(id)
		if err != nil {
			return "", nil, false, err
		}
	}

	valid = true
	if name == "" {
		valid = false
	}

	data, err := s.getGameData(id, sc)
	if err != nil {
		return name, nil, valid, err
	}

	s.appidCache.Add(id, name, valid)
	return name, data, valid, nil
}

func (s *Steam) getGameData(id string, isShortcut bool) (*gameData, error) {
	var isInstalled bool
	var err error

	if isShortcut {
		isInstalled = true
	} else {
		isInstalled, err = s.isInstalled(id)
		if err != nil {
			return nil, err
		}
	}

	return &gameData{id, isInstalled, isShortcut}, nil
}

// SaveCache writes app ID cache to disk
func (s *Steam) SaveCache() error {
	err := s.appidCache.Write()
	if err != nil {
		return err
	}

	err = s.versionNameCache.Write()
	if err != nil {
		return err
	}

	return nil
}
