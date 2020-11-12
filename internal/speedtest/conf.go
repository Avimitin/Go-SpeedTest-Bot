package speedtest

import (
	"encoding/json"
	"fmt"
	"go-speedtest-bot/internal/web"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"os"
)

func GetHost() *Host {
	if !configExistTest() {
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
	s := cfg.Section("host")
	p, _ := s.Key("port").Int()
	return &Host{
		IP:    s.Key("ip").String(),
		Port:  p,
		Token: s.Key("token").String(),
	}

}

func configExistTest() bool {
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
		if file.Name() == "host.ini" {
			match++
		}
	}
	return match == 1
}

func Ping(h *Host) bool {
	resp, err := web.Get(fmt.Sprintf("http://%s:%d", h.IP, h.Port))
	if err != nil {
		log.Printf("[PingError]Unable to connect to backend")
		return false
	}
	var v Version
	err = json.Unmarshal(resp, &v)
	if err != nil {
		log.Printf("[ParseError]Unable to unmarshall json data")
		return false
	}

	return v.Main != "" && v.WebAPI != ""
}
