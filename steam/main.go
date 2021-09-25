package steam

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"

	"github.com/andygrunwald/vdf"

	. "github.com/nning/list_proton_versions"
	"github.com/nning/list_proton_versions/cache"
)

type Steam struct {
	cache         *cache.Cache
	libraryConfig map[string]interface{}
}

type JsonAppData struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

func lookup(m map[string]interface{}, x []string) (map[string]interface{}, error) {
	y := m

	for _, s := range x {
		if y[s] == nil {
			return nil, errors.New("Key not found: " + s)
		} else {
			y = y[s].(map[string]interface{})
		}
	}

	return y, nil
}

func vdfLookup(file string, x ...string) (map[string]interface{}, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	file = path.Join(dir, ".steam", "root", file)

	f, err := os.Open(file)
	PanicOnError(err)

	p := vdf.NewParser(f)
	m, err := p.Parse()
	PanicOnError(err)

	return lookup(m, x)
}

func New() *Steam {
	c := cache.New("steam-appids")
	return &Steam{
		cache: c,
	}
}

func (self *Steam) GetName(id string) string {
	name := self.cache.Get(id)

	if name != "" {
		return name
	}

	res, err := http.Get("https://store.steampowered.com/api/appdetails/?appids=" + id)
	PanicOnError(err)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	PanicOnError(err)

	data := make(map[string]JsonAppData)
	err = json.Unmarshal(body, &data)
	PanicOnError(err)

	name = data[id].Data.Name
	val := name
	if val == "" {
		val = "💩"
	}

	self.cache.Add(id, val)

	return name
}

func (self *Steam) SaveCache() {
	self.cache.Write()
}

func (self *Steam) GetCompatToolMapping(x ...string) (map[string]interface{}, error) {
	return vdfLookup("config/config.vdf", "InstallConfigStore", "Software", "Valve", "Steam", "CompatToolMapping")
}

func (self *Steam) GetLibraryConfig() (map[string]interface{}, error) {
	return vdfLookup("steamapps/libraryfolders.vdf", "libraryfolders")
}

func (self *Steam) IsInstalled(id string) bool {
	m := self.libraryConfig
	var err error

	if m == nil {
		m, err = self.GetLibraryConfig()
		PanicOnError(err)
		self.libraryConfig = m
	}

	installed := false

	for i := 0; i < 10; i++ {
		x := m[fmt.Sprint(i)]
		if x == nil {
			break
		}

		apps := x.(map[string]interface{})["apps"].(map[string]interface{})
		for app := range apps {
			if app == id {
				return true
			}
		}
	}

	return installed
}
