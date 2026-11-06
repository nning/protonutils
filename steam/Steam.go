package steam

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nning/protonutils/cache"
	"github.com/nning/protonutils/config"
	"github.com/nning/protonutils/utils"
)

// LegacyDefaultRoot is the Steam root that was used as default by older
// versions of protonutils (and written to the config file)
const LegacyDefaultRoot = "~/.local/share/Steam"

// RootCandidates are the Steam root locations that are tried if no Steam root
// is configured
var RootCandidates = []string{
	"~/.local/share/Steam",
	"~/.steam/root",
	"~/.steam/steam",
	"~/.steam/debian-installation",
	"~/.var/app/com.valvesoftware.Steam/data/Steam",
	"~/.var/app/com.valvesoftware.Steam/.local/share/Steam",
}

// versionNameCacheMaxAge is the time after which names of compatibility tools
// are looked up again
const versionNameCacheMaxAge = 6 * time.Hour

// Steam struct wraps caches and exposes functions for Steam data retrieval
type Steam struct {
	AppidCache       *cache.Cache
	VersionNameCache *cache.Cache

	AppInfo           *AppInfo
	CompatToolMapping *CompatToolMappingVdf
	LibraryConfig     *LibraryConfigVdf
	LocalConfig       *LocalConfigVdf
	LoginUsers        *LoginUsersVdf
	Shortcuts         *Shortcuts
	CustomTools       CustomTools

	CompatTools CompatTools

	UID                   string
	Root                  string
	EnableViewedSteamPlay bool

	// Stderr receives warnings (defaults to os.Stderr)
	Stderr io.Writer

	apiUnavailable bool
}

func isDir(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

func isSteamRoot(p string) bool {
	return isDir(filepath.Join(p, "config")) && isDir(filepath.Join(p, "steamapps"))
}

// FindRoot returns the Steam root directory. If configured is empty (or the
// legacy default), the RootCandidates are tried.
func FindRoot(configured string) (string, error) {
	if configured != "" && configured != LegacyDefaultRoot {
		root := utils.ExpandHome(configured)
		if !isDir(root) {
			return "", fmt.Errorf("Steam root not a directory: %v", root)
		}

		return root, nil
	}

	for _, candidate := range RootCandidates {
		root := utils.ExpandHome(candidate)
		if isSteamRoot(root) {
			return root, nil
		}
	}

	return "", errors.New("Steam root not found (tried " + strings.Join(RootCandidates, ", ") + "), " +
		"set it with: protonutils config steam_root <path>")
}

// New instantiates Steam struct
func New(user string, cfg *config.Config, ignoreCache bool) (*Steam, error) {
	root, err := FindRoot(cfg.SteamRoot)
	if err != nil {
		return nil, err
	}

	var appidMaxAge, versionNameMaxAge time.Duration = -1, versionNameCacheMaxAge
	if ignoreCache {
		appidMaxAge, versionNameMaxAge = 0, 0
	}

	appidCache, err := cache.New("appids", appidMaxAge)
	if err != nil {
		return nil, err
	}

	versionNameCache, err := cache.New("proton-names", versionNameMaxAge)
	if err != nil {
		return nil, err
	}

	s := &Steam{
		AppidCache:            appidCache,
		VersionNameCache:      versionNameCache,
		Root:                  root,
		EnableViewedSteamPlay: cfg.EnableViewedSteamPlay,
		CompatTools:           make(CompatTools),
		Stderr:                os.Stderr,
	}

	err = s.initLoginUsers()
	if err != nil {
		return nil, err
	}

	s.UID, err = s.resolveUID(user, s.Stderr)
	if err != nil {
		return nil, err
	}

	slog.Debug("steam.New", "user", user, "uid", s.UID, "root", root, "cache", !ignoreCache, "enableViewedSteamPlay", s.EnableViewedSteamPlay)

	err = s.initConfigs()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Steam) initConfigs() error {
	inits := []func() error{
		s.initAppInfo,
		s.initCustomTools,
		s.initCompatToolMapping,
		s.initLibraryConfig,
		s.initLocalConfig,
		s.initShortcuts,
	}

	for _, init := range inits {
		if err := init(); err != nil {
			return err
		}
	}

	return nil
}

// SaveCache writes caches to disk
func (s *Steam) SaveCache() error {
	err := s.AppidCache.Write()
	if err != nil {
		return err
	}

	return s.VersionNameCache.Write()
}

// GetCompatibilityToolsDir returns compatibilitytools.d location
func (s *Steam) GetCompatibilityToolsDir() string {
	return filepath.Join(s.Root, "compatibilitytools.d")
}

// IsFlatpak returns whether protonutils runs in a Flatpak sandbox
func IsFlatpak() bool {
	_, err := os.Stat("/.flatpak-info")
	return err == nil
}

// IsRunning returns whether the Steam client is running (a process with an
// executable named "steam"). Always returns false if run in Flatpak sandbox.
func (s *Steam) IsRunning() (bool, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		link, err := os.Readlink(filepath.Join("/proc", entry.Name(), "exe"))
		if err != nil {
			continue
		}

		if filepath.Base(strings.TrimSuffix(link, " (deleted)")) == "steam" {
			return true, nil
		}
	}

	return false, nil
}
