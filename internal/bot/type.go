package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// BotConf store all bot config
type Conf struct {
	token string
}

type B = tgbotapi.BotAPI
type M = tgbotapi.Message
