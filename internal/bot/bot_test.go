package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"testing"
)

func TestNewBot(t *testing.T) {
	bot := NewBot()
	SendT(bot, 649191333, "New Bot Test")
}

func NewMsg() *M {
	return &tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: 649191333},
		Text:     "/default OxygenProxy",
		Entities: &[]tgbotapi.MessageEntity{{Offset: 0, Type: "bot_command", Length: 7}},
	}
}

func TestStart(t *testing.T) {
	cmdReadSub(NewBot(), NewMsg())
}

func TestCmdList(t *testing.T) {
	cmdListSubs(NewBot(), NewMsg())
}

func TestCMDSelectDef(t *testing.T) {
	cmdSelectDefaultSub(NewBot(), NewMsg())
}
