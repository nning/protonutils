package config

import (
	"os"
	"path"

	"github.com/nning/protonutils/utils"
	"gopkg.in/yaml.v3"
)

// Config represents YAML config file structure
type Config struct {
	User      string `yaml:"user"`
	SteamRoot string `yaml:"steam_root"`
	dir       string
	file      string
}

// New instantiates and loads config from file
func New() (*Config, error) {
	dir := utils.GetConfigDir()
	file := path.Join(dir, "protonutils.yml")

	cfg := &Config{dir: dir, file: file}
	err := cfg.Load()

	return cfg, err
}

// String returns YAML serialization of cfg
func (cfg *Config) String() string {
	bytes, _ := yaml.Marshal(cfg)
	return string(bytes)
}

// Load loads cfg values from file
func (cfg *Config) Load() error {
	_, err := os.Stat(cfg.file)
	if err == nil {
		content, err := os.ReadFile(cfg.file)
		if err != nil {
			return err
		}

		yaml.Unmarshal(content, cfg)
	}

	if cfg.SteamRoot == "" {
		cfg.SteamRoot = "~/.local/share/Steam"
	}

	return nil
}

// Save saves cfg values into file
func (cfg *Config) Save() error {
	_, err := os.Stat(cfg.dir)
	if err != nil {
		err := os.MkdirAll(cfg.dir, 0700)
		if err != nil {
			return err
		}
	}

	content, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(cfg.file, content, 0600)
	if err != nil {
		return err
	}

	return nil
}
