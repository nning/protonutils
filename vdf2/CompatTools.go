package vdf2

import (
	"sort"

	"github.com/nning/protonutils/steam"
)

// CompatTool holds info about a compatibility tool (like human-readable name
// and a list of the games that are using it)
type CompatTool struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Games     Games  `json:"games"`
	IsDefault bool   `json:"isDefault"`
	IsCustom  bool   `json:"isCustom"`
}

// CompatTools maps compatibility tool version IDs to CompatTool objects
// containing info like name and games.
type CompatTools map[string]*CompatTool

// IsValid checks whether a version ID (v) exists in the CompatTools config
func (c CompatTools) IsValid(v string) bool {
	for id := range c {
		if id == v {
			return true
		}
	}

	return false
}

// Add adds an entry to the CompatTools config (by tool id and (display) name)
func (c CompatTools) Add(id, name string) {
	if c[id] != nil {
		return
	}

	c[id] = &CompatTool{
		ID:        id,
		Name:      name,
		IsDefault: id == "",
		Games:     make(Games),
	}
}

// AddGame adds a Game entry to an existing CompatTool entry
func (c CompatTools) AddGame(toolID string, game *Game) bool {
	if c[toolID] == nil {
		return false
	}

	c[toolID].Games[game.Name] = game
	return true
}

func (c CompatTools) Read(s *steam.Steam) (CompatTools, error) {
	ctm, err := GetCompatToolMapping(s)
	if err != nil {
		return nil, err
	}

	tools, err := ctm.ReadCompatTools()
	if err != nil {
		return nil, err
	}

	c.Merge(&tools)

	lc, err := GetLocalConfig(s)
	if err != nil {
		return nil, err
	}

	games, err := lc.GetViewedSteamPlay()
	if err != nil {
		return nil, err
	}

	def := c.GetDefault()
	if def != nil {
		for _, game := range games {
			if c.Includes(game.ID) {
				continue
			}

			c.AddGame(def.ID, game)
		}
	}

	// fmt.Println(path.Join(s.Root, "compatibilitytools.d"))
	// files, err := ioutil.ReadDir(path.Join(s.Root, "compatibilitytools.d"))
	// if err != nil {
	// 	return nil, err
	// }

	// for _, file := range files {
	// 	id := file.Name()
	// 	if strings.HasPrefix(id, ".") {
	// 		continue
	// 	}
	// 	c.Add(id, id)
	// }

	return c, nil
}

// Merge adds all entries from other to c (without duplicates)
func (c CompatTools) Merge(other *CompatTools) CompatTools {
	for _, tool := range *other {
		for _, game := range tool.Games {
			c.Add(tool.ID, tool.Name)
			c.AddGame(tool.ID, game)
		}
	}

	return c
}

// GetDefault returns the first entry from c that is a default compatibility tool
func (c CompatTools) GetDefault() *CompatTool {
	for _, tool := range c {
		if tool.IsDefault {
			return tool
		}
	}

	return nil
}

// Sort returns slice of alphabetically sorted CompatTools IDs
func (c CompatTools) Sort() []string {
	type kv struct {
		key   string
		value *CompatTool
	}

	var tmp []kv
	for k, v := range c {
		tmp = append(tmp, kv{k, v})
	}

	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].value.Name < tmp[j].value.Name
	})

	var ids []string
	for _, tool := range tmp {
		ids = append(ids, tool.key)
	}

	return ids
}

// Includes returns whether a game is included in c
func (c CompatTools) Includes(appID string) bool {
	for _, tool := range c {
		if tool.Games.Includes(appID) {
			return true
		}
	}

	return false
}

// NewCompatTools returns new CompatTools struct (optionally initialized with
// data, which is a map from compatibility tool version IDs to slices of game
// IDs)
func NewCompatTools(s *steam.Steam, data ...map[string][]string) (*CompatTools, error) {
	compatTools := make(CompatTools)

	if len(data) > 0 {
		for versionID, games := range data[0] {
			for _, gameID := range games {
				game, _, err := GetGameData(s, gameID)
				if err != nil {
					return nil, err
				}

				compatTools.Add(versionID, s.GetCompatToolName(versionID))
				compatTools.AddGame(versionID, game)
			}
		}
	}

	return &compatTools, nil
}
