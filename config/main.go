package config

import (
	"os"
	"os/user"
	"path"

	"gopkg.in/yaml.v3"
)

// Config represents YAML config file structure
type Config struct {
	User string `yaml:"user"`
}

// New instantiates and loads config from file
func New() (*Config, error) {
	cfg := &Config{}
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
	u, _ := user.Current()
	p := path.Join(u.HomeDir, ".config", "protonutils", "protonutils.yml")

	_, err := os.Stat(p)
	if err == nil {
		content, err := os.ReadFile(p)
		if err != nil {
			return err
		}

		yaml.Unmarshal(content, cfg)
	}

	return nil
}

// Save saves cfg values into file
func (cfg *Config) Save() error {
	u, _ := user.Current()
	dirPath := path.Join(u.HomeDir, ".config", "protonutils")
	filePath := path.Join(dirPath, "protonutils.yml")

	_, err := os.Stat(dirPath)
	if err != nil {
		err := os.MkdirAll(dirPath, 0700)
		if err != nil {
			return err
		}
	}

	content, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, content, 0600)
	if err != nil {
		return err
	}

	return nil

}
