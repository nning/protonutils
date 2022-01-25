package steam2

import (
	"fmt"
	"io/fs"
	"os"
	osUser "os/user"
	"path"
	"strings"

	"github.com/nning/protonutils/cache"
)

const testSteamRoot = "../test/root"

// Steam struct wraps caches and exposes functions for Steam data retrieval
type Steam struct {
	AppidCache       *cache.Cache
	VersionNameCache *cache.Cache

	AppInfo           *BinaryVdf
	CompatToolMapping *CompatToolMappingVdf
	LibraryConfig     *LibraryConfigVdf
	LocalConfig       *LocalConfigVdf
	Shortcuts         *BinaryVdf

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

	err = s.initConfigs()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Steam) initConfigs() error {
	err := s.initAppInfo()
	if err != nil {
		return err
	}

	err = s.initCompatToolMapping()
	if err != nil {
		return err
	}

	err = s.initLibraryConfig()
	if err != nil {
		return err
	}

	err = s.initLocalConfig()
	if err != nil {
		return err
	}

	err = s.initShortcuts()
	if err != nil {
		return err
	}

	return nil
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

func (s *Steam) GetCompatibilityToolsDir() string {
	return path.Join(s.Root, "compatibilitytools.d")
}
