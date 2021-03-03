package config

// Configuration contains bot, runner configuration
// at run time; Also store user define default
// subscription config
type Configuration struct {
	Global        Global     `json:"global"`
	Runner        []*Runner  `json:"runner"`
	DefaultConfig []*Default `json:"default_config"`
}

type Global struct {
	Token string `json:"token"`
	Admin int    `json:"admin"`
}

type Runner struct {
	Name string `json:"name"`
	Host struct {
		Address string `json:"address"`
		Key     string `json:"key"`
	} `json:"host"`
	Admins []int `json:"admins"`
}

type Default struct {
	Name   string `json:"name"`
	Link   string `json:"link"`
	Chat   int    `json:"chat"`
	Admins []int  `json:"admins"`
}
