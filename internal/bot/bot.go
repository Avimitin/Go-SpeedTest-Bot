package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go-speedtest-bot/internal/config"
	"log"
	"os"
)

func loadConf() *Conf {
	cfg := config.GetBotFile()
	s := cfg.Section("bot")
	return &Conf{token: s.Key("token").String()}
}

// NewBot return a bot instance
func NewBot() *B {
	bot, err := tgbotapi.NewBotAPI(loadConf().token)
	if err != nil {
		log.Println("[NewBotError]", err)
		os.Exit(-1)
	}
	return bot
}

// SendT send text message
func SendT(bot *B, cid int64, text string) {
	msg := tgbotapi.NewMessage(cid, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("SendError", err)
	}
}

func CMDHandler(bot *B, msg *M) {
	if msg.IsCommand() {
		if cmd, ok := Commands[msg.Text]; ok {
			cmd(bot, msg)
		}
	}
}

func cmdStart(b *B, m *M) {
	text := "Here is a bot who can help you manage all your proxy."
	SendT(b, m.Chat.ID, text)
}
