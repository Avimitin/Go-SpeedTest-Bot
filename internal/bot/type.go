package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type B = tgbotapi.BotAPI
type M = tgbotapi.Message

type CMDFunc func(*M)
type CMD = map[string]CMDFunc

var Commands = CMD{
	"start":     cmdStart,
	"ping":      cmdPing,
	"status":    cmdStatus,
	"read_sub":  cmdReadSub,
	"result":    cmdResult,
	"run_url":   cmdStartTestWithURL,
	"list_subs": cmdListSubs,
	"schedule":  cmdSchedule,
}

var CfgFlags = map[string]string{
	"-u": "exp",
	"-M": "TCP_PING",
	"-m": "ST_ASYNC",
}

// DefaultConfig contains all the default speed test setting.
type DefaultConfig struct {
	Remarks  string
	Url      string
	Method   string
	Mode     string
	Interval int
	Chat     int64
	Include  []string
	Exclude  []string
}
