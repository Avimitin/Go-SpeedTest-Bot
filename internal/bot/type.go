package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// BotConf store all bot config
type Conf struct {
	token string
}

type B = tgbotapi.BotAPI
type M = tgbotapi.Message

type CMDFunc func(*B, *M)
type CMD = map[string]CMDFunc

var Commands = CMD{
	"start":     cmdStart,
	"ping":      cmdPing,
	"status":    cmdStatus,
	"read_sub":  cmdReadSub,
	"result":    cmdResult,
	"run_url":   cmdStartTestWithURL,
	"list_subs": cmdListSubs,
}
