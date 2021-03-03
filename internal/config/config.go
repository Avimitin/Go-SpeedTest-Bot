package config

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/ini.v1"
)

func findConfigFilePath() string {
	var route string
	if route = os.Getenv("SPT_CFG_PATH"); route != "" {
		return path.Join(route, "config.ini")
	}
	if route = os.Getenv("HOME"); route != "" {
		return path.Join(route, ".config", "spt_bot", "config.ini")
	}
	return path.Join(".", "config", "config.ini")
}

func GetConfig() (*ini.File, error) {
	path := findConfigFilePath()
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %v", path, err)
	}
	return cfg, nil
}
