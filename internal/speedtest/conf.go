package speedtest

import (
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
		log.Printf("[ConfigError]Unable to read host.ini: %v", err)
		os.Exit(1)
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
		log.Printf("[ReadDirError]Can't read config directory: %v", err)
		os.Exit(1)
	}
	var match int
	for _, file := range files {
		if file.Name() == "host.ini" {
			match++
		}
	}
	return match == 1
}

func Ping() {
	// TODO: use get version api to test host can be established
}
