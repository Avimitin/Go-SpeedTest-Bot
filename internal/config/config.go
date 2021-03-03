package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

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
