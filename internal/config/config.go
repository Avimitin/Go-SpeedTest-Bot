package config

import (
	"encoding/json"
	"go-speedtest-bot/internal/runner"
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
	path := findConfigFilePath()
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("read %q: %v", path, err)
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

// ListAllRunners return all the usable runner name
func ListAllRunners() string {
	runners := GetAllRunner()
	var text string = "available runner:\n"
	for _, r := range runners {
		text += r.Name + "\n"
	}
	return text
}

// GetRunner return a specific runner
func GetRunner(runnername string) *runner.Runner {
	runners := GetAllRunner()
	for _, f := range runners {
		if f.Name == runnername {
			return f
		}
	}
	return nil
}

// GetDefaultConfig return default speedtest config
func GetDefaultConfig(configname string) *Default {
	defaultConfig := GetAllDefaultConfig()
	for _, f := range defaultConfig {
		if f.Name == configname {
			return f
		}
	}
	return nil
}

// GetAllDefaultConfig return all the default config
func GetAllDefaultConfig() []*Default {
	return userSetting.DefaultConfig
}
