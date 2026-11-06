package steam

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CustomTool represents a compatibility tool installed in compatibilitytools.d
type CustomTool struct {
	// ID is the internal name used in compatibility tool mappings
	ID string
	// DisplayName is the human-readable name (may be empty)
	DisplayName string
	// Dir is the directory (in compatibilitytools.d) containing the tool
	Dir string
}

// CustomTools maps internal names of custom compatibility tools to CustomTool
type CustomTools map[string]*CustomTool

// Get returns custom tool by internal name (nil if not found)
func (c CustomTools) Get(id string) *CustomTool {
	if c == nil {
		return nil
	}

	return c[id]
}

// Dirs returns sorted list of directories containing custom tools
func (c CustomTools) Dirs() []string {
	set := make(map[string]bool)
	for _, tool := range c {
		set[tool.Dir] = true
	}

	dirs := make([]string, 0, len(set))
	for dir := range set {
		dirs = append(dirs, dir)
	}

	sort.Strings(dirs)
	return dirs
}

// ReadCustomTools reads all compatibility tools from dir (usually
// compatibilitytools.d). Internal and display names are read from
// compatibilitytool.vdf. If that file is missing or cannot be parsed, the
// directory name is used as internal name.
func ReadCustomTools(dir string) (CustomTools, error) {
	tools := make(CustomTools)

	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return tools, nil
	}
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		p := filepath.Join(dir, name)
		if !isDir(p) {
			continue
		}

		found := false
		n, err := ParseTextConfig(filepath.Join(p, "compatibilitytool.vdf"))
		if err == nil {
			x := n.FirstByName("compat_tools")
			if n.Name() == "compat_tools" {
				x = n
			}

			for t := x.FirstSubTree(); t != nil; t = t.NextSubTree() {
				tools[t.Name()] = &CustomTool{
					ID:          t.Name(),
					DisplayName: t.FirstByName("display_name").String(),
					Dir:         p,
				}
				found = true
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			slog.Debug("ReadCustomTools: ignoring invalid compatibilitytool.vdf", "dir", p, "err", err)
		}

		if !found {
			tools[name] = &CustomTool{ID: name, Dir: p}
		}
	}

	return tools, nil
}

func (s *Steam) initCustomTools() error {
	var err error
	s.CustomTools, err = ReadCustomTools(s.GetCompatibilityToolsDir())
	return err
}

// ReloadCustomTools reads compatibilitytools.d again (e.g. after installing a
// tool)
func (s *Steam) ReloadCustomTools() error {
	return s.initCustomTools()
}
