package config

import (
	"encoding/json"
	"go-speedtest-bot/module/runner"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var (
	userSetting *Configuration
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

// LoadConfig is a reusable initialize config function
func LoadConfig() {
	configPath := findConfigFilePath()
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("read %q: %v", configPath, err)
	}
	err = json.Unmarshal(configFile, &userSetting)
	if err != nil {
		log.Fatalf("parse file: %v", err)
	}
}

// GetToken return bot token
func GetToken() string {
	return userSetting.Global.Token
}

// GetAllRunner return all predefine runner
func GetAllRunner() []*runner.Runner {
	return userSetting.Runner
}

// GetRunner return a specific runner
func GetRunner(name string) *runner.Runner {
	runners := GetAllRunner()
	for _, f := range runners {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// GetDefaultConfig return default speedtest config
func GetDefaultConfig(name string) *Default {
	defaultConfig := GetAllDefaultConfig()
	for _, f := range defaultConfig {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// GetAllDefaultConfig return all the default config
func GetAllDefaultConfig() []*Default {
	return userSetting.DefaultConfig
}

// GetPasteBinSetting return setting of pastebin.com
func GetPasteBinKey() string {
	return userSetting.PB.Key
}

// PasteBinEnabled return if pastebin feature is enabled
func PasteBinEnabled() bool {
	if userSetting.PB == nil {
		return false
	}
	return userSetting.PB.Enable
}
