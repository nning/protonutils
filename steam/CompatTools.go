package steam

import (
	"sort"
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

// IsInstalled returns if a tool is installed
func (t CompatTool) IsInstalled(s *Steam) bool {
	// TODO Implement for non custom tools
	if !t.IsCustom {
		return false
	}

	tool := s.CustomTools.Get(t.ID)
	return tool != nil && isDir(tool.Dir)
}

// IsValid checks whether a version ID (v) exists in the CompatTools config
func (c CompatTools) IsValid(v string) bool {
	_, ok := c[v]
	return ok
}

// Add adds an entry to the CompatTools config (by tool id and (display) name)
func (c CompatTools) Add(id, name string) {
	if c[id] != nil {
		if name != "" {
			c[id].Name = name
		}

		return
	}

	c[id] = &CompatTool{
		ID:        id,
		Name:      name,
		IsDefault: id == "",
		IsCustom:  id != "" && id == name,
		Games:     make(Games),
	}
}

// AddGame adds a Game entry to an existing CompatTool entry
func (c CompatTools) AddGame(toolID string, game *Game) bool {
	if c[toolID] == nil {
		return false
	}

	c[toolID].Games[game.ID] = game
	return true
}

// ReadCompatTools reads compatibility tool mappings from config, and sets
// s.CompatTools accordingly.
func (s *Steam) ReadCompatTools() error {
	tools, err := s.CompatToolMapping.ReadCompatTools()
	if err != nil {
		return err
	}

	games, err := s.LocalConfig.GetGames(s.EnableViewedSteamPlay)
	if err != nil {
		return err
	}

	def := tools.GetDefault()
	if def == nil {
		// No explicit default configured, Steam uses its built-in default
		tools.Add("", "Default")
		def = tools[""]
	}

	mapped := tools.GameIDs()
	for _, game := range games {
		if mapped[game.ID] {
			continue
		}

		tools.AddGame(def.ID, game)
	}

	for id, tool := range s.CustomTools {
		name := tool.DisplayName
		if name == "" {
			name = id
		}

		tools.Add(id, name)
		tools[id].IsCustom = true
	}

	s.CompatTools = tools
	return nil
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

// Sort returns slice of CompatTools IDs sorted by name (and ID)
func (c CompatTools) Sort() []string {
	ids := make([]string, 0, len(c))
	for id := range c {
		ids = append(ids, id)
	}

	sort.Slice(ids, func(i, j int) bool {
		a, b := c[ids[i]], c[ids[j]]
		if a.Name != b.Name {
			return a.Name < b.Name
		}

		return a.ID < b.ID
	})

	return ids
}

// GameIDs returns set of app IDs of all games in c
func (c CompatTools) GameIDs() map[string]bool {
	ids := make(map[string]bool)
	for _, tool := range c {
		for id := range tool.Games {
			ids[id] = true
		}
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

// UsedToolIDs returns set of IDs of compatibility tools that are used by at
// least one installed game or configured as default
func (s *Steam) UsedToolIDs() map[string]bool {
	used := make(map[string]bool)

	for id, tool := range s.CompatTools {
		if id != "" && tool.Games.CountInstalled() > 0 {
			used[id] = true
		}
	}

	if def := s.CompatToolMapping.GetDefaultID(); def != "" {
		used[def] = true
	}

	return used
}

// UnusedCustomToolDirs returns sorted list of directories in
// compatibilitytools.d that only contain tools that are neither used by an
// installed game nor configured as default. s.ReadCompatTools has to be called
// before.
func (s *Steam) UnusedCustomToolDirs() []string {
	used := s.UsedToolIDs()

	usedDirs := make(map[string]bool)
	for id, tool := range s.CustomTools {
		if used[id] {
			usedDirs[tool.Dir] = true
		}
	}

	var dirs []string
	for _, dir := range s.CustomTools.Dirs() {
		if !usedDirs[dir] {
			dirs = append(dirs, dir)
		}
	}

	return dirs
}

// NewCompatTools returns new CompatTools struct (optionally initialized with
// data, which is a map from compatibility tool version IDs to slices of game
// IDs)
func (s *Steam) NewCompatTools(data ...map[string][]string) (*CompatTools, error) {
	compatTools := make(CompatTools)

	if len(data) > 0 {
		for versionID, games := range data[0] {
			for _, gameID := range games {
				game, _, err := s.GetGameData(gameID)
				if err != nil {
					return nil, err
				}

				versionName, _ := s.GetCompatToolName(versionID)

				compatTools.Add(versionID, versionName)
				compatTools.AddGame(versionID, game)
			}
		}
	}

	return &compatTools, nil
}
