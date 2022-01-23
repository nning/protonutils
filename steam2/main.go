package steam2

import (
	"io/ioutil"

	"fmt"
	"io/fs"
	"os"
	osUser "os/user"
	"path"
	"strings"

	"github.com/BenLubar/vdf"

	"github.com/nning/protonutils/cache"
)

const testSteamRoot = "../test/root"

// Steam struct wraps caches and exposes functions for Steam data retrieval
type Steam struct {
	AppidCache       *cache.Cache
	VersionNameCache *cache.Cache

	AppInfo           *AppInfoVdf
	CompatToolMapping *CompatToolMappingVdf
	LibraryConfig     *LibraryConfigVdf
	LocalConfig       *LocalConfigVdf

	CompatTools CompatTools

	UID  string
	Root string
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
		AppidCache:       appidCache,
		VersionNameCache: protonNameCache,
		Root:             root,
	}

	s.CompatTools = make(CompatTools)

	uid, _ := s.userToID32(user)
	s.UID, err = s.getUID(uid)
	if err != nil {
		return nil, err
	}

	err = s.InitAppInfo()
	if err != nil {
		return nil, err
	}

	err = s.InitCompatToolMapping()
	if err != nil {
		return nil, err
	}

	err = s.InitLibraryConfig()
	if err != nil {
		return nil, err
	}

	err = s.InitLocalConfig()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// SaveCache writes caches to disk
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

// Lookup looks up a "key path" in a parsed VDF tree
func Lookup(n *vdf.Node, x []string) (*vdf.Node, error) {
	y := n

	for _, key := range x {
		y = y.FirstByName(key)
		if y == nil {
			return nil, &KeyNotFoundError{Name: key}
		}
	}

	return y, nil
}

// ParseTextConfig reads a file and parses it as text VDF
func ParseTextConfig(p string) (*vdf.Node, error) {
	in, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var n vdf.Node
	err = n.UnmarshalText(in)

	return &n, nil
}
