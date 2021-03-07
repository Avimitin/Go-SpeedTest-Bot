package config

import (
	"encoding/json"
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

func GetToken() string {
	return userSetting.Global.Token
}

func GetAllRunner() []*Runner {
	return userSetting.Runner
}

func ListAllRunners() string {
	runners := GetAllRunner()
	var text string = "available runner:\n"
	for _, r := range runners {
		text += r.Name + "\n"
	}
	return text
}

func GetRunner(runnername string) *Runner {
	runners := GetAllRunner()
	for _, f := range runners {
		if f.Name == runnername {
			return f
		}
	}
	return nil
}

func GetDefaultConfig(configname string) *Default {
	defaultConfig := GetAllDefaultConfig()
	for _, f := range defaultConfig {
		if f.Name == configname {
			return f
		}
	}
	return nil
}

func GetAllDefaultConfig() []*Default {
	return userSetting.DefaultConfig
}
