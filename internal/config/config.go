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

func GetToken() string {
	file, err := GetConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return file.Section("Bot").Key("token").String()
}

type HostConfig struct {
	address string
	key     string
}

func (hc *HostConfig) GetURL() string {
	return hc.address
}

func GetHost() *HostConfig {
	file, err := GetConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &HostConfig{
		address: file.Section("Host").Key("address").String(),
		key:     file.Section("Host").Key("key").String(),
	}
}
