package steam

import (
	"github.com/nning/list_proton_versions/cache"
)

type Steam struct {
	appidCache         *cache.Cache
	vdfCache           MapLevel
	CompatToolVersions CompatToolVersions
}

func New(notFake bool) (*Steam, error) {
	c, err := cache.New("steam-appids", notFake)
	if err != nil {
		return nil, err
	}

	return &Steam{
		appidCache:         c,
		vdfCache:           make(MapLevel),
		CompatToolVersions: make(CompatToolVersions),
	}, nil
}
