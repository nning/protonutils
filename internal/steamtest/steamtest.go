// Package steamtest creates synthetic Steam root directories for tests
package steamtest

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// KV is a binary KeyValues node used to build appinfo.vdf and shortcuts.vdf
type KV struct {
	Key      string
	Value    any // string, int32, uint64 or nil for sub trees
	Children []KV
}

// Str returns string node
func Str(k, v string) KV { return KV{Key: k, Value: v} }

// Int returns int32 node
func Int(k string, v int32) KV { return KV{Key: k, Value: v} }

// Tree returns sub tree node
func Tree(k string, children ...KV) KV { return KV{Key: k, Children: children} }

type kvWriter struct {
	buf   bytes.Buffer
	table map[string]uint32
	keys  []string
}

func (w *kvWriter) key(k string) {
	if w.table == nil {
		w.buf.WriteString(k)
		w.buf.WriteByte(0)
		return
	}

	i, ok := w.table[k]
	if !ok {
		i = uint32(len(w.keys))
		w.table[k] = i
		w.keys = append(w.keys, k)
	}

	binary.Write(&w.buf, binary.LittleEndian, i)
}

func (w *kvWriter) node(n KV) {
	switch v := n.Value.(type) {
	case nil:
		w.buf.WriteByte(0x00)
		w.key(n.Key)
		for _, c := range n.Children {
			w.node(c)
		}
		w.buf.WriteByte(0x08)
	case string:
		w.buf.WriteByte(0x01)
		w.key(n.Key)
		w.buf.WriteString(v)
		w.buf.WriteByte(0)
	case int32:
		w.buf.WriteByte(0x02)
		w.key(n.Key)
		binary.Write(&w.buf, binary.LittleEndian, v)
	case uint64:
		w.buf.WriteByte(0x07)
		w.key(n.Key)
		binary.Write(&w.buf, binary.LittleEndian, v)
	default:
		panic(fmt.Sprintf("unsupported value type %T", v))
	}
}

// BinaryVdf encodes nodes as binary VDF (with trailing end marker)
func BinaryVdf(nodes ...KV) []byte {
	var w kvWriter
	for _, n := range nodes {
		w.node(n)
	}
	w.buf.WriteByte(0x08)
	return w.buf.Bytes()
}

// App is an entry of appinfo.vdf
type App struct {
	ID         uint32
	Name       string
	InstallDir string
	// DeckCategory is the Steam Deck compatibility category (0 = none)
	DeckCategory int32
	// Type is the app type (e.g. "Game", "Tool")
	Type string
	// Extra nodes appended to the "appinfo" node
	Extra []KV
}

func (a App) kv() KV {
	common := []KV{Str("name", a.Name)}
	if a.Type != "" {
		common = append(common, Str("type", a.Type))
	}
	if a.DeckCategory != 0 {
		common = append(common, Tree("steam_deck_compatibility", Int("category", a.DeckCategory)))
	}

	children := []KV{Int("appid", int32(a.ID)), Tree("common", common...)}
	if a.InstallDir != "" {
		children = append(children, Tree("config", Str("installdir", a.InstallDir)))
	}
	children = append(children, a.Extra...)

	return Tree("appinfo", children...)
}

// AppInfo encodes apps as appinfo.vdf of given version (27, 28 or 29)
func AppInfo(version int, apps ...App) []byte {
	var out bytes.Buffer
	le := binary.LittleEndian

	binary.Write(&out, le, uint32(0x07564400+version/10*16+version%10))
	binary.Write(&out, le, uint32(1))
	if version == 29 {
		binary.Write(&out, le, uint64(0)) // patched below
	}

	w := &kvWriter{}
	if version == 29 {
		w.table = make(map[string]uint32)
	}

	headerLen := 60
	if version == 27 {
		headerLen = 40
	}

	for _, a := range apps {
		w.buf.Reset()
		w.node(a.kv())
		w.buf.WriteByte(0x08)

		binary.Write(&out, le, a.ID)
		binary.Write(&out, le, uint32(headerLen+w.buf.Len()))
		out.Write(make([]byte, headerLen))
		out.Write(w.buf.Bytes())
	}

	binary.Write(&out, le, uint32(0))

	b := out.Bytes()
	if version == 29 {
		le.PutUint64(b[8:], uint64(len(b)))

		var table bytes.Buffer
		binary.Write(&table, le, uint32(len(w.keys)))
		for _, k := range w.keys {
			table.WriteString(k)
			table.WriteByte(0)
		}
		b = append(b, table.Bytes()...)
	}

	return b
}

// Shortcut is a non-Steam game shortcut
type Shortcut struct {
	ID   uint32
	Name string
}

