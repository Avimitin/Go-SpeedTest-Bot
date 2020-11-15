package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"os"
)

// ConfigExistTest test the integrity of the configuration file
func CFGExistTest() bool {
	files, err := ioutil.ReadDir("./config")
	if err != nil {
		if env := os.Getenv("SPT_BOT_PATH"); env != "" {
			files, err = ioutil.ReadDir(fmt.Sprintf("%s/config", os.Getenv("SPT_BOT_PATH")))
		}
		if err != nil {
			log.Printf("[ReadDirError]Can't read config directory: %v", err)
			os.Exit(1)
		}
	}
	var match int
	for _, file := range files {
		if fname := file.Name(); fname == "host.ini" || fname == "bot.ini" {
			match++
		}
	}
	return match == 2
}

// GetHostFile return parsing host configuration file
func GetHostFile() *ini.File {
	if !CFGExistTest() {
		log.Printf("[ConfigError]No enough config")
		os.Exit(1)
	}
	cfg, err := ini.Load("./config/host.ini")
	if err != nil {
		cfg, err = ini.Load(fmt.Sprintf("%s/config/host.ini", os.Getenv("SPT_BOT_PATH")))
		if err != nil {
			log.Printf("[ConfigError]Unable to read host.ini: %v", err)
			os.Exit(1)
		}
	}
	return cfg
}

// GetBotFile return parsing bot configuration file
func GetBotFile() *ini.File {
	if !CFGExistTest() {
		log.Printf("[ConfigError]No enough config")
		os.Exit(1)
	}
	cfg, err := ini.Load("./config/bot.ini")
	if err != nil {
		cfg, err = ini.Load(fmt.Sprintf("%s/config/bot.ini", os.Getenv("SPT_BOT_PATH")))
		if err != nil {
			log.Printf("[ConfigError]Unable to read subs.ini: %v", err)
			os.Exit(1)
		}
	}
	return cfg
}

func GetSubsFile() *ini.File {
	if !CFGExistTest() {
		log.Printf("[ConfigError]No enough config")
		os.Exit(1)
	}
	cfg, err := ini.Load("./config/subs.ini")
	if err != nil {
		cfg, err = ini.Load(fmt.Sprintf("%s/config/subs.ini", os.Getenv("SPT_BOT_PATH")))
		if err != nil {
			log.Printf("[ConfigError]Unable to read subs.ini: %v", err)
			os.Exit(1)
		}
	}
	return cfg
}
