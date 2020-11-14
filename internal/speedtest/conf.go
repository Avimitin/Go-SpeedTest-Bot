package speedtest

import (
	"encoding/json"
	"fmt"
	"go-speedtest-bot/internal/config"
	"go-speedtest-bot/internal/web"
	"log"
)

// GetHost will fetch local ini config file information
// If method can't find any config file it will try to fetch environmental variable "SPT_BOT_PATH"
func GetHost() *Host {
	cfg := config.GetHostFile()
	s := cfg.Section("host")
	p, _ := s.Key("port").Int()
	return &Host{
		IP:    s.Key("ip").String(),
		Port:  p,
		Token: s.Key("token").String(),
	}

}

// Ping will test if the given host is accessible or not
func Ping(h *Host) bool {
	resp, err := web.Get(fmt.Sprintf("http://%s:%d/getversion", h.IP, h.Port))
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
