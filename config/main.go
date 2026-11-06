package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/nning/protonutils/utils"
	"gopkg.in/yaml.v3"
)

// Config represents YAML config file structure
type Config struct {
	User                  string `yaml:"user"`
	SteamRoot             string `yaml:"steam_root"`
	EnableViewedSteamPlay bool   `yaml:"enable_viewed_steam_play"`

	dir  string
	file string
}

// UnknownKeyError is returned if a config key does not exist
type UnknownKeyError struct {
	Key string
}

func (e *UnknownKeyError) Error() string {
	return fmt.Sprintf("unknown config key %q (valid keys: %v)", e.Key, Keys())
}

// New instantiates and loads config from file
func New() (*Config, error) {
	dir := utils.GetConfigDir()
	file := filepath.Join(dir, "protonutils.yml")

	cfg := &Config{dir: dir, file: file}

	err := migrateLegacyFile(file)
	if err != nil {
		return cfg, err
	}

	err = cfg.Load()

	return cfg, err
}

// migrateLegacyFile moves a config file from the location used by older
// versions ($XDG_CONFIG_HOME/protonutils.yml) to file
func migrateLegacyFile(file string) error {
	legacy := utils.GetLegacyConfigFile()
	if legacy == "" {
		return nil
	}

	if _, err := os.Stat(file); err == nil {
		return nil
	}

	if _, err := os.Stat(legacy); err != nil {
		return nil
	}

	err := os.MkdirAll(filepath.Dir(file), 0700)
	if err != nil {
		return err
	}

	return os.Rename(legacy, file)
}

// Path returns the location of the config file
func (cfg *Config) Path() string {
	return cfg.file
}

// String returns YAML serialization of cfg
func (cfg *Config) String() string {
	bytes, _ := yaml.Marshal(cfg)
	return string(bytes)
}

// Load loads cfg values from file
func (cfg *Config) Load() error {
	content, err := os.ReadFile(cfg.file)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", cfg.file, err)
	}

	return nil
}

// Keys returns sorted list of valid config keys
func Keys() []string {
	keys := []string{"user", "steam_root", "enable_viewed_steam_play"}
	sort.Strings(keys)
	return keys
}

// Get returns the value of config key as string
func (cfg *Config) Get(key string) (string, error) {
	switch key {
	case "user":
		return cfg.User, nil
	case "steam_root":
		return cfg.SteamRoot, nil
	case "enable_viewed_steam_play":
		return strconv.FormatBool(cfg.EnableViewedSteamPlay), nil
	}

	return "", &UnknownKeyError{key}
}

// Set validates and sets the value of config key (without saving)
func (cfg *Config) Set(key, value string) error {
	switch key {
	case "user":
		cfg.User = value
	case "steam_root":
		cfg.SteamRoot = value
	case "enable_viewed_steam_play":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value %q for %v", value, key)
		}
		cfg.EnableViewedSteamPlay = b
	default:
		return &UnknownKeyError{key}
	}

	return nil
}

// Save saves cfg values into file
func (cfg *Config) Save() error {
	err := os.MkdirAll(cfg.dir, 0700)
	if err != nil {
		return err
	}

	content, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(cfg.file, content, 0600)
}
