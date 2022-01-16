package vdf2

import (
	"fmt"

	"github.com/nning/protonutils/steam"
)

// CompatTools maps compatibility tool version IDs to CompatTool objects
// containing info like name and games.
type CompatTools map[string]*CompatTool

// CompatTool holds info about a compatibility tool (like human-readable name
// and a list of the games that are using it)
type CompatTool = steam.Version

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
		Games:     make(steam.Games),
	}
}

// AddGame adds a Game entry to an existing CompatTool entry
func (c CompatTools) AddGame(toolID string, game *steam.Game) bool {
	if c[toolID] == nil {
		return false
	}

	c[toolID].Games[game.Name] = game
	return true
}

func (c CompatTools) Read(s *steam.Steam) (*CompatTools, error) {
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

	return &c, nil
}

func (c CompatTools) Merge(other *CompatTools) CompatTools {
	for _, tool := range *other {
		for _, game := range tool.Games {
			c.Add(tool.ID, tool.Name)
			c.AddGame(tool.ID, game)
		}
	}

	return c
}

func (c CompatTools) GetDefault() *CompatTool {
	fmt.Println(len(c)) // TODO c empty?!

	for _, tool := range c {
		if tool.IsDefault {
			return tool
		}
	}

	return nil
}

func NewCompatTools(s *steam.Steam, data map[string][]string) (*CompatTools, error) {
	compatTools := make(CompatTools)

	for versionID, games := range data {
		for _, gameID := range games {
			game, _, err := s.GetGameData(gameID)
			if err != nil {
				return nil, err
			}

			compatTools.Add(versionID, s.GetCompatToolName(versionID))
			compatTools.AddGame(versionID, game)
		}
	}

	return &compatTools, nil
}
