package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func findConfigFilePath() string {
	var route string
	if route = os.Getenv("SPT_CFG_PATH"); route != "" {
		return path.Join(route, "config.json")
	}
	if route = os.Getenv("HOME"); route != "" {
		return path.Join(route, ".config", "spt_bot", "config.json")
	}
	return path.Join(".", "config", "config.json")
}

func GetConfig() (*Configuration, error) {
	path := findConfigFilePath()
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %q: %v", path, err)
	}
	var cfg *Configuration
	err = json.Unmarshal(configFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parse file: %v", err)
	}
	return cfg, nil
}

func GetToken() string {
	file, err := GetConfig()
	if err != nil {
		fatal(err)
	}
	return file.Global.Token
}

type HostConfig struct {
	address string
	key     string
}

func (hc *HostConfig) GetURL() string {
	return hc.address
}

func GetRunner(runnername string) *Runner {
	file, err := GetConfig()
	if err != nil {
		fatal(err)
	}
	for _, f := range file.Runner {
		if f.Name == runnername {
			return f
		}
	}
	return nil
}

func GetDefaultConfig(configname string) *Default {
	file, err := GetConfig()
	if err != nil {
		fatal(err)
	}
	for _, f := range file.DefaultConfig {
		if f.Name == configname {
			return f
		}
	}
	return nil
}

func fatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}
