package steam

import (
	"errors"
	"strconv"
)

// InvalidID is a placeholder for cached entries with invalid app ID
// TODO Transparently cache invalid app IDs
const InvalidID = "💩"

// GetName returns name for app ID
func (s *Steam) GetName(id string) (string, error) {
	name := s.appidCache.Get(id)
	if name != "" {
		return name, nil
	}

	name, err := s.findNameInAppInfo(id)
	if err != nil && !errors.Is(err, strconv.ErrRange) {
		return "", err
	}

	if name == "" {
		name = InvalidID
	}

	s.appidCache.Add(id, name)
	return name, nil
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
