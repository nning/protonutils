package steam

import (
	"github.com/nning/list_proton_versions/cache"
)

// Steam struct wraps caches and exposes functions for Steam data retrieval
type Steam struct {
	appidCache         *cache.Cache
	vdfCache           mapLevel
	CompatToolVersions CompatToolVersions
}

// New instantiates Steam struct
func New(notFake bool) (*Steam, error) {
	c, err := cache.New("steam-appids", notFake)
	if err != nil {
		return nil, err
	}

	return &Steam{
		appidCache:         c,
		vdfCache:           make(mapLevel),
		CompatToolVersions: make(CompatToolVersions),
	}, nil
}
