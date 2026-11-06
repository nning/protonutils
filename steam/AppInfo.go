package steam

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BenLubar/vdf"
	"github.com/dustin/go-humanize"
)

// Magic numbers (versions) of appinfo.vdf
const (
	AppInfoMagicV27 = 0x07564427
	AppInfoMagicV28 = 0x07564428
	AppInfoMagicV29 = 0x07564429
)

// steamPlayManifestsAppID is the app ID of "SteamPlay 2.0 Manifests", which
// contains the display names of official compatibility tools
const steamPlayManifestsAppID = "891390"

// AppInfo provides indexed access to the entries of appinfo.vdf
type AppInfo struct {
	Path    string
	Version uint32

	data    []byte
	keys    []string
	entries map[uint32][]byte
	parsed  map[uint32]*vdf.Node
}

// ParseAppInfo indexes the content of appinfo.vdf (supported versions: 27, 28,
// 29). Entries are only parsed on access.
func ParseAppInfo(data []byte, path string) (*AppInfo, error) {
	a := &AppInfo{
		Path:    path,
		data:    data,
		entries: make(map[uint32][]byte),
		parsed:  make(map[uint32]*vdf.Node),
	}

	if len(data) == 0 {
		return a, nil
	}

	if len(data) < 8 {
		return nil, errors.New("appinfo.vdf: file too short")
	}

	a.Version = binary.LittleEndian.Uint32(data)

	var headerLen, pos int
	end := len(data)

	switch a.Version {
	case AppInfoMagicV27:
		headerLen, pos = 40, 8
	case AppInfoMagicV28:
		headerLen, pos = 60, 8
	case AppInfoMagicV29:
		headerLen, pos = 60, 16
		if len(data) < 16 {
			return nil, errors.New("appinfo.vdf: file too short")
		}

		tableOffset := binary.LittleEndian.Uint64(data[8:])
		if tableOffset > uint64(len(data)) {
			return nil, errors.New("appinfo.vdf: invalid string table offset")
		}

		end = int(tableOffset)

		var err error
		a.keys, err = parseAppInfoStringTable(data[tableOffset:])
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("appinfo.vdf: unsupported version 0x%08x", a.Version)
	}

	for pos+4 <= end {
		appID := binary.LittleEndian.Uint32(data[pos:])
		if appID == 0 {
			break
		}

		if pos+8 > end {
			return nil, fmt.Errorf("appinfo.vdf: truncated entry for app %d", appID)
		}

		size := int(binary.LittleEndian.Uint32(data[pos+4:]))
		start := pos + 8
		if size < headerLen || start+size > end {
			return nil, fmt.Errorf("appinfo.vdf: invalid size of entry for app %d", appID)
		}

		a.entries[appID] = data[start+headerLen : start+size]
		pos = start + size
	}

	return a, nil
}

func parseAppInfoStringTable(data []byte) ([]string, error) {
	r := &binaryVdfReader{data: data}

	count, err := r.readUint32()
	if err != nil {
		return nil, fmt.Errorf("appinfo.vdf: string table: %w", err)
	}

	if int(count) > len(data) {
		return nil, errors.New("appinfo.vdf: string table: invalid count")
	}

	keys := make([]string, count)
	for i := range keys {
		keys[i], err = r.readCString()
		if err != nil {
			return nil, fmt.Errorf("appinfo.vdf: string table: %w", err)
		}
	}

	return keys, nil
}

// Len returns the number of entries in appinfo.vdf
func (a *AppInfo) Len() int {
	return len(a.entries)
}

// Get returns the parsed "appinfo" node for app ID id or nil if there is no
// entry for id
func (a *AppInfo) Get(id string) (*vdf.Node, error) {
	if id == "" || id == "0" {
		return nil, nil
	}

	appID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		// Not a valid app ID, so there is no entry
		return nil, nil
	}

	if n, ok := a.parsed[uint32(appID)]; ok {
		return n, nil
	}

	entry, ok := a.entries[uint32(appID)]
	if !ok {
		return nil, nil
	}

	n, err := parseBinaryVdf(entry, a.keys)
	if err != nil {
		return nil, fmt.Errorf("appinfo.vdf: app %v: %w", id, err)
	}

	a.parsed[uint32(appID)] = n
	return n, nil
}

// GetType returns the app type from appinfo.vdf (e.g. "Game", "Tool",
// "Config"; empty string if unknown)
func (a *AppInfo) GetType(id string) string {
	n, err := a.Get(id)
	if n == nil || err != nil {
		return ""
	}

	return n.FirstByName("common").FirstByName("type").String()
}

// IsInfrastructure returns whether app id is part of the Steam infrastructure
// instead of something that can run with a compatibility tool: configuration
// apps (like "Steam Client") and tools without Windows version (like Proton or
// the Steam Linux Runtime). Windows tools (like authoring tools) are not
// considered infrastructure.
func (a *AppInfo) IsInfrastructure(id string) bool {
	n, err := a.Get(id)
	if n == nil || err != nil {
		return false
	}

	common := n.FirstByName("common")
	t := common.FirstByName("type").String()

	if strings.EqualFold(t, "Config") {
		return true
	}

	if !strings.EqualFold(t, "Tool") {
		return false
	}

	for _, os := range strings.Split(common.FirstByName("oslist").String(), ",") {
		if strings.EqualFold(strings.TrimSpace(os), "windows") {
			return false
		}
	}

	return true
}

// GetCompatToolName returns human-readable name of compatibility tool,
// for example: "proton_63" -> "Proton 6.3-8". If no name is found, id is
// returned.
func (s *Steam) GetCompatToolName(id string) (string, error) {
	if id == "" {
		return "", nil
	}

	if tool := s.CustomTools.Get(id); tool != nil && tool.DisplayName != "" {
		return tool.DisplayName, nil
	}

	str, _ := s.VersionNameCache.Get(id)
	if str != "" {
		return str, nil
	}

	n, err := s.AppInfo.Get(steamPlayManifestsAppID)
	if err != nil {
		s.VersionNameCache.Add(id, id, false)
		return id, err
	}

	name := n.FirstByName("extended").FirstByName("compat_tools").FirstByName(id).FirstByName("display_name").String()
	if name == "" {
		s.VersionNameCache.Add(id, id, false)
		return id, nil
	}

	s.VersionNameCache.Add(id, name, true)
	return name, nil
}

// GetInstalldir returns installation directory of game by app ID (empty string
// if app ID is unknown)
func (s *Steam) GetInstalldir(id string) (string, error) {
	n, err := s.AppInfo.Get(id)
	if n == nil || err != nil {
		return "", err
	}

	dir := n.FirstByName("config").FirstByName("installdir").String()
	if dir == "" {
		return "", errors.New("installdir not found")
	}

	return dir, nil
}

func (s *Steam) initAppInfo() error {
	p := filepath.Join(s.Root, "appcache", "appinfo.vdf")
	in, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	s.AppInfo, err = ParseAppInfo(in, p)
	if err != nil {
		// Do not fail completely (e.g. on a future format change), names
		// can still be retrieved from the store API
		fmt.Fprintf(s.Stderr, "Warning: %v, game and tool names may be missing\n\n", err)
		s.AppInfo, _ = ParseAppInfo(nil, p)
	}

	slog.Debug("Steam.initAppInfo", "path", p, "size", humanize.Bytes(uint64(len(in))), "version", fmt.Sprintf("0x%08x", s.AppInfo.Version), "entries", s.AppInfo.Len())

	return nil
}
