package steam

import (
	"github.com/nning/list_proton_versions/cache"
)

type Steam struct {
	cache         *cache.Cache
	libraryConfig map[string]interface{}
}

func New() *Steam {
	c := cache.New("steam-appids")
	return &Steam{c, nil}
}
