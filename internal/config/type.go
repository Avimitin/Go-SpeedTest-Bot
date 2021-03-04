package config

// Configuration contains bot, runner configuration
// at run time; Also store user define default
// subscription config
type Configuration struct {
	Global        Global     `json:"global"`
	Runner        []*Runner  `json:"runner"`
	DefaultConfig []*Default `json:"default_config"`
}

type Runner struct {
	Name   string `json:"name"`
	Host   Host   `json:"host"`
	Admins []int  `json:"admins"`
}

type Global struct {
	Token string `json:"token"`
	Admin int    `json:"admin"`
}

type Host struct {
	Address string `json:"address"`
	Key     string `json:"key"`
}

func (h *Host) GetURL() string {
	return h.Address
}

type Default struct {
	Name          string `json:"name"`
	Link          string `json:"link"`
	Chat          int64  `json:"chat"`
	Admins        []int  `json:"admins"`
	DefaultRunner string `json:"default_runner"`
}
