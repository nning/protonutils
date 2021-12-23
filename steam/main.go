package steam

import (
	"fmt"
	"io/fs"
	"os"
	osUser "os/user"
	"path"
	"strings"

	"github.com/nning/protonutils/cache"
)

// Steam struct wraps caches and exposes functions for Steam data retrieval
type Steam struct {
	AppidCache         *cache.Cache
	VersionNameCache   *cache.Cache
	vdfCache           mapLevel
	CompatToolVersions CompatToolVersions
	uid                string
	Root               string
}

// New instantiates Steam struct
func New(user string, root string, ignoreCache bool) (*Steam, error) {
	t := -1
	if ignoreCache {
		t = 0
	}

	appidCache, err := cache.New("appids", int64(t))
	if err != nil {
		return nil, err
	}

	t = 6 * 60 * 60 // 6h
	if ignoreCache {
		t = 1
	}

	protonNameCache, err := cache.New("proton-names", int64(t))
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(root, "~/") {
		usr, _ := osUser.Current()
		root = path.Join(usr.HomeDir, root[2:])
	}

	var fInfo fs.FileInfo
	fInfo, err = os.Stat(root)
	if err != nil || !fInfo.IsDir() {
		return nil, fmt.Errorf("Steam root not a directory: %v", root)
	}

	s := &Steam{
		AppidCache:         appidCache,
		VersionNameCache:   protonNameCache,
		vdfCache:           make(mapLevel),
		CompatToolVersions: make(CompatToolVersions),
		Root:               root,
	}

	uid, _ := s.userToID32(user)
	s.uid, err = s.getUID(uid)
	if err != nil {
		return nil, err
	}

	return s, nil
}
