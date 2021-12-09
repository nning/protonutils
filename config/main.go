package config

import (
	"os"
	"os/user"
	"path"

	"gopkg.in/yaml.v3"
)

type Config struct {
	User string `yaml:"user"`
}

func New() (*Config, error) {
	cfg := &Config{}
	err := cfg.Load()
	return cfg, err
}

func (cfg *Config) String() string {
	bytes, _ := yaml.Marshal(cfg)
	return string(bytes)
}

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
