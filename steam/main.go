package steam

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"regexp"
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
}

func getUID(u string) (string, error) {
	usr, _ := user.Current()
	dir := path.Join(usr.HomeDir, ".steam", "root", "userdata")

	// TODO Sort entries by last change?
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	uid := entries[0].Name()

	if len(entries) > 1 {
		users := make([]string, 0)

		for _, entry := range entries {
			name := entry.Name()
			if name == u {
				return name, nil
			}

			isEntryNumeric, err := regexp.MatchString("^[0-9]*$", name)
			if err != nil {
				return "", err
			}

			if name != "0" && isEntryNumeric {
				users = append(users, name)
			}
		}

		uid = users[0]

		fmt.Fprintln(os.Stderr,
			"Warning: Several Steam users available, using "+uid+"\n"+
				"All available users: "+strings.Join(users, ", ")+"\n"+
				"Option \"-u\" can be used to specify user\n")
	}

	return uid, nil
}

// New instantiates Steam struct
func New(user string, ignoreCache bool) (*Steam, error) {
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

	s := &Steam{
		AppidCache:         appidCache,
		VersionNameCache:   protonNameCache,
		vdfCache:           make(mapLevel),
		CompatToolVersions: make(CompatToolVersions),
	}

	uid, _ := s.userToID32(user)
	s.uid, err = getUID(uid)
	if err != nil {
		return nil, err
	}

	return s, nil
}
