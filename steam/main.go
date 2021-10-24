package steam

import (
	"github.com/nning/list_proton_versions/cache"
)

type Steam struct {
	cache              *cache.Cache
	libraryConfig      map[string]interface{}
	CompatToolVersions CompatToolVersions
}

func New(not_fake bool) (*Steam, error) {
	c, err := cache.New("steam-appids", not_fake)
	if err != nil {
		return nil, err
	}

	return &Steam{
		cache:              c,
		libraryConfig:      nil,
		CompatToolVersions: make(CompatToolVersions),
	}, nil
}
