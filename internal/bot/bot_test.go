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
		Text:     "/schedult start",
		Entities: &[]tgbotapi.MessageEntity{{Offset: 0, Type: "bot_command", Length: 9}},
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

func TestSetDefMode(t *testing.T) {
	cmdSetDefaultModeAndMethod(NewBot(), NewMsg())
}

func TestSchedule(t *testing.T) {
	Def.Interval = -1
	Def.Url = "https://oxygenproxy.com"
	Def.Chat = 649191333
	go start(NewBot())
}

func TestCMDSchedule(t *testing.T) {
	cmdSchedule(NewBot(), NewMsg())
}