// Shortcuts encodes shortcuts as shortcuts.vdf
func Shortcuts(shortcuts ...Shortcut) []byte {
	var children []KV
	for i, sc := range shortcuts {
		children = append(children, Tree(fmt.Sprint(i), Int("appid", int32(sc.ID)), Str("AppName", sc.Name)))
	}

	return BinaryVdf(Tree("shortcuts", children...))
}

// User is a Steam user
type User struct {
	ID64        uint64
	AccountName string
	MostRecent  bool
	// Apps are app IDs in localconfig.vdf
	Apps []string
	// Shortcuts of the user
	Shortcuts []Shortcut
}

// ID32 returns the account ID of the user
func (u User) ID32() string {
	return fmt.Sprint(u.ID64 & 0xffffffff)
}

// Tool is a custom compatibility tool in compatibilitytools.d
type Tool struct {
	Dir         string
	ID          string
	DisplayName string
}

// Root describes a synthetic Steam root
type Root struct {
	AppInfoVersion int
	Apps           []App
	Users          []User
	// Mapping maps app IDs to compatibility tool IDs ("0" is the default)
	Mapping map[string]string
	// Installed app IDs (in the library at LibraryPath)
	Installed   []string
	LibraryPath string
	Tools       []Tool
}

func quote(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func write(t testing.TB, p string, data []byte) {
	t.Helper()

	err := os.MkdirAll(filepath.Dir(p), 0700)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(p, data, 0600)
	if err != nil {
		t.Fatal(err)
	}
}

// Write creates the Steam root in dir
func (r Root) Write(t testing.TB, dir string) {
	t.Helper()

	version := r.AppInfoVersion
	if version == 0 {
		version = 29
	}
	write(t, filepath.Join(dir, "appcache", "appinfo.vdf"), AppInfo(version, r.Apps...))

	var b strings.Builder
	b.WriteString("\"users\"\n{\n")
	for _, u := range r.Users {
		mostRecent := "0"
		if u.MostRecent {
			mostRecent = "1"
		}
		fmt.Fprintf(&b, "\t\"%d\"\n\t{\n\t\t\"AccountName\"\t\t%v\n\t\t\"MostRecent\"\t\t\"%v\"\n\t}\n", u.ID64, quote(u.AccountName), mostRecent)
	}
	b.WriteString("}\n")
	write(t, filepath.Join(dir, "config", "loginusers.vdf"), []byte(b.String()))

	b.Reset()
	b.WriteString("\"InstallConfigStore\"\n{\n\t\"Software\"\n\t{\n\t\t\"Valve\"\n\t\t{\n\t\t\t\"Steam\"\n\t\t\t{\n\t\t\t\t\"CompatToolMapping\"\n\t\t\t\t{\n")
	for _, id := range sortedKeys(r.Mapping) {
		fmt.Fprintf(&b, "\t\t\t\t\t%v\n\t\t\t\t\t{\n\t\t\t\t\t\t\"name\"\t\t%v\n\t\t\t\t\t\t\"config\"\t\t\"\"\n\t\t\t\t\t\t\"priority\"\t\t\"250\"\n\t\t\t\t\t}\n", quote(id), quote(r.Mapping[id]))
	}
	b.WriteString("\t\t\t\t}\n\t\t\t}\n\t\t}\n\t}\n}\n")
	write(t, filepath.Join(dir, "config", "config.vdf"), []byte(b.String()))

	libraryPath := r.LibraryPath
	if libraryPath == "" {
		libraryPath = filepath.Join(dir, "library")
	}

	b.Reset()
	fmt.Fprintf(&b, "\"libraryfolders\"\n{\n\t\"0\"\n\t{\n\t\t\"path\"\t\t%v\n\t\t\"apps\"\n\t\t{\n", quote(libraryPath))
	for _, id := range r.Installed {
		fmt.Fprintf(&b, "\t\t\t%v\t\t\"1000\"\n", quote(id))
	}
	b.WriteString("\t\t}\n\t}\n}\n")
	write(t, filepath.Join(dir, "steamapps", "libraryfolders.vdf"), []byte(b.String()))

	err := os.MkdirAll(filepath.Join(dir, "userdata"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	for _, u := range r.Users {
		b.Reset()
		b.WriteString("\"UserLocalConfigStore\"\n{\n\t\"Software\"\n\t{\n\t\t\"Valve\"\n\t\t{\n\t\t\t\"Steam\"\n\t\t\t{\n\t\t\t\t\"apps\"\n\t\t\t\t{\n")
		for _, id := range u.Apps {
			fmt.Fprintf(&b, "\t\t\t\t\t%v\n\t\t\t\t\t{\n\t\t\t\t\t\t\"ViewedSteamPlay\"\t\t\"1\"\n\t\t\t\t\t}\n", quote(id))
		}
		b.WriteString("\t\t\t\t}\n\t\t\t}\n\t\t}\n\t}\n}\n")

		userConfig := filepath.Join(dir, "userdata", u.ID32(), "config")
		write(t, filepath.Join(userConfig, "localconfig.vdf"), []byte(b.String()))

		if len(u.Shortcuts) > 0 {
			write(t, filepath.Join(userConfig, "shortcuts.vdf"), Shortcuts(u.Shortcuts...))
		}
	}

	for _, tool := range r.Tools {
		b.Reset()
		fmt.Fprintf(&b, "\"compatibilitytools\"\n{\n\t\"compat_tools\"\n\t{\n\t\t%v\n\t\t{\n\t\t\t\"install_path\" \".\"\n\t\t\t\"display_name\" %v\n\t\t}\n\t}\n}\n", quote(tool.ID), quote(tool.DisplayName))

		toolDir := filepath.Join(dir, "compatibilitytools.d", tool.Dir)
		vdfPath := filepath.Join(toolDir, "compatibilitytool.vdf")

		if existing, err := os.ReadFile(vdfPath); err == nil {
			// Several tools in one directory
			s := strings.Replace(string(existing), "\t\"compat_tools\"\n\t{\n", "\t\"compat_tools\"\n\t{\n\t\t"+quote(tool.ID)+"\n\t\t{\n\t\t\t\"display_name\" "+quote(tool.DisplayName)+"\n\t\t}\n", 1)
			write(t, vdfPath, []byte(s))
		} else {
			write(t, vdfPath, []byte(b.String()))
		}
	}
}

// Default returns a typical Steam root description:
//
//   - user "alice" (most recent) and user "bob"
//   - default tool: GE-Proton9-1 (custom)
//   - Portal 2 (620, installed) -> proton_63
//   - Half-Life (70, installed) -> GE-Proton8-1 (custom)
//   - Portal (400, not installed) -> GE-Proton7-1 (custom)
//   - Cyberpunk 2077 (1091500, installed) -> default
//   - Proton 6.3 (1580130, installed, type Tool, must not be listed as game)
//   - shortcut "My Game" (3228583970) -> proton_63
//   - unused custom tool GE-Proton6-1
func Default() Root {
	return Root{
		Apps: []App{
			{ID: 70, Name: "Half-Life", InstallDir: "Half-Life", DeckCategory: 3},
			{ID: 400, Name: "Portal", InstallDir: "Portal"},
			{ID: 620, Name: "Portal 2", InstallDir: "Portal 2", DeckCategory: 3},
			{ID: 1091500, Name: "Cyberpunk 2077", InstallDir: "Cyberpunk 2077", DeckCategory: 2, Type: "Game"},
			{ID: 1580130, Name: "Proton 6.3", InstallDir: "Proton 6.3", Type: "Tool"},
			{ID: 891390, Name: "SteamPlay 2.0 Manifests", Extra: []KV{
				Tree("extended", Tree("compat_tools",
					Tree("proton_63", Str("display_name", "Proton 6.3-8")),
					Tree("proton_experimental", Str("display_name", "Proton - Experimental")),
				)),
			}},
		},
		Users: []User{
			{
				ID64: 76561198050517827, AccountName: "alice", MostRecent: true,
				Apps:      []string{"70", "400", "620", "1091500", "1580130", "3228583970"},
				Shortcuts: []Shortcut{{ID: 3228583970, Name: "My Game"}},
			},
			{ID64: 76561198050517828, AccountName: "bob", Apps: []string{"620"}},
		},
		Mapping: map[string]string{
			"0":          "GE-Proton9-1",
			"620":        "proton_63",
			"70":         "GE-Proton8-1",
			"400":        "GE-Proton7-1",
			"3228583970": "proton_63",
		},
		Installed: []string{"70", "620", "1091500", "1580130"},
		Tools: []Tool{
			{Dir: "GE-Proton6-1", ID: "GE-Proton6-1", DisplayName: "GE-Proton6-1"},
			{Dir: "GE-Proton7-1", ID: "GE-Proton7-1", DisplayName: "GE-Proton7-1"},
			{Dir: "GE-Proton8-1", ID: "GE-Proton8-1", DisplayName: "GE-Proton8-1"},
			{Dir: "GE-Proton9-1", ID: "GE-Proton9-1", DisplayName: "GE-Proton9-1"},
		},
	}
}

// Setup writes the Default Steam root to a temporary directory, isolates
// config and cache directories via XDG environment variables and returns the
// Steam root path
func Setup(t testing.TB, r Root) string {
	t.Helper()

	tmp := t.TempDir()
	t.Setenv("HOME", filepath.Join(tmp, "home"))
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, "config"))
	t.Setenv("XDG_CACHE_HOME", filepath.Join(tmp, "cache"))

	dir := filepath.Join(tmp, "steam")
	r.Write(t, dir)

	return dir
}
