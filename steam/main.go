package steam

import (
	"github.com/nning/list_proton_versions/cache"
)

type Steam struct {
	cache              *cache.Cache
	libraryConfig      map[string]interface{}
	CompatToolVersions CompatToolVersions
}

func New() *Steam {
	c := cache.New("steam-appids")

	return &Steam{
		cache:              c,
		libraryConfig:      nil,
		CompatToolVersions: make(CompatToolVersions),
	}
}
