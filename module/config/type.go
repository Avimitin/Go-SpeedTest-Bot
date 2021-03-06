package config

import "go-speedtest-bot/module/runner"

// Configuration contains bot, runner configuration
// at run time; Also store user define default
// subscription config
type Configuration struct {
	Global        Global           `json:"global"`
	Runner        []*runner.Runner `json:"runner"`
	DefaultConfig []*Default       `json:"default_config"`
	PB            *PasteBin        `json:"pastebin"`
}

type Global struct {
	Token string `json:"token"`
	Admin int    `json:"admin"`
}

type Default struct {
	Name          string `json:"name"`
	Link          string `json:"link"`
	Chat          int64  `json:"chat"`
	Admins        []int  `json:"admins"`
	DefaultRunner string `json:"default_runner"`
	Interval      int    `json:"interval"`
}

func (def *Default) HasAccess(id int) bool {
	for _, admin := range def.Admins {
		if admin == id {
			return true
		}
	}
	return false
}

type PasteBin struct {
	Enable bool   `json:"enable"`
	Key    string `json:"key"`
}
