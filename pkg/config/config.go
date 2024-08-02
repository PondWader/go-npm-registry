package config

import (
	"errors"
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

type Config struct {
	UserKeys []string `default:"[]"`
	DbPath   string   `default:"./sqlite.db"`
	Port     int      `default:"8080"`
}

func Load(path string) (Config, error) {
	config := Config{}
	defaults.Set(&config)

	if _, err := os.Stat(path); err == nil {
		configData, err := os.ReadFile(path)
		if err != nil {
			return config, err
		}
		err = yaml.Unmarshal(configData, &config)
		if err != nil {
			return config, err
		}
		return config, nil
	} else if errors.Is(err, os.ErrNotExist) {
		defaultConfig, _ := yaml.Marshal(config)
		os.WriteFile(path, defaultConfig, 0777)
		return config, nil
	} else {
		return config, err
	}
}
