package steam

import (
	"io"
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

// GetGameData returns intialized Game struct by app ID
func (s *Steam) GetGameData(id string) (*Game, bool, error) {
	var err error

	sc := isShortcut(id)
	name, valid := s.AppidCache.Get(id)
	if name != "" && valid {
		data, err := s.getGameData(id, name, sc)
		if err != nil {
			return nil, false, err
		}

		return data, true, nil
	}

	if sc {
		name, err = s.findNameInShortcuts(id)
	} else {
		name, err = s.findNameInAppInfo(id)
		if err == io.EOF {
			name, err = s.getNameFromAPI(id)
		}
		if err != nil {
			return nil, false, err
		}
	}

	valid = true
	if name == "" {
		valid = false
	}

	data, err := s.getGameData(id, name, sc)
	if err != nil {
		return nil, valid, err
	}

	s.AppidCache.Add(id, name, valid)
	return data, valid, nil
}

func (s *Steam) getGameData(id string, name string, isShortcut bool) (*Game, error) {
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

	return &Game{id, name, isInstalled, isShortcut}, nil
}

// SaveCache writes app ID cache to disk
func (s *Steam) SaveCache() error {
	err := s.AppidCache.Write()
	if err != nil {
		return err
	}

	err = s.VersionNameCache.Write()
	if err != nil {
		return err
	}

	return nil
}
